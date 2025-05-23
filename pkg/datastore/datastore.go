package datastore

import (
	"context"
	"encoding/json"
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

type Datastore struct {
	dbClient *ent.Client
	logger   *zap.Logger
}

// SaveSample saves a sample prescription with its validated JSON
func (d *Datastore) SaveSample(ctx context.Context, fileID, contentType, jsonContent string, embedding []float32) error {
	var prescription models.Prescription
	if err := json.Unmarshal([]byte(jsonContent), &prescription); err != nil {
		return fmt.Errorf("failed to unmarshal prescription JSON: %w", err)
	}

	return d.SaveSamplePrescription(ctx, contentType, fileID, prescription, embedding)
}

func NewDatastore(cfg config.Config, logger *zap.Logger) (*Datastore, error) {
	dbClient, err := newEntClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to initialize datastore: %w", err)
	}

	return &Datastore{
		dbClient: dbClient,
		logger:   logger,
	}, nil
}

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
