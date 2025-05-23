package models

type Prescription struct {
	DateWritten         string            `json:"date_written" jsonschema_description:"Date the prescription was written (YYYY-MM-DD)"`
	DateNeeded          string            `json:"date_needed" jsonschema_description:"Date by which the medication is needed (YYYY-MM-DD)"`
	Patient             Patient           `json:"patient,omitzero" jsonschema_description:"Demographic and insurance details of the patient"`
	Prescriber          Prescriber        `json:"prescriber,omitzero" jsonschema_description:"Information about the prescribing healthcare provider"`
	Diagnosis           PatientDiagnosis  `json:"diagnosis" jsonschema_description:"Clinical diagnosis details associated with the prescription"`
	ClinicalInfo        []string          `json:"clinical_info" jsonschema_description:"Additional clinical notes, such as lab values, genetic markers, or BSA"`
	Medications         []Medication      `json:"medications" jsonschema_description:"List of medications prescribed on this form"`
	TherapyStatus       string            `json:"therapy_status" jsonschema_description:"Indicates whether the therapy is new, restarted, or ongoing"`
	FailedTherapies     []TherapyHistory  `json:"failed_therapies" jsonschema_description:"List of prior therapies the patient tried and discontinued"`
	Delivery            DeliveryInfo      `json:"delivery,omitzero" jsonschema_description:"Shipping instructions for the medication"`
	PrescriberSignature SignatureInfo     `json:"prescriber_signature" jsonschema_description:"Signature and DAW code authorization from the prescriber"`
	Attachments         AttachmentDetails `json:"attachments" jsonschema_description:"Boolean indicators for supplemental documents provided with the form"`
}
