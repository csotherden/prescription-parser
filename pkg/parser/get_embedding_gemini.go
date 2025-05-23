package parser

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/csotherden/prescription-parser/pkg/models"
	"google.golang.org/genai"
)

var embeddingDimensionality int32 = 1536

func (p *PrescriptionParser) getEmbeddingGemini(ctx context.Context, prescription models.Prescription) ([]float32, error) {
	rxJson, err := json.Marshal(prescription)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal prescription: %w", err)
	}

	resp, err := p.geminiClient.Models.EmbedContent(
		ctx,
		"gemini-embedding-exp-03-07",
		[]*genai.Content{
			genai.NewContentFromText(string(rxJson), genai.RoleUser),
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
