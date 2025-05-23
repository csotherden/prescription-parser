package models

import "github.com/google/uuid"

type SamplePrescription struct {
	ID       uuid.UUID `json:"id"`
	FileID   string    `json:"file_id"`
	MIMEType string    `json:"mime_type"`
	Content  string    `json:"content"`
}
