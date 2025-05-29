package drivenports

import (
	"context"
	"louder/internal/core/domain"
)

type PeopleRepository interface {
	ListAll(ctx context.Context) ([]domain.Person, error)
	ShowPerson(ctx context.Context, personId string) (domain.Person, error)
	AddPerson(ctx context.Context, person domain.Person) (domain.Person, error)
	DelPerson(ctx context.Context, personId string) error
	FindByName(ctx context.Context, name string) ([]domain.Person, error)
	FindByAge(ctx context.Context, min, max int) ([]domain.Person, error)

	FindByBirthCountry(ctx context.Context, country string) ([]domain.Person, error)
	FindByResidingCountry(ctx context.Context, country string) ([]domain.Person, error)
	FindByVisitedCountries(ctx context.Context, countries ...string) ([]domain.Person, error)
}
