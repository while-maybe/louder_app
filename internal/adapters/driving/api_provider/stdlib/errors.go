package stdlibapiadapter

import "errors"

// Sentinel errors for randomnumbers
var (
	ErrMissingNumDice  = errors.New("Missing required parameter: 'numdice'")
	ErrMissingNumSides = errors.New("Missing required parameter: 'numsides'")
	ErrFormatNumDice   = errors.New("Invalid format for 'numdice': must be a valid integer.")
	ErrFormatNumSides  = errors.New("Invalid format for 'numsides': must be a valid integer.")
	ErrValueNumDice    = errors.New("Value error for 'numdice': must be a positive number")
	ErrValueNumSides   = errors.New("Value error for 'numsides': must be a positive number")
)
