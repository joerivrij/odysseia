package models

// ErrorModel is the base model used for handling errors
type ErrorModel struct {
	UniqueCode string `json:"uniqueCode"`
}

// ValidationMessages messages used in validation error
type ValidationMessages struct {
	Field   string `json:"validationField"`
	Message string `json:"validationMessage"`
}

// ValidationError validation errors occur when data is malformed
type ValidationError struct {
	ErrorModel
	Messages []ValidationMessages `json:"errorModel"`
}

type NotFoundError struct {
	ErrorModel
	Message NotFoundMessage `json:"errorModel"`
}

func (m *NotFoundError) Error() string {
	return m.Error()
}

type NotFoundMessage struct {
	Type   string `json:"type"`
	Reason string `json:"reason"`
}

type ElasticSearchError struct {
	ErrorModel
	Message ElasticErrorMessage `json:"errorModel"`
}

func (m *ElasticSearchError) Error() string {
	return m.Error()
}

type ElasticErrorMessage struct {
	ElasticError string `json:"elasticError"`
}

// MethodMessages messages used in method error
type MethodMessages struct {
	Methods string `json:"allowedMethods"`
	Message string `json:"methodError"`
}

// method errors occur when calling endpoints with an unallowed method
type MethodError struct {
	ErrorModel
	Messages []MethodMessages `json:"errorModel"`
}
