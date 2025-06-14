package domain

import (
	"errors"
)

type RandomNumber uint

// Not very happy about the exported fields...
type RandomDice struct {
	Roll     []uint
	RollSum  uint
	MaxDice  uint
	MaxSides uint
}

var (
	ErrInvalidNumDice  = errors.New("Value error for 'numdice': out of range")
	ErrInvalidNumSides = errors.New("Value error for 'numsides': out of range")
)

func (d *RandomDice) ValidateDiceParameters(numDice, sides uint) error {
	allErrors := make([]error, 0, 2)

	if numDice < 1 || numDice > d.MaxDice {
		allErrors = append(allErrors, ErrInvalidNumDice)
	}
	if sides < 2 || sides > d.MaxSides {
		allErrors = append(allErrors, ErrInvalidNumSides)
	}
	if len(allErrors) > 0 {
		return errors.Join(allErrors...)
	}

	return nil
}
