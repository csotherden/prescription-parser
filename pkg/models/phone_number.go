package models

type PhoneNumber struct {
	Label     string `json:"label" jsonschema_description:"Label used to identify the phone number e.g. Home, Mobile, Work"`
	Number    string `json:"number" jsonschema_description:"The numeric phone number without any spaces or formating characters"`
	Extension string `json:"extension" jsonschema_description:"The extension to dial (if any)"`
}
