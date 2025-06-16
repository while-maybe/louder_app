package domain

type CurrencyCode string

type Currency struct {
	code CurrencyCode
	name string
}

func NewCurrency(code CurrencyCode, name string) *Currency {
	return &Currency{
		code: code,
		name: name,
	}
}

func (c *Currency) Code() CurrencyCode {
	return c.code
}

func (c *Currency) Name() string {
	return c.name
}
