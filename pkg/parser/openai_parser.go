package parser

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/csotherden/prescription-parser/pkg/config"
	"github.com/csotherden/prescription-parser/pkg/datastore"
	"github.com/csotherden/prescription-parser/pkg/jobs"
	"github.com/csotherden/prescription-parser/pkg/models"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/responses"
	"go.uber.org/zap"
)

// openAIEmbedding represents the structure of an embedding response from the OpenAI API.
// It contains the vector representation of a text input.
// We're using this to have a consistent float32 slice across all parsers.
type openAIEmbedding struct {
	Object    string    `json:"object"`
	Index     int       `json:"index"`
	Embedding []float32 `json:"embedding"`
}

// OpenAIParser implements the Parser interface using OpenAI services.
// It uses OpenAI's vision and embedding capabilities to process prescription images
// and extract structured data from them.
type OpenAIParser struct {
	ds     datastore.Datastore
	logger *zap.Logger
	client openai.Client
}

// NewOpenAIParser creates a new OpenAI-based parser.
// It initializes a connection to the OpenAI API with the provided API key
// and returns a parser instance ready for processing prescription images.
func NewOpenAIParser(cfg config.Config, ds datastore.Datastore, logger *zap.Logger) (*OpenAIParser, error) {
	client := openai.NewClient(
		option.WithAPIKey(cfg.OpenAIAPIKey),
	)

	return &OpenAIParser{
		ds:     ds,
		logger: logger,
		client: client,
	}, nil
}

// ParseImage handles parsing a prescription image using OpenAI vision API.
// It creates an asynchronous job to process the image and returns the job ID.
// The actual processing is done in a separate goroutine.
func (p *OpenAIParser) ParseImage(ctx context.Context, fileName string, file io.Reader) (string, error) {
	// Create a job for asynchronous processing
	jobID := jobs.GlobalTracker.CreateJob(
		JobTypeParsePrescription,
		fmt.Sprintf("Processing image: %s", fileName),
	)

	p.logger.Info("starting image parsing", zap.String("job_id", jobID), zap.String("file_name", fileName))

	go p.parseImageProcess(context.Background(), jobID, fileName, file)

	return jobID, nil
}

// deleteImage removes an image from the OpenAI API.
// It's used to clean up resources after processing is complete.
func (p *OpenAIParser) deleteImage(ctx context.Context, imageID string) error {
	_, err := p.client.Files.Delete(ctx, imageID)
	if err != nil {
		return fmt.Errorf("failed to delete image: %w", err)
	}
	return nil
}

// parseImageProcess processes the image asynchronously.
// It validates the file type, uploads it to OpenAI, performs parsing passes,
// and updates the job status throughout the process. It also cleans up
// the uploaded files when done.
func (p *OpenAIParser) parseImageProcess(ctx context.Context, jobID, fileName string, file io.Reader) {
	fileExt := strings.ToLower(filepath.Ext(fileName))
	var contentType string
	switch fileExt {
	case ".pdf":
		contentType = "application/pdf"
	default:
		jobs.GlobalTracker.UpdateJob(jobID, jobs.JobStatusFailed, fmt.Errorf("unsupported file type. file must be PDF not %s", fileExt), nil)
		return
	}

	inputFile := openai.File(file, fileName, contentType)

	storedFile, err := p.client.Files.New(ctx, openai.FileNewParams{
		File:    inputFile,
		Purpose: openai.FilePurposeUserData,
	})
	if err != nil {
		jobs.GlobalTracker.UpdateJob(jobID, jobs.JobStatusFailed, fmt.Errorf("failed to upload file: %w", err), nil)
		p.logger.Error("failed to upload file", zap.String("job_id", jobID), zap.String("file_name", fileName), zap.Error(err))
		return
	}

	// Update job status to processing
	jobs.GlobalTracker.UpdateJob(jobID, jobs.JobStatusProcessing, nil, nil)

	defer func() {
		err := p.deleteImage(ctx, storedFile.ID)
		if err != nil {
			p.logger.Error("failed to delete image", zap.String("image_id", storedFile.ID), zap.Error(err))
		}
	}()

	// Initial parsing pass
	rx, err := p.firstParsingPass(ctx, storedFile.ID)
	if err != nil {
		p.logger.Error("failed in first parsing pass", zap.String("job_id", jobID), zap.Error(err))
		jobs.GlobalTracker.UpdateJob(jobID, jobs.JobStatusFailed, fmt.Errorf("failed in first parsing pass: %w", err), nil)
		return
	}

	p.logger.Info("first parsing pass completed", zap.String("job_id", jobID), zap.String("file_name", fileName))

	// Get embedding for the parsed prescription
	embedding, err := p.GetEmbedding(ctx, rx)
	if err != nil {
		p.logger.Error("failed to get embedding", zap.String("job_id", jobID), zap.String("file_name", fileName), zap.Error(err))
		jobs.GlobalTracker.UpdateJob(jobID, jobs.JobStatusComplete, nil, rx)
		return
	}

	// Get similar samples
	samples, err := p.ds.GetSamples(ctx, embedding)
	if err != nil {
		p.logger.Error("failed to get samples", zap.String("job_id", jobID), zap.String("file_name", fileName), zap.Error(err))
		jobs.GlobalTracker.UpdateJob(jobID, jobs.JobStatusComplete, nil, rx)
		return
	}

	p.logger.Info("sample images loaded", zap.String("job_id", jobID), zap.String("file_name", fileName), zap.Int("sample_count", len(samples)))

	// Second parsing pass with examples
	if len(samples) > 0 {
		secondPassRx, err := p.secondParsingPass(ctx, storedFile.ID, samples, rx)
		if err != nil {
			p.logger.Error("failed in second parsing pass", zap.String("job_id", jobID), zap.Error(err))
			jobs.GlobalTracker.UpdateJob(jobID, jobs.JobStatusComplete, nil, rx)
			return
		}
		rx = secondPassRx
	}

	p.logger.Info("successfully processed image", zap.String("job_id", jobID), zap.String("file_name", fileName))
	jobs.GlobalTracker.UpdateJob(jobID, jobs.JobStatusComplete, nil, rx)
}

// firstParsingPass performs the initial parsing of the prescription.
// It sends the prescription image to OpenAI API with system and user prompts
// to extract structured data from the image.
func (p *OpenAIParser) firstParsingPass(ctx context.Context, fileID string) (models.Prescription, error) {
	messages := []responses.ResponseInputItemUnionParam{
		responses.ResponseInputItemParamOfMessage(
			systemPrompt,
			"system"),
	}

	imageMessage := responses.ResponseInputItemParamOfMessage(
		responses.ResponseInputMessageContentListParam{
			responses.ResponseInputContentUnionParam{
				OfInputFile: &responses.ResponseInputFileParam{
					FileID: openai.String(fileID),
					Type:   "input_file",
				},
			},
			responses.ResponseInputContentUnionParam{
				OfInputText: &responses.ResponseInputTextParam{
					Text: parsePrompt,
					Type: "input_text",
				},
			},
		},
		"user",
	)

	messages = append(messages, imageMessage)

	params := responses.ResponseNewParams{
		Text: responses.ResponseTextConfigParam{
			Format: responses.ResponseFormatTextConfigUnionParam{
				OfJSONSchema: &responses.ResponseFormatTextJSONSchemaConfigParam{
					Name:        "Prescription",
					Schema:      PrescriptionResponseSchema,
					Strict:      openai.Bool(true),
					Description: openai.String("Prescription Image Parser Prescription JSON"),
					Type:        "json_schema",
				},
			},
		},
		Model: "gpt-4.1-2025-04-14",
		Input: responses.ResponseNewParamsInputUnion{
			OfInputItemList: messages,
		},
		MaxOutputTokens: openai.Int(10240),
	}

	resp, err := p.client.Responses.New(ctx, params)
	if err != nil {
		return models.Prescription{}, fmt.Errorf("failed to process image: %w", err)
	}

	var rx models.Prescription
	err = json.Unmarshal([]byte(resp.OutputText()), &rx)
	if err != nil {
		return models.Prescription{}, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return rx, nil
}

// secondParsingPass performs a fresh parsing pass with example context.
// It uses similar prescription samples to refine the initial parsing results,
// potentially improving accuracy by learning from precedents.
func (p *OpenAIParser) secondParsingPass(ctx context.Context, fileID string, samples []models.SamplePrescription, firstPassRx models.Prescription) (models.Prescription, error) {
	messages := []responses.ResponseInputItemUnionParam{
		responses.ResponseInputItemParamOfMessage(
			systemPrompt,
			"system"),
	}

	for _, sample := range samples {
		sampleMessage := responses.ResponseInputItemParamOfMessage(
			responses.ResponseInputMessageContentListParam{
				responses.ResponseInputContentUnionParam{
					OfInputFile: &responses.ResponseInputFileParam{
						FileID: openai.String(sample.FileID),
						Type:   "input_file",
					},
				},
				responses.ResponseInputContentUnionParam{
					OfInputText: &responses.ResponseInputTextParam{
						Text: parsePrompt,
						Type: "input_text",
					},
				},
			},
			"user",
		)

		sampleResponse := responses.ResponseInputItemParamOfMessage(
			responses.ResponseInputMessageContentListParam{
				responses.ResponseInputContentUnionParam{
					OfInputText: &responses.ResponseInputTextParam{
						Text: sample.Content,
						Type: "output_text",
					},
				},
			},
			"assistant",
		)

		messages = append(messages, sampleMessage)
		messages = append(messages, sampleResponse)
	}

	reviewMessage := responses.ResponseInputItemParamOfMessage(
		responses.ResponseInputMessageContentListParam{
			responses.ResponseInputContentUnionParam{
				OfInputFile: &responses.ResponseInputFileParam{
					FileID: openai.String(fileID),
					Type:   "input_file",
				},
			},
			responses.ResponseInputContentUnionParam{
				OfInputText: &responses.ResponseInputTextParam{
					Text: reviewPrompt,
					Type: "input_text",
				},
			},
		},
		"user",
	)

	messages = append(messages, reviewMessage)

	params := responses.ResponseNewParams{
		Text: responses.ResponseTextConfigParam{
			Format: responses.ResponseFormatTextConfigUnionParam{
				OfJSONSchema: &responses.ResponseFormatTextJSONSchemaConfigParam{
					Name:        "Prescription",
					Schema:      PrescriptionResponseSchema,
					Strict:      openai.Bool(true),
					Description: openai.String("Prescription Image Parser Prescription JSON"),
					Type:        "json_schema",
				},
			},
		},
		Model: "gpt-4.1-2025-04-14",
		Input: responses.ResponseNewParamsInputUnion{
			OfInputItemList: messages,
		},
		MaxOutputTokens: openai.Int(10240),
	}

	resp, err := p.client.Responses.New(ctx, params)
	if err != nil {
		return firstPassRx, fmt.Errorf("failed to run second pass: %w", err)
	}

	var secondPassRx models.Prescription
	err = json.Unmarshal([]byte(resp.OutputText()), &secondPassRx)
	if err != nil {
		return firstPassRx, fmt.Errorf("failed to unmarshal second pass response: %w", err)
	}

	return secondPassRx, nil
}

// GetEmbedding generates embeddings for a prescription using OpenAI.
// It converts the prescription to JSON and sends it to the OpenAI API to generate
// a vector representation for similarity search.
func (p *OpenAIParser) GetEmbedding(ctx context.Context, prescription models.Prescription) ([]float32, error) {
	jsonBytes, err := json.Marshal(prescription)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal prescription: %w", err)
	}

	resp, err := p.client.Embeddings.New(ctx, openai.EmbeddingNewParams{
		Input: openai.EmbeddingNewParamsInputUnion{
			OfString: openai.String(string(jsonBytes)),
		},
		Model:          openai.EmbeddingModelTextEmbedding3Small,
		Dimensions:     openai.Int(1536),
		EncodingFormat: "float",
	})
	if err != nil || len(resp.Data) == 0 {
		return nil, fmt.Errorf("failed to generate prescription embedding: %w", err)
	}

	var emb openAIEmbedding

	err = json.Unmarshal([]byte(resp.Data[0].RawJSON()), &emb)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal prescription embedding: %w", err)
	}

	return emb.Embedding, nil
}

// UploadImage uploads an image to the OpenAI API.
// It validates the file type, uploads the file to OpenAI, and returns an ID
// that can be used to reference the image in subsequent API calls.
func (p *OpenAIParser) UploadImage(ctx context.Context, fileName string, file io.Reader) (string, error) {
	fileExt := strings.ToLower(filepath.Ext(fileName))
	var contentType string
	switch fileExt {
	case ".pdf":
		contentType = "application/pdf"
	default:
		return "", fmt.Errorf("unsupported file type: %s. file must be PDF", fileExt)
	}

	inputFile := openai.File(file, fileName, contentType)

	storedFile, err := p.client.Files.New(ctx, openai.FileNewParams{
		File:    inputFile,
		Purpose: openai.FilePurposeUserData,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload image: %w", err)
	}

	return storedFile.ID, nil
}
