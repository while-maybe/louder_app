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

// NewPersonID generates a new unique PersonID (UUID v7)
func NewPersonID() (PersonID, error) {
	id, err := uuid.NewV7() // V7 is time ordered
	if err != nil {
		return PersonID(uuid.Nil), err
	}
	return PersonID(id), nil
}

// PersonIDFromString converts a string to a PersonID
func PersonIDFromString(s string) (PersonID, error) {
	id, err := uuid.FromString(s)
	if err != nil {
		return PersonID(uuid.Nil), err
	}
	return PersonID(id), nil
}

// String returns the string representation of the PersonID
func (pid PersonID) String() string {
	return uuid.UUID(pid).String()
}

// isNil checks if the PersonID is a "zero" or nil UUID
func (pid PersonID) isNil() bool {
	return uuid.UUID(pid).IsNil()
}

// NewPerson factory function
func NewPerson(firstName, lastName, email string, dob time.Time /* ... */) (*Person, error) {
	personID, err := NewPersonID()
	if err != nil {
		return nil, err // Propagate error from ID generation
	}

	// ... (rest of your validation and initialization)
	return &Person{
		id:        personID,
		firstName: firstName,
		lastName:  lastName,
		email:     email,
		dob:       dob,
		// ...
	}, nil
}

// ID returns the PersonID object of a Person entity
func (p *Person) ID() PersonID {
	return p.id
}
