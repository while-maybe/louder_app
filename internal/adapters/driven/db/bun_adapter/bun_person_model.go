package bunadapter

import (
	"fmt"
	dbcommon "louder/internal/adapters/driven/db/db_common"
	"louder/internal/core/domain"
	"time"

	"github.com/uptrace/bun"
)

type BunModelPerson struct {
	bun.BaseModel `bun:"table:person, alias:p"`

	ID        domain.PersonID `bun:"id, pk"`
	FirstName string          `bun:"first_name"`
	LastName  string          `bun:"last_name"`
	Email     string          `bun:"email"`
	DOB       time.Time       `bun:"dob"`
}

// mappers
func toBunModelPerson(p *domain.Person) *BunModelPerson {
	if p == nil {
		return nil
	}

	return &BunModelPerson{
		ID:        p.ID(),
		FirstName: p.FirstName(),
		LastName:  p.LastName(),
		Email:     p.Email(),
		DOB:       p.DOB(),
	}
}

func (m *BunModelPerson) toDomainPerson() (*domain.Person, error) {
	if m == nil {
		return nil, fmt.Errorf("%w", dbcommon.ErrHydrateWithNil)
	}

	return domain.HydratePerson(
		m.ID, m.FirstName, m.LastName, m.Email, m.DOB), nil
}
