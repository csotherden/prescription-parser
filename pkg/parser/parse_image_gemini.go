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
	"go.uber.org/zap"
	"google.golang.org/genai"
)

func (p *PrescriptionParser) parseImageGemini(ctx context.Context, jobID, fileName string, file io.Reader) {
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

	resp, err := p.geminiClient.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash-preview-05-20",
		contents,
		config,
	)
	if err != nil {
		return
	}

	var rx models.Prescription
	err = json.Unmarshal([]byte(resp.Text()), &rx)
	if err != nil {
		p.logger.Error("failed to process image", zap.String("job_id", jobID), zap.String("file_name", fileName), zap.Error(err))
		jobs.GlobalTracker.UpdateJob(jobID, jobs.JobStatusFailed, fmt.Errorf("failed to process image: %w", err), "")
		return
	}

	p.logger.Info("first parsing pass completed", zap.String("job_id", jobID), zap.String("file_name", fileName))

	embedding, err := p.getEmbeddingGemini(ctx, rx)
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
	history = append(history, genai.NewContentFromText(resp.Text(), genai.RoleModel))

	userParts = []*genai.Part{
		{
			InlineData: &genai.Blob{
				MIMEType: contentType,
				Data:     fileBytes,
			},
		},
		genai.NewPartFromText(parsePrompt),
	}

	contents = []*genai.Content{
		genai.NewContentFromParts(userParts, genai.RoleUser),
	}

	p.logger.Info("starting second parsing pass (review)", zap.String("job_id", jobID), zap.String("file_name", fileName))

	chat, err := p.geminiClient.Chats.Create(
		ctx,
		"gemini-2.5-flash-preview-05-20",
		config,
		history,
	)
	if err != nil {
		p.logger.Error("failed to initiate chat session", zap.String("job_id", jobID), zap.String("file_name", fileName), zap.Error(err))
		jobs.GlobalTracker.UpdateJob(jobID, jobs.JobStatusComplete, nil, rx)
		return
	}

	resp, err = chat.SendMessage(
		ctx,
		genai.Part{
			Text: reviewPrompt,
		},
	)
	if err != nil {
		p.logger.Error("failed to run second pass for prescription", zap.String("job_id", jobID), zap.String("file_name", fileName), zap.Error(err))
		jobs.GlobalTracker.UpdateJob(jobID, jobs.JobStatusComplete, nil, rx)
		return
	}

	var secondPassRx models.Prescription
	err = json.Unmarshal([]byte(resp.Text()), &secondPassRx)
	if err != nil {
		p.logger.Error("failed to process second pass image", zap.String("job_id", jobID), zap.String("file_name", fileName), zap.Error(err))
		jobs.GlobalTracker.UpdateJob(jobID, jobs.JobStatusComplete, nil, rx)
	}

	p.logger.Info("successfully processed image", zap.String("job_id", jobID), zap.String("file_name", fileName))

	jobs.GlobalTracker.UpdateJob(jobID, jobs.JobStatusComplete, nil, &secondPassRx)
}
