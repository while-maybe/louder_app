package domain

import (
	"time"
)

type birthCountry string

type Person struct {
	firstName        string
	lastName         string
	email            string
	dob              time.Time
	pets             []Pet
	birthCountry     Country
	residentCountry  Country
	visitedCountries []Country
}
