package currencycore

import (
	"context"
	"louder/internal/core/domain"
)

type Repository interface {
	Save(ctx context.Context, currency *domain.Currency) (*domain.Currency, error)
	GetByID(ctx context.Context, cc domain.CurrencyCode) (*domain.Currency, error) // ID is the Country's ISO code
	CountAll(ctx context.Context) (int, error)
	GetRandom(ctx context.Context) (*domain.Currency, error)

	// GetByName(ctx context.Context, name string) (*domain.Currency, error)
	// Search(ctx context.Context, terms string) ([]*domain.Currency, error) // get a list of countries when search terms are given, like Google?
}
