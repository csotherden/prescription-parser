package models

type Medication struct {
	DrugName            string `json:"drug_name" jsonschema_description:"Name of the prescribed drug"`
	Ndc                 string `json:"ndc" jsonschema_description:"National Drug Code (NDC) for the medication"`
	Form                string `json:"form" jsonschema_description:"Dosage form (e.g., tablet, injection, packet)"`
	Strength            string `json:"strength" jsonschema_description:"Drug strength (e.g., 40 mg/0.4 mL)"`
	SIG                 string `json:"sig" jsonschema_description:"Verbatim instructions for administration (SIG) from the form"`
	Quantity            string `json:"quantity" jsonschema_description:"Amount of medication to dispense"`
	Refills             string `json:"refills" jsonschema_description:"Number of authorized refills"`
	StartDate           string `json:"start_date" jsonschema_description:"Date the patient should begin the medication (YYYY-MM-DD)"`
	Duration            string `json:"duration" jsonschema_description:"Intended treatment duration (e.g., 12 weeks)"`
	AdministrationNotes string `json:"administration_notes" jsonschema_description:"Plain English translation of SIG directions"`
	Indication          string `json:"indication" jsonschema_description:"Diagnosis or condition the drug is intended to treat"`
}
