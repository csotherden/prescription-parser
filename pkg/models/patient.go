package models

type Patient struct {
	FirstName        string        `json:"first_name" jsonschema_description:"Patient's first name"`
	MiddleName       string        `json:"middle_name" jsonschema_description:"Patient's middle name"`
	LastName         string        `json:"last_name" jsonschema_description:"Patient's last name"`
	Dob              string        `json:"dob" jsonschema_description:"Patient's date of birth (YYYY-MM-DD)"`
	Sex              string        `json:"sex" jsonschema_description:"Patient's biological sex (e.g., Male, Female, Other)"`
	Weight           Measurement   `json:"weight" jsonschema_description:"Patient's recorded weight"`
	Height           Measurement   `json:"height" jsonschema_description:"Patient's recorded height"`
	Address          Address       `json:"address" jsonschema_description:"Patient's residential address"`
	PhoneNumbers     []PhoneNumber `json:"phone_numbers" jsonschema_description:"Patient's contact phone numbers"`
	Allergies        []string      `json:"allergies" jsonschema_description:"List of known allergies"`
	EmergencyContact Contact       `json:"emergency_contact" jsonschema_description:"Emergency contact details for the patient"`
	Insurance        []Insurance   `json:"insurance" jsonschema_description:"List of the patient's insurance policies"`
}
