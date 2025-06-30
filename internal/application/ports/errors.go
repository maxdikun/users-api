package ports

import "fmt"

type DuplicationError struct {
	Source string
	Object string
	Field  string
	Value  any
}

var _ error = (*DuplicationError)(nil)

// Error implements error.
func (err *DuplicationError) Error() string {
	return fmt.Sprintf(
		"duplication error occurred in %s storage: object '%s' with value '%v' in field '%s' already exists",
		err.Source,
		err.Object,
		err.Value,
		err.Field,
	)
}

type NotFoundError struct {
	Source string
	Object string
	Field  string
}

var _ error = (*NotFoundError)(nil)

// Error implements error.
func (err *NotFoundError) Error() string {
	return fmt.Sprintf(
		"object '%s' is not found in the storage %s, by field '%s'",
		err.Object,
		err.Source,
		err.Field,
	)
}
