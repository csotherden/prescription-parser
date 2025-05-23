package models

type TherapyHistory struct {
	Name                     string `json:"name" jsonschema_description:"Name of the previous therapy or medication"`
	ReasonForDiscontinuation string `json:"reason_for_discontinuation" jsonschema_description:"Reason why the previous therapy was stopped"`
}
