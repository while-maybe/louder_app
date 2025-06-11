package randomnumbers

import "louder/internal/core/domain"

// What services does random number provide to others? what can we ask out of it?
type Port interface {
	GetRandomNumber() domain.RandomNumber
	RollDice(numDice, numSides uint) (*domain.RandomDice, error)
}
