package datastore

import (
	"context"
	"fmt"

	"github.com/csotherden/prescription-parser/pkg/models"
	"github.com/pgvector/pgvector-go"
	"go.uber.org/zap"
)

// SaveSamplePrescription stores a prescription and its vector embedding in the database.
// It creates both the prescription record and its associated embedding in a single transaction.
//
// Parameters:
//   - ctx: Context for the database operation
//   - mimeType: MIME type of the prescription image
//   - imageID: ID of the image file associated with this prescription
//   - prescription: Prescription data to store
//   - embedding: Vector embedding representing the prescription content for similarity search
//
// Returns:
//   - An error if the database operation fails, nil on success
func (d *PgEntDatastore) SaveSamplePrescription(ctx context.Context, mimeType, imageID string, prescription models.Prescription, embedding []float32) error {
	tx, err := d.dbClient.Tx(ctx)
	if err != nil {
		d.logger.Error("failed to create transaction", zap.Error(err))
		return fmt.Errorf("failed to create transaction: %w", err)
	}
	defer tx.Rollback()

	dbPrescription, err := tx.Prescription.Create().
		SetFileID(imageID).
		SetMimeType(mimeType).
		SetContent(prescription).
		Save(ctx)
	if err != nil {
		d.logger.Error("failed to create prescription", zap.Error(err))
		return fmt.Errorf("failed to create prescription: %w", err)
	}

	_, err = tx.Embedding.Create().
		SetPrescriptionID(dbPrescription.ID).
		SetEmbedding(pgvector.NewVector(embedding)).
		Save(ctx)
	if err != nil {
		d.logger.Error("failed to create embedding", zap.Error(err))
		return fmt.Errorf("failed to create embedding: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		d.logger.Error("failed to commit transaction", zap.Error(err))
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
