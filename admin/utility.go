package admin

import (
	"errors"
)

var (
	// ErrParse error message for a parsing error
	ErrParse = errors.New("An error occurred parsing the response")
)
