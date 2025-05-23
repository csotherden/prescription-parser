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

// GetSamples retrieves prescription samples that are most similar to the provided embedding vector.
// It uses pgvector's similarity search to find the closest matches in the embedding space.
// The method returns up to 3 most similar samples, ordered by vector similarity.
//
// Parameters:
//   - ctx: Context for the database operation
//   - embedding: Vector embedding to use for similarity search
//
// Returns:
//   - A slice of SamplePrescription that are most similar to the embedding
//   - An error if the database operation fails
func (d *PgEntDatastore) GetSamples(ctx context.Context, embedding []float32) ([]models.SamplePrescription, error) {
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
