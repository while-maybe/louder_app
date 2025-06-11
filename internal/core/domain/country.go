package domain

import "strings"

type CountryCode string

type WikiCode string

type Country struct {
	code       CountryCode
	name       string
	currencies []CurrencyCode
	wikidataid WikiCode
}

// NewCountry creates a Country object
func NewCountry(code CountryCode, name string, curr []CurrencyCode, wikidataid WikiCode) Country {
	return Country{
		code:       code,
		name:       name,
		currencies: curr,
		wikidataid: wikidataid,
	}
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

// Currency returns a slice of Currency codes used in the Country
func (c Country) Currency() []CurrencyCode {
	return c.currencies
}
