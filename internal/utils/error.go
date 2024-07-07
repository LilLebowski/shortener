// Package utils contains errors
package utils

import "fmt"

// UniqueConstraintError - struct error for UniqueConstraint
type UniqueConstraintError struct {
	Err error
}

// Error - error for UniqueConstraintError
func (uce *UniqueConstraintError) Error() string {
	return fmt.Sprintf("%v", uce.Err)
}

// Unwrap - error for UniqueConstraintError
func (uce *UniqueConstraintError) Unwrap() error {
	return uce.Err
}

// NewUniqueConstraintError - UniqueConstraintError initializer
func NewUniqueConstraintError(err error) error {
	return &UniqueConstraintError{
		Err: err,
	}
}

// DeletedError - struct error for already deleted item
type DeletedError struct {
	Message string
	Err     error
}

// Error for DeletedError
func (de *DeletedError) Error() string {
	return fmt.Sprintf("[%s] %v", de.Message, de.Err)
}

// Unwrap for DeletedError
func (de *DeletedError) Unwrap() error {
	return de.Err
}

// NewDeletedError - DeletedError initializer
func NewDeletedError(msg string, err error) error {
	return &DeletedError{
		Message: msg,
		Err:     err,
	}
}

// NotFoundError - struct error for not found item
type NotFoundError struct {
	Message string
	Err     error
}

// Error for NotFoundError
func (nfe *NotFoundError) Error() string {
	return fmt.Sprintf("[%s] %v", nfe.Message, nfe.Err)
}

// Unwrap for NotFoundError
func (nfe *NotFoundError) Unwrap() error {
	return nfe.Err
}

// NewNotFoundError - NotFoundError initializer
func NewNotFoundError(msg string, err error) error {
	return &NotFoundError{
		Message: msg,
		Err:     err,
	}
}
