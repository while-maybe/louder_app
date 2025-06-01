package drivenports

import (
	"context"
	"louder/internal/core/domain"
)

type PersonRepository interface {
	GetAll(ctx context.Context) ([]domain.Person, error)
	GetByID(ctx context.Context, personId string) (*domain.Person, error)
	Save(ctx context.Context, person *domain.Person) error
	// DelPersonFromRepo(ctx context.Context, personId string) error
	// GetByNameFromRepo(ctx context.Context, name string) ([]domain.Person, error)
	// GetByAgeFromRepo(ctx context.Context, min, max int) ([]domain.Person, error)

	// GetByBirthCountryFromRepo(ctx context.Context, country string) ([]domain.Person, error)
	// GetByResidingCountryFromRepo(ctx context.Context, country string) ([]domain.Person, error)
	// GetByVisitedCountriesFromRepo(ctx context.Context, countries ...string) ([]domain.Person, error)
}
