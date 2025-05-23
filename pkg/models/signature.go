package models

type SignatureInfo struct {
	Date    string `json:"date" jsonschema_description:"Date the prescription was signed by the prescriber (YYYY-MM-DD)"`
	DawCode string `json:"daw_code" jsonschema_description:"Dispense As Written (DAW) code indicating substitution permission (e.g., 0 = substitution allowed, 1 = brand medically necessary)"`
}
