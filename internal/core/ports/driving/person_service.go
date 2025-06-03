package drivingports

import (
	"context"
	"louder/internal/core/domain"
)

// PersonService defines the primary use case for Person - What do we do with Person?
type PersonService interface {
	CreatePerson(ctx context.Context, firstName, lastName, email string) (*domain.Person, error)
	GetPersonByID(ctx context.Context, pid domain.PersonID) (*domain.Person, error)
	// GetAll(context.Context) ([]domain.Person, error)
}
