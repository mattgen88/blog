package models

import "errors"

var (
	ErrValidation   = errors.New("An error occurred validating the model.")
	ErrSave         = errors.New("An error occurred saving the model.")
	ErrDoesNotExist = errors.New("An error occurred finding the requested model.")
)
