package countrycore

import (
	"context"
	"louder/internal/core/domain"
)

type Repository interface {
	Save(ctx context.Context, country *domain.Country) (*domain.Country, error)
	GetByID(ctx context.Context, cc domain.CountryCode) (*domain.Country, error) // ID is the Country's ISO code
	CountAll(ctx context.Context) (int, error)
	GetRandom(ctx context.Context) (*domain.Country, error)

	// GetByName(ctx context.Context, name string) (*domain.Country, error)
	// Search(ctx context.Context, terms string) ([]*domain.Country, error) // get a list of countries when search terms are given, like Google?
}
