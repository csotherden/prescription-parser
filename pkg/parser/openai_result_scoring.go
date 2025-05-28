package parser

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/csotherden/prescription-parser/pkg/models"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/responses"
)

// ScoreResult scores the result of a parser against a validated expected JSON.
func (p *OpenAIParser) ScoreResult(ctx context.Context, expectedJSON, outputJSON string) (models.ParserResultScore, error) {
	messages := []responses.ResponseInputItemUnionParam{
		responses.ResponseInputItemParamOfMessage(
			responses.ResponseInputMessageContentListParam{
				responses.ResponseInputContentUnionParam{
					OfInputText: &responses.ResponseInputTextParam{
						Text: fmt.Sprintf("%sHere are the JSON objects to compare:\n\nValidated Expected JSON:\n%s\n\nParser Output JSON:\n%s", scoringPrompt, expectedJSON, outputJSON),
						Type: "input_text",
					},
				},
			},
			"user"),
	}

	params := responses.ResponseNewParams{
		Text: responses.ResponseTextConfigParam{
			Format: responses.ResponseFormatTextConfigUnionParam{
				OfJSONSchema: &responses.ResponseFormatTextJSONSchemaConfigParam{
					Name:        "ParserResultScore",
					Schema:      ResultScoreSchema,
					Strict:      openai.Bool(true),
					Description: openai.String("Parser Result Score JSON"),
					Type:        "json_schema",
				},
			},
		},
		Model: "gpt-4.1-2025-04-14",
		Input: responses.ResponseNewParamsInputUnion{
			OfInputItemList: messages,
		},
		MaxOutputTokens: openai.Int(10240),
	}

	resp, err := p.client.Responses.New(ctx, params)
	if err != nil {
		return models.ParserResultScore{}, fmt.Errorf("failed to score result: %w", err)
	}

	var score models.ParserResultScore
	err = json.Unmarshal([]byte(resp.OutputText()), &score)
	if err != nil {
		return models.ParserResultScore{}, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return score, nil
}
