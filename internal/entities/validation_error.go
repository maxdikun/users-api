package entities

type ValidationError struct {
	Field   string
	Message string
}

var _ error = (*ValidationError)(nil)

// Error implements error.
func (v *ValidationError) Error() string {
	panic("unimplemented")
}

func newValidationError(field, message string) *ValidationError {
	return &ValidationError{
		Field:   field,
		Message: message,
	}
}
