package models

type Address struct {
	Street string `json:"street" jsonschema_description:"Street address of the location"`
	City   string `json:"city" jsonschema_description:"City name"`
	State  string `json:"state" jsonschema_description:"Two-letter state abbreviation (e.g., NY, CA)"`
	Zip    string `json:"zip" jsonschema_description:"ZIP or postal code"`
}
