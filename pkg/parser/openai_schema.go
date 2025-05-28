package parser

import (
	"encoding/json"

	"github.com/csotherden/prescription-parser/pkg/models"
	"github.com/invopop/jsonschema"
)

// PrescriptionResponseSchema is a JSON schema that defines the structure for prescription data extraction
// when using the OpenAI API. It's generated from the models.Prescription struct using reflection.
var PrescriptionResponseSchema = GenerateSchema[models.Prescription]()

// ResultScoreSchema is a JSON schema that defines the structure for result score data extraction
// when using the OpenAI API. It's generated from the models.ParserResultScore struct using reflection.
var ResultScoreSchema = GenerateSchema[models.ParserResultScore]()

// GenerateSchema creates a JSON schema from a Go type using reflection.
// It configures the schema generator to comply with the subset of JSON schema
// that is supported by the OpenAI API.
//
// The generic type T is the struct type from which to generate the schema.
// Returns a map representation of the JSON schema that can be provided to the API.
func GenerateSchema[T any]() map[string]interface{} {
	// Structured Outputs uses a subset of JSON schema
	// These flags are necessary to comply with the subset
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            true,
	}
	var v T
	schema := reflector.Reflect(v)
	schemaJson, err := schema.MarshalJSON()
	if err != nil {
		panic(err)
	}

	var schemaObj map[string]interface{}
	err = json.Unmarshal(schemaJson, &schemaObj)
	if err != nil {
		panic(err)
	}

	return schemaObj
}
