package parser

import (
	"context"
	"fmt"
	"io"

	"github.com/csotherden/prescription-parser/pkg/config"
	"github.com/csotherden/prescription-parser/pkg/datastore"
	"github.com/csotherden/prescription-parser/pkg/models"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"go.uber.org/zap"
	"google.golang.org/genai"
)

// JobTypeParserConversation is the job type for Parser conversation processing
const JobTypeParserConversation = "parse_prescription"

// PrescriptionParser represents the AI prescription parser
type PrescriptionParser struct {
	ds              *datastore.Datastore
	logger          *zap.Logger
	openAIClient    openai.Client
	geminiClient    *genai.Client
	parseImageFunc  func(context.Context, string, string, io.Reader)
	embeddingFunc   func(context.Context, models.Prescription) ([]float32, error)
	uploadImageFunc func(context.Context, string, io.Reader) (string, error)
}

// NewParser creates a new instance of Parser
func NewParser(cfg config.Config, ds *datastore.Datastore, logger *zap.Logger) (*PrescriptionParser, error) {
	p := &PrescriptionParser{
		ds:     ds,
		logger: logger,
	}

	p.logger.Info("initializing parser", zap.String("parser_backend", cfg.ParserBackend))

	switch cfg.ParserBackend {
	case "OpenAI":
		p.openAIClient = openai.NewClient(
			option.WithAPIKey(cfg.OpenAIAPIKey),
		)
		p.parseImageFunc = p.parseImageOpenAI
		p.embeddingFunc = p.getEmbeddingOpenAI
		p.uploadImageFunc = p.uploadImageOpenAI
	case "Gemini":
		geminiClient, err := genai.NewClient(context.Background(), &genai.ClientConfig{
			APIKey:  cfg.GeminiAPIKey,
			Backend: genai.BackendGeminiAPI,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to initialize gemini client: %w", err)
		}

		p.geminiClient = geminiClient
		p.parseImageFunc = p.parseImageGemini
		p.embeddingFunc = p.getEmbeddingGemini
		p.uploadImageFunc = p.uploadImageGemini
	default:
		return nil, fmt.Errorf("unknown parser backend: %s. Must be OpenAI or Gemini", cfg.ParserBackend)
	}

	return p, nil
}
