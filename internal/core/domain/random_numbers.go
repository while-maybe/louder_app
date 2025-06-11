package domain

import (
	"errors"
)

type RandomNumber uint

// Not very happy about the exported fields...

type RandomDice struct {
	Roll    []uint
	RollSum uint
}

func ValidateDiceParameters(numDice, sides uint) error {
	allErrors := make([]error, 0, 2)

	if numDice < 1 || numDice > 10 {
		allErrors = append(allErrors, errors.New("Dice number must be between 1 and 10"))
	}
	if sides < 3 || sides > 20 {
		allErrors = append(allErrors, errors.New("Dice sides must be between 3 and 20"))
	}
	if len(allErrors) > 0 {
		return errors.Join(allErrors...)
	}

	return nil
}

// Roll returns the result of a dice roll as a slice of numbers
func (rd *RandomDice) RollDice() []uint {
	// remember returning a slice happens by reference, recipient could change the it
	// we can ensure protection with a defensive copy
	rollCopy := make([]uint, len(rd.Roll), len(rd.Roll))
	copy(rollCopy, rd.Roll)

	return rollCopy
}

// Sum returns the sum of all dice points in the roll
func (rd *RandomDice) Sum() uint {
	return rd.RollSum
}
