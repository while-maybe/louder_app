package countrycore

import (
	"context"
	"louder/internal/core/domain"
)

type ExternalCountryProvider interface {
	FetchAllCountries(ctx context.Context) ([]*domain.Country, error)
	GetTotalCountryCountFromAPI(ctx context.Context) (int, error)
}
