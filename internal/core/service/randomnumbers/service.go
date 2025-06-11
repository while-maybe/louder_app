// Implements the driving ports. Contains the core application logic. It depends on driven port interfaces, not concrete implementations. This is key for DI.
package randomnumbers

import (
	"fmt"
	"louder/internal/core/domain"
)

type Service struct {
	repo Repository // will represent the injected dependency
}

// check if Service implements the Port
var _ Port = (*Service)(nil)

// NewRandNumberService creates an instance of the Service struct
func NewRandNumberService(repo Repository) Port {
	return &Service{
		repo: repo,
	}
}

// what behaviours/actions/jobs should this service be able to perform and offer to whoever calls it?

// GetRandomNumber returns a uint random number
func (s *Service) GetRandomNumber() domain.RandomNumber {
	rn := s.repo.GenerateRandomNumber()
	return rn
}

// RollDice takes the number of dice and number of sides per dice and return the result as a RandomDice object or an error
func (s *Service) RollDice(numDice, numSides uint) (*domain.RandomDice, error) {
	roll, err := s.repo.GenerateDiceRoll(numDice, numSides)
	if err != nil {
		return nil, fmt.Errorf("Error rolling dice: %w", err)
	}

	return roll, nil
}
