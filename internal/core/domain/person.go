package domain

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

type PersonID uuid.UUID

type Person struct {
	id               PersonID
	firstName        string
	lastName         string
	email            string
	dob              time.Time
	pets             []Pet
	birthCountry     Country
	residentCountry  Country
	visitedCountries []Country
}
