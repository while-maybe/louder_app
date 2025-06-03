package domain

import (
	"database/sql/driver"
	"fmt"
	"math/rand"
	"time"

	"github.com/gofrs/uuid/v5"
)

type PersonID uuid.UUID

type Person struct {
	id        PersonID
	firstName string
	lastName  string
	email     string
	dob       time.Time

	// TODO - implement later
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

// Bytes returns a slice of bytes representation of the PersonID
func (pid PersonID) Bytes() []byte {
	return uuid.UUID(pid).Bytes()
}

// Value implements the driver.Valuer interface for database storage. This tells the SQL driver how to store PersonID in the database.
func (pid PersonID) Value() (driver.Value, error) {
	return uuid.UUID(pid).Bytes(), nil
	// return uuid.UUID(pid).String(), nil // Store as string
}

// isNil checks if the PersonID is a "zero" or nil UUID
func (pid PersonID) isNil() bool {
	return uuid.UUID(pid).IsNil()
}

// Scan implements the sql.Scanner interface for reading from the database. This tells the SQL driver how to convert the database value back to PersonID.
func (pid *PersonID) Scan(value any) error {
	if value == nil {
		*pid = PersonID(uuid.Nil)
		return nil
	}
	var u uuid.UUID
	switch v := value.(type) {
	case string:
		var err error
		u, err = uuid.FromString(v)
		if err != nil {
			return fmt.Errorf("Person ID Scan: failed to parse UUID from string %s: %w", v, err)
		}
	case []byte:
		var err error
		u, err = uuid.FromBytes(v)
		if err != nil {
			return fmt.Errorf("PersonID Scan: failed to parse UUID from bytes %s: %w", string(v), err)
		}
	default:
		return fmt.Errorf("PersonID Scan: unsupported type %T for PersonID", value)
	}

	*pid = PersonID(u)
	return nil
}

// NewRandomDOB return a moment in time in the past maxAge
func NewRandomDOB() time.Time {
	maxAge := 100
	maxApproxDays := int(float64(maxAge) * 365.2425) // acceptable precision
	randDays := rand.Intn(maxApproxDays) + 1
	randomTimeOfDay := time.Duration(rand.Intn(24*3600)) * time.Second
	// .AddDate on its own, returns a time.Time with 00h00m00s so we add a random time of day
	return time.Now().AddDate(0, 0, -randDays).Add(randomTimeOfDay).UTC()
}

// NewPerson factory function
func NewPerson(firstName, lastName, email string, dob time.Time) (*Person, error) {
	personID, err := NewPersonID()
	if err != nil {
		return nil, err
	}

	return &Person{
		id:        personID,
		firstName: firstName,
		lastName:  lastName,
		email:     email,
		dob:       dob,
	}, nil
}

// HydratePerson accepts data from repository and creates a new Person object from it
func HydratePerson(id PersonID, firstName, lastName, email string, dob time.Time) *Person {
	return &Person{
		id:        id,
		firstName: firstName,
		lastName:  lastName,
		email:     email,
		dob:       dob,
	}
}

// ID returns the PersonID object of a Person entity
func (p *Person) ID() PersonID {
	return p.id
}

func (p *Person) FirstName() string {
	return p.firstName
}

func (p *Person) LastName() string {
	return p.lastName
}

func (p *Person) Email() string {
	return p.email
}

func (p *Person) DOB() time.Time {
	return p.dob.UTC()
}
