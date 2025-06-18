package sqlxadapter

import (
	"fmt"
	dbcommon "louder/internal/adapters/driven/db/dbcommon"
	"louder/internal/core/domain"
)

type CurrencyModel struct {
	Code domain.CurrencyCode `db:"code"`
	Name string              `db:"name"`
}

// toModelCurrency takes a Currency domain entity and returns its equivalent SQLx model
func toModelCurrency(c *domain.Currency) *CurrencyModel {
	if c == nil {
		return nil
	}

	return &CurrencyModel{
		Code: c.Code(),
		Name: c.Name(),
	}
}

// toDomainCurrency takes a SQLx currency model and returns its equivalent domain entity
func (m *CurrencyModel) toDomainCurrency() (*domain.Currency, error) {
	if m == nil {
		return nil, dbcommon.ErrConvertCurrency
	}

	newCurrency, err := domain.NewCurrency(m.Code, m.Name)
	if err != nil {
		return nil, fmt.Errorf("cannot create this domain currency: %w", err)
	}

	return newCurrency, nil
}
