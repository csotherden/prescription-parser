package models

type AttachmentDetails struct {
	InsuranceCards   bool `json:"insurance_cards" jsonschema_description:"Whether a copy of the insurance card (front and back) is attached"`
	LabResults       bool `json:"lab_results" jsonschema_description:"Whether recent laboratory results are included"`
	PathologyReports bool `json:"pathology_reports" jsonschema_description:"Whether a pathology report is attached"`
	ClinicalNotes    bool `json:"clinical_notes" jsonschema_description:"Whether recent clinical or office notes are included"`
	OtherDocuments   bool `json:"other_documents" jsonschema_description:"Whether any other relevant documents are attached"`
}
