package models

// ErrorModel is the base model used for handling errors
type ErrorModel struct {
	UniqueCode string `json:"uniqueCode"`
}

// messages used in validation error
type ValidationMessages struct {
	Field    string `json:"validationField"`
	Message    string `json:"validationMessage"`
}

// validation errors occur when data is malformed
type ValidationError struct {
	ErrorModel
	Messages   []ValidationMessages `json:"errorModel"`
}

type NotFoundError struct {
	ErrorModel
	Message NotFoundMessage `json:"errorModel"`
}

type NotFoundMessage struct {
	Type string `json:"type"`
	Reason string `json:"reason"`
}

//  messages used in method error
type MethodMessages struct {
	Methods    string `json:"allowedMethods"`
	Message    string `json:"methodError"`
}

// method errors occur when calling endpoints with an unallowed method
type MethodError struct {
	ErrorModel
	Messages   []MethodMessages `json:"errorModel"`
}
