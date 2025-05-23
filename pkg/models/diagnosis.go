package models

type Diagnosis struct {
	Description string `json:"description" jsonschema_description:"Text description of the diagnosis (e.g., Psoriatic Arthritis)"`
	Icd10Code   string `json:"icd10_code" jsonschema_description:"ICD-10 code for the diagnosis (e.g., L40.50)"`
}

type PatientDiagnosis struct {
	DateOfDiagnosis     string      `json:"date_of_diagnosis" jsonschema_description:"Date when the diagnosis was made (YYYY-MM-DD)"`
	PrimaryDiagnosis    Diagnosis   `json:"primary_diagnosis" jsonschema_description:"The primary diagnosis for which the medication is prescribed"`
	AdditionalDiagnoses []Diagnosis `json:"additional_diagnoses" jsonschema_description:"Any additional diagnoses relevant to the patient"`
}
