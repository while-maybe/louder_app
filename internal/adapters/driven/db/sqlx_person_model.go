package sqlitedbadapter

import (
	"fmt"
	"louder/internal/core/domain"
	"time"
)

const (
	ErrConvertIDfromDB = Error("error converting person ID from DB")
	ErrHydrateWithNil  = Error("error attempted to hydrate person without data")
)

// SQLxPersonModel is the data structure used for interacting with the 'person' table using SQLx
// 'db' tags are added to the struct fields for SQLx mapping
type SQLxModelPerson struct {
	ID        domain.PersonID `db:"id"`
	FirstName string          `db:"first_name"`
	LastName  string          `db:"last_name"`
	Email     string          `db:"email"`
	DOB       time.Time       `db:"dob"`
}

// mappers
func toSQLxModelPerson(p *domain.Person) *SQLxModelPerson {
	if p == nil {
		return nil
	}

	return &SQLxModelPerson{
		ID:        p.ID(),
		FirstName: p.FirstName(),
		LastName:  p.LastName(),
		Email:     p.Email(),
		DOB:       p.DOB(),
	}
}

func toDomainPerson(smp *SQLxModelPerson) (*domain.Person, error) {
	if smp == nil {
		return nil, fmt.Errorf("%w", ErrHydrateWithNil)
	}

	// this is actually reduntant at this stage as smp.ID is already the correct type?

	// personID, err := domain.PersonIDFromString(smp.ID.String())
	// if err != nil {
	// 	return nil, fmt.Errorf("%w: %w", ErrConvertIDfromDB, err)
	// }

	return domain.HydratePerson(
		smp.ID, smp.FirstName, smp.LastName, smp.Email, smp.DOB), nil
}
