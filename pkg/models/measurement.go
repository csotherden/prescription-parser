package models

type Measurement struct {
	Unit  string `json:"unit" jsonschema_description:"The unit the measurement was taken in"`
	Value string `json:"value" jsonschema_description:"The recorded value of the measurement"`
}
