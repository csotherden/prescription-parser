package parser

import (
	"google.golang.org/genai"
)

// geminiSchema defines the structure for prescription data extraction using structured output.
// It describes the expected JSON format for the prescription data with detailed field descriptions
// to help the model correctly parse and structure the prescription information.
var geminiSchema = genai.Schema{
	Type:        "OBJECT",
	Description: "Prescription form data structure",
	Properties: map[string]*genai.Schema{
		"date_written": {
			Type:        "STRING",
			Description: "Date the prescription was written (YYYY-MM-DD)",
		},
		"date_needed": {
			Type:        "STRING",
			Description: "Date by which the medication is needed (YYYY-MM-DD)",
		},
		"patient": {
			Type:        "OBJECT",
			Description: "Demographic and insurance details of the patient",
			Properties: map[string]*genai.Schema{
				"first_name": {
					Type:        "STRING",
					Description: "Patient's first name",
				},
				"middle_name": {
					Type:        "STRING",
					Description: "Patient's middle name",
				},
				"last_name": {
					Type:        "STRING",
					Description: "Patient's last name",
				},
				"dob": {
					Type:        "STRING",
					Description: "Patient's date of birth (YYYY-MM-DD)",
				},
				"sex": {
					Type:        "STRING",
					Description: "Patient's biological sex (e.g., Male, Female, Other)",
				},
				"weight": {
					Type:        "OBJECT",
					Description: "Patient's recorded weight",
					Properties: map[string]*genai.Schema{
						"unit": {
							Type:        "STRING",
							Description: "The unit the measurement was taken in",
						},
						"value": {
							Type:        "STRING",
							Description: "The recorded value of the measurement",
						},
					},
				},
				"height": {
					Type:        "OBJECT",
					Description: "Patient's recorded height",
					Properties: map[string]*genai.Schema{
						"unit": {
							Type:        "STRING",
							Description: "The unit the measurement was taken in",
						},
						"value": {
							Type:        "STRING",
							Description: "The recorded value of the measurement",
						},
					},
				},
				"address": {
					Type:        "OBJECT",
					Description: "Patient's residential address",
					Properties: map[string]*genai.Schema{
						"street": {
							Type:        "STRING",
							Description: "Street address of the location",
						},
						"city": {
							Type:        "STRING",
							Description: "City name",
						},
						"state": {
							Type:        "STRING",
							Description: "Two-letter state abbreviation (e.g., NY, CA)",
						},
						"zip": {
							Type:        "STRING",
							Description: "ZIP or postal code",
						},
					},
				},
				"phone_numbers": {
					Type:        "ARRAY",
					Description: "Patient's contact phone numbers",
					Items: &genai.Schema{
						Type:        "OBJECT",
						Description: "Phone number details",
						Properties: map[string]*genai.Schema{
							"label": {
								Type:        "STRING",
								Description: "Label used to identify the phone number e.g. Home, Mobile, Work",
							},
							"number": {
								Type:        "STRING",
								Description: "The numeric phone number without any spaces or formating characters",
							},
							"extension": {
								Type:        "STRING",
								Description: "The extension to dial (if any)",
							},
						},
					},
				},
				"allergies": {
					Type:        "ARRAY",
					Description: "List of known allergies",
					Items: &genai.Schema{
						Type: "STRING",
					},
				},
				"emergency_contact": {
					Type:        "OBJECT",
					Description: "Emergency contact details for the patient",
					Properties: map[string]*genai.Schema{
						"name": {
							Type:        "STRING",
							Description: "Full name of the contact person",
						},
						"relationship": {
							Type:        "STRING",
							Description: "The patient's relationship to the contact person (e.g., spouse, parent)",
						},
						"phone": {
							Type:        "STRING",
							Description: "Phone number for the emergency or alternate contact",
						},
					},
				},
				"insurance": {
					Type:        "ARRAY",
					Description: "List of the patient's insurance policies",
					Items: &genai.Schema{
						Type:        "OBJECT",
						Description: "Insurance policy details",
						Properties: map[string]*genai.Schema{
							"type": {
								Type:        "STRING",
								Description: "Primary or Secondary insurance type",
							},
							"provider": {
								Type:        "STRING",
								Description: "Name of the insurance provider",
							},
							"id_number": {
								Type:        "STRING",
								Description: "Patient's insurance ID number",
							},
							"group_number": {
								Type:        "STRING",
								Description: "Insurance group number",
							},
							"rx_bin": {
								Type:        "STRING",
								Description: "Prescription BIN (Bank Identification Number)",
							},
							"pcn": {
								Type:        "STRING",
								Description: "Processor Control Number for pharmacy claims",
							},
							"policyholder_name": {
								Type:        "STRING",
								Description: "Full name of the insurance policyholder",
							},
							"policyholder_dob": {
								Type:        "STRING",
								Description: "Date of birth of the policyholder (YYYY-MM-DD)",
							},
							"phone_number": {
								Type:        "STRING",
								Description: "Phone number of the insurance provider",
							},
						},
					},
				},
			},
		},
		"prescriber": {
			Type:        "OBJECT",
			Description: "Information about the prescribing healthcare provider",
			Properties: map[string]*genai.Schema{
				"name": {
					Type:        "STRING",
					Description: "Full name of the prescriber",
				},
				"specialty": {
					Type:        "STRING",
					Description: "Medical specialty of the prescriber (e.g., Oncology, Dermatology)",
				},
				"npi": {
					Type:        "STRING",
					Description: "National Provider Identifier (NPI) of the prescriber",
				},
				"state_license": {
					Type:        "STRING",
					Description: "Prescriber's state license number",
				},
				"dea": {
					Type:        "STRING",
					Description: "Prescriber's DEA number for controlled substances",
				},
				"office": {
					Type:        "OBJECT",
					Description: "Details of the prescriber's office or practice",
					Properties: map[string]*genai.Schema{
						"name": {
							Type:        "STRING",
							Description: "Name of the prescriber's office or medical facility",
						},
						"address": {
							Type:        "OBJECT",
							Description: "Physical address of the prescriber's office",
							Properties: map[string]*genai.Schema{
								"street": {
									Type:        "STRING",
									Description: "Street address of the location",
								},
								"city": {
									Type:        "STRING",
									Description: "City name",
								},
								"state": {
									Type:        "STRING",
									Description: "Two-letter state abbreviation (e.g., NY, CA)",
								},
								"zip": {
									Type:        "STRING",
									Description: "ZIP or postal code",
								},
							},
						},
						"phone": {
							Type:        "STRING",
							Description: "Main phone number for the office",
						},
						"fax": {
							Type:        "STRING",
							Description: "Fax number for the office",
						},
						"contact_name": {
							Type:        "STRING",
							Description: "Name of the designated office contact person",
						},
						"contact_email": {
							Type:        "STRING",
							Description: "Email address of the office contact person",
						},
					},
				},
			},
		},
		"diagnosis": {
			Type:        "OBJECT",
			Description: "Clinical diagnosis details associated with the prescription",
			Properties: map[string]*genai.Schema{
				"date_of_diagnosis": {
					Type:        "STRING",
					Description: "Date when the diagnosis was made (YYYY-MM-DD)",
				},
				"primary_diagnosis": {
					Type:        "OBJECT",
					Description: "The primary diagnosis for which the medication is prescribed",
					Properties: map[string]*genai.Schema{
						"description": {
							Type:        "STRING",
							Description: "Text description of the diagnosis (e.g., Psoriatic Arthritis)",
						},
						"icd10_code": {
							Type:        "STRING",
							Description: "ICD-10 code for the diagnosis (e.g., L40.50)",
						},
					},
				},
				"additional_diagnoses": {
					Type:        "ARRAY",
					Description: "Any additional diagnoses relevant to the patient",
					Items: &genai.Schema{
						Type:        "OBJECT",
						Description: "Additional diagnosis",
						Properties: map[string]*genai.Schema{
							"description": {
								Type:        "STRING",
								Description: "Text description of the diagnosis",
							},
							"icd10_code": {
								Type:        "STRING",
								Description: "ICD-10 code for the diagnosis",
							},
						},
					},
				},
			},
		},
		"clinical_info": {
			Type:        "ARRAY",
			Description: "Additional clinical notes, such as lab values, genetic markers, or BSA",
			Items: &genai.Schema{
				Type: "STRING",
			},
		},
		"medications": {
			Type:        "ARRAY",
			Description: "List of medications prescribed on this form",
			Items: &genai.Schema{
				Type:        "OBJECT",
				Description: "Medication details",
				Properties: map[string]*genai.Schema{
					"drug_name": {
						Type:        "STRING",
						Description: "Name of the prescribed drug",
					},
					"ndc": {
						Type:        "STRING",
						Description: "National Drug Code (NDC) for the medication",
					},
					"form": {
						Type:        "STRING",
						Description: "Dosage form (e.g., tablet, injection, packet)",
					},
					"strength": {
						Type:        "STRING",
						Description: "Drug strength (e.g., 40 mg/0.4 mL)",
					},
					"sig": {
						Type:        "STRING",
						Description: "Verbatim instructions for administration (SIG) from the form",
					},
					"quantity": {
						Type:        "STRING",
						Description: "Amount of medication to dispense",
					},
					"refills": {
						Type:        "STRING",
						Description: "Number of authorized refills",
					},
					"start_date": {
						Type:        "STRING",
						Description: "Date the patient should begin the medication (YYYY-MM-DD)",
					},
					"duration": {
						Type:        "STRING",
						Description: "Intended treatment duration (e.g., 12 weeks)",
					},
					"administration_notes": {
						Type:        "STRING",
						Description: "Plain English translation of SIG directions",
					},
					"indication": {
						Type:        "STRING",
						Description: "Diagnosis or condition the drug is intended to treat",
					},
				},
			},
		},
		"therapy_status": {
			Type:        "STRING",
			Description: "Indicates whether the therapy is new, restarted, or ongoing",
		},
		"failed_therapies": {
			Type:        "ARRAY",
			Description: "List of prior therapies the patient tried and discontinued",
			Items: &genai.Schema{
				Type:        "OBJECT",
				Description: "Previous therapy history",
				Properties: map[string]*genai.Schema{
					"name": {
						Type:        "STRING",
						Description: "Name of the previous therapy or medication",
					},
					"reason_for_discontinuation": {
						Type:        "STRING",
						Description: "Reason why the previous therapy was stopped",
					},
				},
			},
		},
		"delivery": {
			Type:        "OBJECT",
			Description: "Shipping instructions for the medication",
			Properties: map[string]*genai.Schema{
				"destination": {
					Type:        "STRING",
					Description: "Where the prescription should be shipped (e.g., Patient's Home, Prescriber's Office)",
				},
				"notes": {
					Type:        "STRING",
					Description: "Additional delivery instructions or details",
				},
			},
		},
		"prescriber_signature": {
			Type:        "OBJECT",
			Description: "Signature and DAW code authorization from the prescriber",
			Properties: map[string]*genai.Schema{
				"date": {
					Type:        "STRING",
					Description: "Date the prescription was signed by the prescriber (YYYY-MM-DD)",
				},
				"daw_code": {
					Type:        "STRING",
					Description: "Dispense As Written (DAW) code indicating substitution permission (e.g., 0 = substitution allowed, 1 = brand medically necessary)",
				},
			},
		},
		"attachments": {
			Type:        "OBJECT",
			Description: "Boolean indicators for supplemental documents provided with the form",
			Properties: map[string]*genai.Schema{
				"insurance_cards": {
					Type:        "BOOLEAN",
					Description: "Whether a copy of the insurance card (front and back) is attached",
				},
				"lab_results": {
					Type:        "BOOLEAN",
					Description: "Whether recent laboratory results are included",
				},
				"pathology_reports": {
					Type:        "BOOLEAN",
					Description: "Whether a pathology report is attached",
				},
				"clinical_notes": {
					Type:        "BOOLEAN",
					Description: "Whether recent clinical or office notes are included",
				},
				"other_documents": {
					Type:        "BOOLEAN",
					Description: "Whether any other relevant documents are attached",
				},
			},
		},
	},
}
