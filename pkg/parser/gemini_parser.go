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
	"go.uber.org/zap"
	"google.golang.org/genai"
)

// GeminiParser implements the Parser interface using Google Gemini services.
// It leverages Gemini's multimodal capabilities to process prescription images
// and extract structured data from them.
type GeminiParser struct {
	ds     datastore.Datastore
	logger *zap.Logger
	client *genai.Client
}

// NewGeminiParser creates a new Gemini-based parser.
// It initializes a client for the Gemini API with the provided API key
// and returns a parser instance ready for processing prescription images.
func NewGeminiParser(cfg config.Config, ds datastore.Datastore, logger *zap.Logger) (*GeminiParser, error) {
	client, err := genai.NewClient(context.Background(), &genai.ClientConfig{
		APIKey:  cfg.GeminiAPIKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize gemini client: %w", err)
	}

	return &GeminiParser{
		ds:     ds,
		logger: logger,
		client: client,
	}, nil
}

// ParseImage handles parsing a prescription image using Gemini multimodal API.
// It creates an asynchronous job to process the image and returns the job ID.
// The actual processing is done in a separate goroutine.
func (p *GeminiParser) ParseImage(ctx context.Context, fileName string, file io.Reader) (string, error) {
	// Create a job for asynchronous processing
	jobID := jobs.GlobalTracker.CreateJob(
		JobTypeParsePrescription,
		fmt.Sprintf("Processing image: %s", fileName),
	)

	p.logger.Info("starting image parsing", zap.String("job_id", jobID), zap.String("file_name", fileName))

	go p.parseImageProcess(context.Background(), jobID, fileName, file)

	return jobID, nil
}

// parseImageProcess processes the image asynchronously.
// It reads the file contents, validates the file type, performs parsing passes,
// and updates the job status throughout the process.
func (p *GeminiParser) parseImageProcess(ctx context.Context, jobID, fileName string, file io.Reader) {
	fileExt := strings.ToLower(filepath.Ext(fileName))
	var contentType string
	switch fileExt {
	case ".pdf":
		contentType = "application/pdf"
	default:
		jobs.GlobalTracker.UpdateJob(jobID, jobs.JobStatusFailed, fmt.Errorf("unsupported file type. file must be PDF not %s", fileExt), nil)
		return
	}

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		jobs.GlobalTracker.UpdateJob(jobID, jobs.JobStatusFailed, fmt.Errorf("failed to read file contents: %w", err), nil)
		return
	}

	// Update job status to processing
	jobs.GlobalTracker.UpdateJob(jobID, jobs.JobStatusProcessing, nil, nil)

	// Initial parsing pass
	rx, err := p.firstParsingPass(ctx, contentType, fileBytes)
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
		secondPassRx, err := p.secondParsingPass(ctx, contentType, fileBytes, samples, rx)
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
// It sends the prescription image to Gemini API with system and user prompts
// to extract structured data from the image.
func (p *GeminiParser) firstParsingPass(ctx context.Context, contentType string, fileBytes []byte) (models.Prescription, error) {
	config := &genai.GenerateContentConfig{
		SystemInstruction: genai.NewContentFromText(systemPrompt, genai.RoleUser),
		ResponseMIMEType:  "application/json",
		ResponseSchema:    &geminiSchema,
	}

	userParts := []*genai.Part{
		{
			InlineData: &genai.Blob{
				MIMEType: contentType,
				Data:     fileBytes,
			},
		},
		genai.NewPartFromText(parsePrompt),
	}

	contents := []*genai.Content{
		genai.NewContentFromParts(userParts, genai.RoleUser),
	}

	resp, err := p.client.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash-preview-05-20",
		contents,
		config,
	)
	if err != nil {
		return models.Prescription{}, fmt.Errorf("failed to process image: %w", err)
	}

	var rx models.Prescription
	err = json.Unmarshal([]byte(resp.Text()), &rx)
	if err != nil {
		return models.Prescription{}, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return rx, nil
}

// secondParsingPass performs a review with example context.
// It uses similar prescription samples to refine the initial parsing results,
// potentially improving accuracy by learning from precedents.
func (p *GeminiParser) secondParsingPass(ctx context.Context, contentType string, fileBytes []byte, samples []models.SamplePrescription, firstPassRx models.Prescription) (models.Prescription, error) {
	history := []*genai.Content{}

	for _, sample := range samples {
		sampleParts := []*genai.Part{
			genai.NewPartFromURI(sample.FileID, sample.MIMEType),
			genai.NewPartFromText(parsePrompt),
		}

		history = append(history, genai.NewContentFromParts(sampleParts, genai.RoleUser))
		history = append(history, genai.NewContentFromText(sample.Content, genai.RoleModel))
	}

	firstPassParts := []*genai.Part{
		{
			InlineData: &genai.Blob{
				MIMEType: contentType,
				Data:     fileBytes,
			},
		},
		genai.NewPartFromText(parsePrompt),
	}

	history = append(history, genai.NewContentFromParts(firstPassParts, genai.RoleUser))

	firstPassText, err := json.Marshal(firstPassRx)
	if err != nil {
		return firstPassRx, fmt.Errorf("failed to marshal first pass response: %w", err)
	}

	history = append(history, genai.NewContentFromText(string(firstPassText), genai.RoleModel))

	config := &genai.GenerateContentConfig{
		SystemInstruction: genai.NewContentFromText(systemPrompt, genai.RoleUser),
		ResponseMIMEType:  "application/json",
		ResponseSchema:    &geminiSchema,
	}

	chat, err := p.client.Chats.Create(
		ctx,
		"gemini-2.5-flash-preview-05-20",
		config,
		history,
	)
	if err != nil {
		return firstPassRx, fmt.Errorf("failed to initiate chat session: %w", err)
	}

	resp, err := chat.SendMessage(
		ctx,
		genai.Part{
			Text: reviewPrompt,
		},
	)
	if err != nil {
		return firstPassRx, fmt.Errorf("failed to run second pass: %w", err)
	}

	var secondPassRx models.Prescription
	err = json.Unmarshal([]byte(resp.Text()), &secondPassRx)
	if err != nil {
		return firstPassRx, fmt.Errorf("failed to unmarshal second pass response: %w", err)
	}

	return secondPassRx, nil
}

// GetEmbedding generates embeddings for a prescription using Gemini embeddings API.
// It converts the prescription to JSON and sends it to the Gemini API to generate
// a vector representation for similarity search.
func (p *GeminiParser) GetEmbedding(ctx context.Context, prescription models.Prescription) ([]float32, error) {
	jsonBytes, err := json.Marshal(prescription)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal prescription: %w", err)
	}

	embeddingDimensionality := int32(1536)
	resp, err := p.client.Models.EmbedContent(
		ctx,
		"gemini-embedding-exp-03-07",
		[]*genai.Content{
			genai.NewContentFromText(string(jsonBytes), genai.RoleUser),
		},
		&genai.EmbedContentConfig{
			TaskType:             "SEMANTIC_SIMILARITY",
			OutputDimensionality: &embeddingDimensionality,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate prescription embedding: %w", err)
	}

	return resp.Embeddings[0].Values, nil
}

// UploadImage uploads an image using the Gemini Files API.
// It validates the file type, uploads the file to Gemini, and returns a URI
// that can be used to reference the image in subsequent API calls.
func (p *GeminiParser) UploadImage(ctx context.Context, fileName string, file io.Reader) (string, error) {
	fileExt := strings.ToLower(filepath.Ext(fileName))
	var contentType string
	switch fileExt {
	case ".pdf":
		contentType = "application/pdf"
	default:
		return "", fmt.Errorf("unsupported file type: %s. file must be PDF", fileExt)
	}

	resp, err := p.client.Files.Upload(ctx, file, &genai.UploadFileConfig{
		MIMEType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload image: %w", err)
	}

	return resp.URI, nil
}
