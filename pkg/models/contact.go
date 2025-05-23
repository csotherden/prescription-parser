package models

type Contact struct {
	Name         string `json:"name" jsonschema_description:"Full name of the contact person"`
	Relationship string `json:"relationship" jsonschema_description:"The patient's relationship to the contact person (e.g., spouse, parent)"`
	Phone        string `json:"phone" jsonschema_description:"Phone number for the emergency or alternate contact"`
}
