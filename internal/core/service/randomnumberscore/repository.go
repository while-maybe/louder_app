// Interfaces that the core application service will use to interact with external systems (things the application "drives"). Example: item_repository.go

package randomnumberscore

import "louder/internal/core/domain"

// Repository could be a port for database interactions
type Repository interface {
	GenerateRandomNumber() domain.RandomNumber
	GenerateDiceRoll(numDice, sides uint) (*domain.RandomDice, error)
}
