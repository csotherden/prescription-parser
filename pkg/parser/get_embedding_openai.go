package parser

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/openai/openai-go"

	"github.com/csotherden/prescription-parser/pkg/models"
)

type openAIEmbedding struct {
	Object    string    `json:"object"`
	Index     int       `json:"index"`
	Embedding []float32 `json:"embedding"`
}

func (p *PrescriptionParser) getEmbeddingOpenAI(ctx context.Context, prescription models.Prescription) ([]float32, error) {
	rxJson, err := json.Marshal(prescription)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal prescription: %w", err)
	}

	resp, err := p.openAIClient.Embeddings.New(ctx, openai.EmbeddingNewParams{
		Input: openai.EmbeddingNewParamsInputUnion{
			OfString: openai.String(string(rxJson)),
		},
		Model:          openai.EmbeddingModelTextEmbedding3Small,
		Dimensions:     openai.Int(1536),
		EncodingFormat: "float",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate prescription embedding: %w", err)
	}

	var emb openAIEmbedding
	err = json.Unmarshal([]byte(resp.Data[0].RawJSON()), &emb)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal prescription embedding: %w", err)
	}

	return emb.Embedding, nil
}
