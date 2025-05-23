package parser

import (
	"context"
	"fmt"
	"io"

	"github.com/csotherden/prescription-parser/pkg/config"
	"github.com/csotherden/prescription-parser/pkg/datastore"
	"github.com/csotherden/prescription-parser/pkg/models"
	"go.uber.org/zap"
)

// Parser defines the interface for prescription parsing services
type Parser interface {
	// ParseImage processes a prescription image and returns a job ID for tracking
	ParseImage(ctx context.Context, fileName string, file io.Reader) (string, error)

	// GetEmbedding generates an embedding vector for a prescription
	GetEmbedding(ctx context.Context, prescription models.Prescription) ([]float32, error)

	// UploadImage uploads an image to the backend service and returns its ID
	UploadImage(ctx context.Context, fileName string, file io.Reader) (string, error)
}

// NewParser creates a new instance of a Parser implementation based on config
func NewParser(cfg config.Config, ds datastore.Datastore, logger *zap.Logger) (Parser, error) {
	logger.Info("initializing parser", zap.String("parser_backend", cfg.ParserBackend))

	switch cfg.ParserBackend {
	case "OpenAI":
		return NewOpenAIParser(cfg, ds, logger)
	case "Gemini":
		return NewGeminiParser(cfg, ds, logger)
	default:
		return nil, fmt.Errorf("unknown parser backend: %s. Must be OpenAI or Gemini", cfg.ParserBackend)
	}
}
