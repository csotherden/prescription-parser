package models

import (
	"encoding/json"
	"testing"
)

func TestPrescriptionValidation(t *testing.T) {
	tests := []struct {
		name         string
		prescription Prescription
		wantErr      bool
	}{
		{
			name: "valid prescription",
			prescription: Prescription{
				Medications: []Medication{
					{
						DrugName: "Metformin",
						Strength: "500mg",
						Quantity: "30",
						Refills:  "3",
						SIG:      "Take 1 tablet by mouth daily",
					},
				},
				Patient: Patient{
					FirstName: "John",
					LastName:  "Doe",
					Dob:       "1990-01-01",
				},
				Prescriber: Prescriber{
					Name:      "Jane Smith",
					Npi:       "1234567890",
					Specialty: "MD",
				},
			},
			wantErr: false,
		},
		{
			name: "missing medication",
			prescription: Prescription{
				Medications: []Medication{},
				Patient: Patient{
					FirstName: "John",
					LastName:  "Doe",
					Dob:       "1990-01-01",
				},
				Prescriber: Prescriber{
					Name:      "Jane Smith",
					Npi:       "1234567890",
					Specialty: "MD",
				},
			},
			wantErr: true,
		},
		{
			name: "missing patient information",
			prescription: Prescription{
				Medications: []Medication{
					{
						DrugName: "Metformin",
						Strength: "500mg",
						Quantity: "30",
						Refills:  "3",
						SIG:      "Take 1 tablet by mouth daily",
					},
				},
				Patient: Patient{},
				Prescriber: Prescriber{
					Name:      "Jane Smith",
					Npi:       "1234567890",
					Specialty: "MD",
				},
			},
			wantErr: true,
		},
		{
			name: "missing prescriber information",
			prescription: Prescription{
				Medications: []Medication{
					{
						DrugName: "Metformin",
						Strength: "500mg",
						Quantity: "30",
						Refills:  "3",
						SIG:      "Take 1 tablet by mouth daily",
					},
				},
				Patient: Patient{
					FirstName: "John",
					LastName:  "Doe",
					Dob:       "1990-01-01",
				},
				Prescriber: Prescriber{},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test JSON marshaling and unmarshaling
			data, err := json.Marshal(tt.prescription)
			if err != nil {
				t.Fatalf("Failed to marshal prescription: %v", err)
			}

			var unmarshaled Prescription
			err = json.Unmarshal(data, &unmarshaled)
			if err != nil {
				t.Fatalf("Failed to unmarshal prescription: %v", err)
			}

			// Verify basic validation
			err = validatePrescription(tt.prescription)

			if (err != nil) != tt.wantErr {
				t.Errorf("Prescription validation error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// Helper function to validate a prescription
func validatePrescription(p Prescription) error {
	if len(p.Medications) == 0 {
		return &ValidationError{Field: "medications", Message: "at least one medication is required"}
	}

	if p.Patient.FirstName == "" || p.Patient.LastName == "" {
		return &ValidationError{Field: "patient", Message: "patient first and last name are required"}
	}

	if p.Prescriber.Name == "" || p.Prescriber.Npi == "" {
		return &ValidationError{Field: "prescriber", Message: "prescriber information is incomplete"}
	}

	return nil
}
