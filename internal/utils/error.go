package utils

import "fmt"

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
