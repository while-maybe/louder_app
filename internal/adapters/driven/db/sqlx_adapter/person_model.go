package sqlxadapter

import (
	"errors"
	"fmt"
	dbcommon "louder/internal/adapters/driven/db/dbcommon"
	"louder/internal/core/domain"
	"louder/pkg/types"
)

var (
	ErrConvertIDfromDB = errors.New("error converting person ID from DB")
)

// SQLxPersonModel is the data structure used for interacting with the 'person' table using SQLx
// 'db' tags are added to the struct fields for SQLx mapping
type SQLxModelPerson struct {
	ID        domain.PersonID `db:"id"`
	FirstName string          `db:"first_name"`
	LastName  string          `db:"last_name"`
	Email     string          `db:"email"`
	DOB       types.UTCTime   `db:"dob"`
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

func (m *SQLxModelPerson) toDomainPerson() (*domain.Person, error) {
	if m == nil {
		return nil, fmt.Errorf("%w", dbcommon.ErrHydrateWithNil)
	}

	return domain.HydratePerson(
		m.ID, m.FirstName, m.LastName, m.Email, m.DOB), nil
}
