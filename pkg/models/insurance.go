package models

type Insurance struct {
	Type             string `json:"type" jsonschema_description:"Primary or Secondary insurance type"`
	Provider         string `json:"provider" jsonschema_description:"Name of the insurance provider"`
	IdNumber         string `json:"id_number" jsonschema_description:"Patient's insurance ID number"`
	GroupNumber      string `json:"group_number" jsonschema_description:"Insurance group number"`
	RxBin            string `json:"rx_bin" jsonschema_description:"Prescription BIN (Bank Identification Number)"`
	Pcn              string `json:"pcn" jsonschema_description:"Processor Control Number for pharmacy claims"`
	PolicyholderName string `json:"policyholder_name" jsonschema_description:"Full name of the insurance policyholder"`
	PolicyholderDob  string `json:"policyholder_dob" jsonschema_description:"Date of birth of the policyholder (YYYY-MM-DD)"`
	PhoneNumber      string `json:"phone_number" jsonschema_description:"Phone number of the insurance provider"`
}
