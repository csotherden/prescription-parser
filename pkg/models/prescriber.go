package models

type Prescriber struct {
	Name         string           `json:"name" jsonschema_description:"Full name of the prescriber"`
	Specialty    string           `json:"specialty" jsonschema_description:"Medical specialty of the prescriber (e.g., Oncology, Dermatology)"`
	Npi          string           `json:"npi" jsonschema_description:"National Provider Identifier (NPI) of the prescriber"`
	StateLicense string           `json:"state_license" jsonschema_description:"Prescriber's state license number"`
	Dea          string           `json:"dea" jsonschema_description:"Prescriber's DEA number for controlled substances"`
	Office       PrescriberOffice `json:"office" jsonschema_description:"Details of the prescriber's office or practice"`
}

type PrescriberOffice struct {
	Name         string  `json:"name" jsonschema_description:"Name of the prescriber's office or medical facility"`
	Address      Address `json:"address" jsonschema_description:"Physical address of the prescriber's office"`
	Phone        string  `json:"phone" jsonschema_description:"Main phone number for the office"`
	Fax          string  `json:"fax" jsonschema_description:"Fax number for the office"`
	ContactName  string  `json:"contact_name" jsonschema_description:"Name of the designated office contact person"`
	ContactEmail string  `json:"contact_email" jsonschema_description:"Email address of the office contact person"`
}
