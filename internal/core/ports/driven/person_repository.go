package drivenports

import (
	"context"
	"louder/internal/core/domain"
)

type PersonRepository interface {
	GetAllPersonsFromRepo(ctx context.Context) ([]domain.Person, error)
	// GetPersonFromRepo(ctx context.Context, personId string) (domain.Person, error)
	AddPersonToRepo(ctx context.Context, person domain.Person) (domain.Person, error)
	// DelPersonFromRepo(ctx context.Context, personId string) error
	// GetByNameFromRepo(ctx context.Context, name string) ([]domain.Person, error)
	// GetByAgeFromRepo(ctx context.Context, min, max int) ([]domain.Person, error)

	// GetByBirthCountryFromRepo(ctx context.Context, country string) ([]domain.Person, error)
	// GetByResidingCountryFromRepo(ctx context.Context, country string) ([]domain.Person, error)
	// GetByVisitedCountriesFromRepo(ctx context.Context, countries ...string) ([]domain.Person, error)
}
