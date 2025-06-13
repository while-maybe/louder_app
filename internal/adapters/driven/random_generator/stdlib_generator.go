package randomgenerator

import (
	"fmt"
	"louder/internal/core/domain"
	"louder/internal/core/service/randomnumbers"
	"math/rand/v2"
)

type StdLibGenerator struct{}

var _ randomnumbers.Repository = (*StdLibGenerator)(nil)

func NewStdLibGenerator() *StdLibGenerator {
	return &StdLibGenerator{}
}

func (s *StdLibGenerator) GenerateRandomNumber() domain.RandomNumber {
	return domain.RandomNumber(rand.Uint())
}

const (
	maxDice  = 10
	maxSides = 20
)

func (s *StdLibGenerator) GenerateDiceRoll(numDice, sides uint) (*domain.RandomDice, error) {

	newDiceRoll := domain.RandomDice{
		MaxDice:  maxDice,
		MaxSides: maxSides,
	}

	if err := newDiceRoll.ValidateDiceParameters(numDice, sides); err != nil {
		return nil, fmt.Errorf("error adapter failed: %w", err)
	}

	newDiceRoll.Roll = make([]uint, numDice, numDice)
	// var sum uint

	for i := range numDice {
		newDiceRoll.Roll[i] = uint(rand.IntN(int(sides)) + 1)
		newDiceRoll.RollSum += newDiceRoll.Roll[i]
	}

	return &newDiceRoll, nil
}
