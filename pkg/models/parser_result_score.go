package models

// FieldScore represents the evaluation of a single field in the parsing results
type FieldScore struct {
	FieldPath     string  `json:"field_path" jsonschema_description:"Path to the evaluated field in dot notation (e.g. patient.first_name, medications[0].form)"`
	ExpectedValue any     `json:"expected_value" jsonschema_description:"The value expected for this field, can be null if not present in expected data"`
	OutputValue   any     `json:"output_value" jsonschema_description:"The value produced by the parser, can be null if missing"`
	Score         float64 `json:"score" jsonschema_description:"Score between 0.0 and 1.0 indicating correctness of the parsed value. Possible scores are 0.0, 0.25, 0.75, and 1.0."`
	Reasoning     string  `json:"reasoning" jsonschema_description:"Explanation for why this score was assigned"`
}

// ParserResultScore represents the overall evaluation of parser results
type ParserResultScore struct {
	FieldScores            []FieldScore `json:"field_scores" jsonschema_description:"Detailed scoring information for each evaluated field"`
	TotalAwardedPoints     float64      `json:"total_awarded_points" jsonschema_description:"Sum of points awarded across all evaluated fields. Each field is scored between 0.0 and 1.0. Possible scores are 0.0, 0.25, 0.75, and 1.0."`
	TotalPossiblePoints    float64      `json:"total_possible_points" jsonschema_description:"Maximum possible points if all fields were perfectly parsed. Up to one point is awarded for each field."`
	OverallScorePercentage float64      `json:"overall_score_percentage" jsonschema_description:"Percentage score calculated as (awarded/possible)*100"`
	SummaryCritique        string       `json:"summary_critique" jsonschema_description:"Overall assessment of the parser output quality based on scores. This is a human-readable summary of the parser's performance."`
}
