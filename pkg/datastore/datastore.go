// Package datastore provides database access functionality for prescription data storage and retrieval.
// It implements a PostgreSQL-based storage system using Ent ORM framework with pgvector support.
package datastore

import (
	"context"
	"fmt"
	"time"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	"github.com/csotherden/prescription-parser/ent"
	"github.com/csotherden/prescription-parser/pkg/config"
	"github.com/csotherden/prescription-parser/pkg/models"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

// Datastore defines the interface for data persistence operations.
// It provides methods for retrieving and storing prescription data along with vector embeddings.
type Datastore interface {
	// GetSamples retrieves prescription samples similar to the provided embedding vector.
	// It returns a list of samples ordered by vector similarity.
	GetSamples(ctx context.Context, embedding []float32) ([]models.SamplePrescription, error)

	// SaveSamplePrescription stores a prescription sample along with its vector embedding.
	// It associates the prescription with the given image ID and MIME type.
	SaveSamplePrescription(ctx context.Context, mimeType, imageID string, prescription models.Prescription, embedding []float32) error
}

// PgEntDatastore implements the Datastore interface using PostgreSQL with Ent ORM.
// It uses pgvector for vector embedding storage and similarity search.
type PgEntDatastore struct {
	dbClient *ent.Client
	logger   *zap.Logger
}

// NewPgEntDatastore creates a new PostgreSQL-based Datastore implementation.
// It initializes the database connection and returns a ready-to-use datastore.
func NewPgEntDatastore(cfg config.Config, logger *zap.Logger) (Datastore, error) {
	dbClient, err := newEntClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to initialize datastore: %w", err)
	}

	return &PgEntDatastore{
		dbClient: dbClient,
		logger:   logger,
	}, nil
}

// newEntClient creates and configures a new Ent client with the provided configuration.
// It sets up connection pooling and runs schema migrations if needed.
func newEntClient(cfg config.Config) (*ent.Client, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		cfg.DatabaseHost, cfg.DatabasePort, cfg.DatabaseUser, cfg.DatabasePassword, cfg.DatabaseName)

	// Create driver with MaxIdleConns and MaxOpenConns
	drv, err := sql.Open(dialect.Postgres, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed opening connection to postgres: %w", err)
	}

	// Get the underlying sql.DB object
	db := drv.DB()
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)
	db.SetConnMaxLifetime(time.Hour)

	// Create the ent client
	client := ent.NewClient(ent.Driver(drv))

	// Run the auto migration
	if err := client.Schema.Create(context.Background()); err != nil {
		client.Close()
		return nil, fmt.Errorf("failed creating schema resources: %w", err)
	}
	return client, nil
}
