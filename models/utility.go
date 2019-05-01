package models

import "errors"

// Error messages
var (
	ErrValidation   = errors.New("an error occurred validating the model")
	ErrSave         = errors.New("an error occurred saving the model")
	ErrDoesNotExist = errors.New("an error occurred finding the requested model")
	ErrDelete       = errors.New("an error occurred in deleting the model")
)
