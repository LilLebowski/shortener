// Package utils contains errors
package utils

import "fmt"

// UniqueConstraintError - error for UniqueConstraint
type UniqueConstraintError struct {
	Err error
}

func (uce *UniqueConstraintError) Error() string {
	return fmt.Sprintf("%v", uce.Err)
}

func (uce *UniqueConstraintError) Unwrap() error {
	return uce.Err
}

func NewUniqueConstraintError(err error) error {
	return &UniqueConstraintError{
		Err: err,
	}
}

// DeletedError - error for already deleted item
type DeletedError struct {
	Message string
	Err     error
}

func (de *DeletedError) Error() string {
	return fmt.Sprintf("[%s] %v", de.Message, de.Err)
}

func (de *DeletedError) Unwrap() error {
	return de.Err
}

func NewDeletedError(msg string, err error) error {
	return &DeletedError{
		Message: msg,
		Err:     err,
	}
}

// NotFoundError - error for not found item
type NotFoundError struct {
	Message string
	Err     error
}

func (nfe *NotFoundError) Error() string {
	return fmt.Sprintf("[%s] %v", nfe.Message, nfe.Err)
}

func (nfe *NotFoundError) Unwrap() error {
	return nfe.Err
}

func NewNotFoundError(msg string, err error) error {
	return &NotFoundError{
		Message: msg,
		Err:     err,
	}
}
