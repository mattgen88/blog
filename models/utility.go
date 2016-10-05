package models

import "errors"

var (
	ValidationError = errors.New("An error occurred validating the model.")
	SaveError       = errors.New("An error occurred saving the model.")
	DNEError        = errors.New("An error occurred finding the requested model.")
)
