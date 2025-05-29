package drivingports

import "louder/internal/core/domain"

// RandomNumberService defines the primary use case for Random Numbers - What do we do with Messages?
type RandomNumberService interface {
	GetRandomNumber() domain.RandomNumber
}
