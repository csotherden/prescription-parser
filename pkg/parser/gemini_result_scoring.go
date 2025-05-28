package parser

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/csotherden/prescription-parser/pkg/models"
	"google.golang.org/genai"
)

// ScoreResult scores the result of a parser against a validated expected JSON.
func (p *GeminiParser) ScoreResult(ctx context.Context, expectedJSON, outputJSON string) (models.ParserResultScore, error) {
	cfg := &genai.GenerateContentConfig{
		ResponseMIMEType: "application/json",
		ResponseSchema:   &parserResultScoreSchema,
	}

	userParts := []*genai.Part{
		genai.NewPartFromText(fmt.Sprintf("%sHere are the JSON objects to compare:\n\nValidated Expected JSON:\n%s\n\nParser Output JSON:\n%s", scoringPrompt, expectedJSON, outputJSON)),
	}

	contents := []*genai.Content{
		genai.NewContentFromParts(userParts, genai.RoleUser),
	}

	resp, err := p.client.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash-preview-05-20",
		contents,
		cfg,
	)
	if err != nil {
		return models.ParserResultScore{}, fmt.Errorf("failed to score result: %w", err)
	}

	var score models.ParserResultScore
	err = json.Unmarshal([]byte(resp.Text()), &score)
	if err != nil {
		return models.ParserResultScore{}, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return score, nil
}
