package parser

import (
	"testing"

	"github.com/csotherden/prescription-parser/pkg/config"
	"github.com/csotherden/prescription-parser/pkg/mocks"
	"go.uber.org/zap"
)

func TestNewParser(t *testing.T) {
	// Create a logger for testing
	logger := zap.NewNop()

	// Create a mock datastore
	mockDatastore := mocks.NewMockDatastore()

	tests := []struct {
		name         string
		config       config.Config
		expectedType string
		expectError  bool
	}{
		{
			name: "OpenAI parser",
			config: config.Config{
				ParserBackend: "OpenAI",
				OpenAIAPIKey:  "test-key",
			},
			expectedType: "*parser.OpenAIParser",
			expectError:  false,
		},
		{
			name: "Gemini parser",
			config: config.Config{
				ParserBackend: "Gemini",
				GeminiAPIKey:  "test-key",
			},
			expectedType: "*parser.GeminiParser",
			expectError:  false,
		},
		{
			name: "Unknown parser",
			config: config.Config{
				ParserBackend: "Unknown",
			},
			expectedType: "",
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser, err := NewParser(tt.config, mockDatastore, logger)

			// Check error expectation
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got nil")
				}
				return
			} else if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			// Check parser type
			parserType := getType(parser)
			if parserType != tt.expectedType {
				t.Errorf("Expected parser type %s, got %s", tt.expectedType, parserType)
			}
		})
	}
}

// Helper function to get the type name of an interface
func getType(p Parser) string {
	// If we want to check nil directly to avoid panic
	if p == nil {
		return "<nil>"
	}

	switch p.(type) {
	case *OpenAIParser:
		return "*parser.OpenAIParser"
	case *GeminiParser:
		return "*parser.GeminiParser"
	default:
		return "unknown"
	}
}
