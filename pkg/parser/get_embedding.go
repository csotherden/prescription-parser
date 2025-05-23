package parser

import (
	"context"

	"github.com/csotherden/prescription-parser/pkg/models"
)

func (p *PrescriptionParser) GetEmbedding(ctx context.Context, prescription models.Prescription) ([]float32, error) {
	return p.embeddingFunc(ctx, prescription)
}
