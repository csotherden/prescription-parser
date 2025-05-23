package models

type DeliveryInfo struct {
	Destination string `json:"destination" jsonschema_description:"Where the prescription should be shipped (e.g., Patient's Home, Prescriber's Office)"`
	Notes       string `json:"notes" jsonschema_description:"Additional delivery instructions or details"`
}
