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

// // Roll returns the result of a dice roll as a slice of numbers
// func (rd *RandomDice) RollDice() []uint {
// 	// remember returning a slice happens by reference, recipient could change the it
// 	// we can ensure protection with a defensive copy
// 	rollCopy := make([]uint, len(rd.Roll), len(rd.Roll))
// 	copy(rollCopy, rd.Roll)

// 	return rollCopy
// }

// // Sum returns the sum of all dice points in the roll
// func (rd *RandomDice) Sum() uint {
// 	return rd.RollSum
// }
