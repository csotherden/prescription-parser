package datastore

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pgvector/pgvector-go"

	"entgo.io/ent/dialect/sql"
	"github.com/csotherden/prescription-parser/pkg/models"
	"go.uber.org/zap"
)

func (d *Datastore) GetSamples(ctx context.Context, embedding []float32) ([]models.SamplePrescription, error) {
	var samples []models.SamplePrescription

	tx, err := d.dbClient.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	embVec := pgvector.NewVector(embedding)
	embs, err := tx.Embedding.Query().
		Order(func(s *sql.Selector) {
			s.OrderExpr(sql.ExprP("embedding <-> $1", embVec))
		}).
		WithPrescription().
		Limit(3).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get samples: %w", err)
	}

	for _, emb := range embs {
		if emb.Edges.Prescription != nil {
			content, err := json.Marshal(emb.Edges.Prescription.Content)
			if err != nil {
				d.logger.Error("failed to marshal prescription content", zap.Error(err))
				continue
			}

			samples = append(samples, models.SamplePrescription{
				ID:       emb.Edges.Prescription.ID,
				FileID:   emb.Edges.Prescription.FileID,
				MIMEType: emb.Edges.Prescription.MimeType,
				Content:  string(content),
			})
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return samples, nil
}
