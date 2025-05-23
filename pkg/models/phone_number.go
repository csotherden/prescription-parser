package models

type PhoneNumbers struct {
	Daytime string `json:"daytime" jsonschema_description:"Primary daytime contact number for the patient"`
	Evening string `json:"evening" jsonschema_description:"Evening contact number for the patient"`
	Cell    string `json:"cell" jsonschema_description:"Mobile phone number for the patient"`
}
