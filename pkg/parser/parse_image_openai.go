package parser

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/csotherden/prescription-parser/pkg/jobs"
	"github.com/csotherden/prescription-parser/pkg/models"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/responses"
	"go.uber.org/zap"
)

func (p *PrescriptionParser) parseImageOpenAI(ctx context.Context, jobID, fileName string, file io.Reader) {
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

	storedFile, err := p.openAIClient.Files.New(ctx, openai.FileNewParams{
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
		err := p.DeleteImage(ctx, storedFile.ID)
		if err != nil {
			p.logger.Error("failed to delete image", zap.String("image_id", storedFile.ID), zap.Error(err))
		}
	}()

	messages := []responses.ResponseInputItemUnionParam{
		responses.ResponseInputItemParamOfMessage(
			systemPrompt,
			"system"),
	}

	imageMessage := responses.ResponseInputItemParamOfMessage(
		responses.ResponseInputMessageContentListParam{
			responses.ResponseInputContentUnionParam{
				OfInputFile: &responses.ResponseInputFileParam{
					FileID: openai.String(storedFile.ID),
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

	resp, err := p.openAIClient.Responses.New(ctx, params)
	if err != nil {
		p.logger.Error("failed to process image", zap.String("image_id", storedFile.ID), zap.Error(err))
		jobs.GlobalTracker.UpdateJob(jobID, jobs.JobStatusFailed, fmt.Errorf("failed to process image: %w", err), "")
		return
	}

	var rx models.Prescription
	err = json.Unmarshal([]byte(resp.OutputText()), &rx)
	if err != nil {
		p.logger.Error("failed to process image", zap.String("job_id", jobID), zap.String("file_name", fileName), zap.Error(err))
		jobs.GlobalTracker.UpdateJob(jobID, jobs.JobStatusFailed, fmt.Errorf("failed to process image: %w", err), "")
		return
	}

	p.logger.Info("first parsing pass completed", zap.String("job_id", jobID), zap.String("file_name", fileName))

	embedding, err := p.getEmbeddingOpenAI(ctx, rx)
	if err != nil {
		p.logger.Error("failed to get embedding", zap.String("job_id", jobID), zap.String("file_name", fileName), zap.Error(err))
		jobs.GlobalTracker.UpdateJob(jobID, jobs.JobStatusComplete, nil, rx)
		return
	}

	samples, err := p.ds.GetSamples(ctx, embedding)
	if err != nil {
		p.logger.Error("failed to get samples", zap.String("job_id", jobID), zap.String("file_name", fileName), zap.Error(err))
		jobs.GlobalTracker.UpdateJob(jobID, jobs.JobStatusComplete, nil, rx)
		return
	}

	p.logger.Info("sample images loaded", zap.String("job_id", jobID), zap.String("file_name", fileName), zap.Int("sample_count", len(samples)))

	messages = []responses.ResponseInputItemUnionParam{
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
						Type: "input_text",
					},
				},
			},
			"assistant",
		)

		messages = append(messages, sampleMessage)
		messages = append(messages, sampleResponse)
	}

	messages = append(messages, imageMessage)

	firstPassResponse := responses.ResponseInputItemParamOfMessage(
		responses.ResponseInputMessageContentListParam{
			responses.ResponseInputContentUnionParam{
				OfInputText: &responses.ResponseInputTextParam{
					Text: resp.OutputText(),
					Type: "input_text",
				},
			},
		},
		"assistant",
	)

	messages = append(messages, firstPassResponse)

	messages = append(
		messages,
		responses.ResponseInputItemParamOfMessage(
			reviewPrompt,
			"user",
		),
	)

	p.logger.Info("starting second parsing pass (review)", zap.String("job_id", jobID), zap.String("file_name", fileName))

	resp, err = p.openAIClient.Responses.New(ctx, params)
	if err != nil {
		p.logger.Error("failed to run second pass for prescription", zap.String("job_id", jobID), zap.String("file_name", fileName), zap.Error(err))
		jobs.GlobalTracker.UpdateJob(jobID, jobs.JobStatusComplete, nil, rx)
		return
	}

	var secondPassRx models.Prescription
	err = json.Unmarshal([]byte(resp.OutputText()), &secondPassRx)
	if err != nil {
		p.logger.Error("failed to process second pass image", zap.String("job_id", jobID), zap.String("file_name", fileName), zap.Error(err))
		jobs.GlobalTracker.UpdateJob(jobID, jobs.JobStatusComplete, nil, rx)
	}

	p.logger.Info("successfully processed image", zap.String("job_id", jobID), zap.String("file_name", fileName))

	jobs.GlobalTracker.UpdateJob(jobID, jobs.JobStatusComplete, nil, &secondPassRx)
}
