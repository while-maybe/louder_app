package drivingports

import (
	"context"
	"louder/internal/core/domain"
)

// PersonService defines the primary use case for Person - What do we do with Person?
type PersonService interface {
	GetAllPersons(context.Context) ([]domain.Person, error)
}
