package domain

import (
	"fmt"
	"strings"
)

type CountryCode string

type WikiCode string

type Country struct {
	code       CountryCode
	name       string
	currencies []Currency
	wikidataid WikiCode
}

// NewCountry creates a Country object
func NewCountry(code CountryCode, name string, currs []Currency, wikidataid WikiCode) (*Country, error) {
	if code.String() == "" {
		return nil, fmt.Errorf("error country code cannot be empty")
	}

	if name == "" {
		return nil, fmt.Errorf("error country name cannot be empty")
	}

	var currenciesCopy []Currency
	if currs != nil {
		currenciesCopy = make([]Currency, 0, len(currs))
		copy(currenciesCopy, currs)
	} else {
		currenciesCopy = []Currency{}
	}

	return &Country{
		code:       code,
		name:       name,
		currencies: currenciesCopy,
		wikidataid: wikidataid,
	}, nil
}

// WikiId returns a string with the wikidataid value
func (c Country) WikiId() WikiCode {
	return c.wikidataid
}

// Name returns the name of the Country as a string
func (c Country) Name() string {
	return strings.ToTitle(c.name)
}

// Code returns the 2 digit Country code as a string
func (c Country) Code() CountryCode {
	return c.code
}

// Currencies returns a slice of Currencies used in the Country
func (c Country) Currencies() []Currency {
	return c.currencies
}

// String returns the string contained in a CountryCode
func (cc CountryCode) String() string {
	return string(cc)
}
