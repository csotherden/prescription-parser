// Package parser provides functionality for parsing prescription images using different AI backends.
// It supports multiple AI services (OpenAI and Gemini) for processing prescription forms and
// extracting structured data from them, with vector embedding generation for similarity search.
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
	// ParseImage processes a prescription image asynchronously and returns a job ID for tracking parsing progress.
	// It takes a filename and file reader, initiates an asynchronous job, and returns the job ID.
	ParseImage(ctx context.Context, fileName string, file io.Reader) (string, error)

	// GetEmbedding generates an embedding vector for a prescription.
	// This vector representation can be used for similarity searches and document clustering.
	GetEmbedding(ctx context.Context, prescription models.Prescription) ([]float32, error)

	// UploadImage uploads an image to persistent storage and returns its ID.
	// The image can then be referenced in subsequent API calls.
	UploadImage(ctx context.Context, fileName string, file io.Reader) (string, error)

	// ScoreResult scores the result of a parser against a validated expected JSON.
	ScoreResult(ctx context.Context, expectedJSON, outputJSON string) (models.ParserResultScore, error)
}

// NewParser creates a new instance of a Parser implementation based on config.
// It returns the appropriate parser implementation (OpenAI or Gemini) based on the configuration.
// Returns an error if the parser backend specified in config is not supported.
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
