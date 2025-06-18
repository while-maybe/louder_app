package domain

import (
	"errors"
	"strings"
)

type CurrencyCode string

type Currency struct {
	code CurrencyCode
	name string
}

func NewCurrency(code CurrencyCode, name string) (*Currency, error) {

	return &Currency{
		code: code,
		name: name,
	}, nil
}

func (c *Currency) Code() CurrencyCode {
	return c.code
}

func (c *Currency) Name() string {
	return c.name
}

func (cc CurrencyCode) String() string {
	return string(cc)
}

func NewCurrencyCode(cc string) (CurrencyCode, error) {
	if cc == "" || len(cc) != 3 {
		return "", errors.New("error creating currency code: must be 3 characters long")
	}
	return CurrencyCode(strings.ToUpper(cc)), nil
}
