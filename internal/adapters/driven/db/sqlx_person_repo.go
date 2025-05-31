package sqlitedbadapter

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"louder/internal/core/domain"
	drivenports "louder/internal/core/ports/driven"

	"github.com/jmoiron/sqlx"
)

const (
	ErrSaveNilPerson   = Error("error cannot save nil person to DB model")
	ErrSqlxSavePerson  = Error("error SQLx save person")
	ErrSavedButNotInDB = Error("error SQLx person saved but can't find in DB")
	ErrNotFoundInDB    = Error("error cannot find this ID in DB")
	ErrEmptyID         = Error("error given ID is empty")
	ErrConvertPerson   = Error("error cannot convert SQLx data to a person")
	ErrDBQueryFailed   = Error("error query has failed")
	ErrInvalidID       = Error("error invalid person ID format")
	ErrNilDomainPerson = Error("error conversion returned nil domain person without error")
	// ErrConvertNilPerson = Error("error convert nil person to DB model")
	// ErrSaveNoRowsAffected = Error("error SQLx can't get rows affected")
)

type SQLxPersonRepo struct {
	db *sqlx.DB
}

// ensure SQLxPersonRepo implements the drivenports.PersonRepository interface with (won't compile otherwise)
var _ drivenports.PersonRepository = (*SQLxPersonRepo)(nil)

func NewSQLxPersonRepo(sqldb *sql.DB) (*SQLxPersonRepo, error) {

	db := sqlx.NewDb(sqldb, "sqlite3")
	return &SQLxPersonRepo{db: db}, nil
}

func (spr *SQLxPersonRepo) Save(ctx context.Context, person *domain.Person) (*domain.Person, error) {
	if person == nil {
		return nil, ErrSaveNilPerson
	}

	// convert from domain.Person to SQLxPersonModel first
	sqlxModel := toSQLxModelPerson(person)

	// yeah... on 2nd look looks redundant
	// if sqlxModel == nil {
	// 	return nil, ErrConvertNilPerson
	// }

	query := `
		INSERT INTO person (id, first_name, last_name, email, dob)
		VALUES (:id, :first_name, :last_name, :email, :dob)
		ON CONFLICT(id) DO UPDATE SET
			first_name = excluded.first_name,
			last_name = excluded.last_name,
			email = excluded.email,
			dob = excluded.dob;`

	result, err := spr.db.NamedExecContext(ctx, query, sqlxModel)
	if err != nil {
		return nil, fmt.Errorf("%w (ID:%s): %w", ErrSqlxSavePerson, person.ID().String(), err)
	}

	rowsAffected, err := result.RowsAffected()
	var retrievedPerson *domain.Person

	switch {
	case err != nil:
		log.Printf("SQLx: Warning - couldn't get rows affected for ID: %s, %v", person.ID().String(), err)

	case rowsAffected == 0:
		log.Printf("SQLx: Info - 0 rows affected for ID: %s. Identical to record?\n", person.ID().String())
		retrievedPerson, err = spr.GetByID(ctx, person.ID().String())

	default:
		log.Printf("SQLx: Successfully saved/updated person ID %s. Fetching current state.", person.ID().String())
		retrievedPerson, err = spr.GetByID(ctx, person.ID().String())
		if err != nil {
			return nil, fmt.Errorf("%w, ID: %s, %w", ErrSavedButNotInDB, person.ID().String(), err)
		}
	}
	return retrievedPerson, nil
}

func (spr *SQLxPersonRepo) GetByID(ctx context.Context, personId string) (*domain.Person, error) {
	if personId == "" {
		return nil, ErrEmptyID
	}

	// validate format
	_, err := domain.PersonIDFromString(personId)
	if err != nil {
		return nil, fmt.Errorf("%w ID: %s, %w", ErrInvalidID, personId, err)
	}

	query := `
		SELECT id, first_name, last_name, email, dob
		FROM person
		WHERE id = ?;`

	var sqlxModel SQLxModelPerson

	err = spr.db.GetContext(ctx, &sqlxModel, query, personId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%w ID: %s", ErrNotFoundInDB, personId)
		}
		return nil, fmt.Errorf("%w ID: %s, %w", ErrDBQueryFailed, personId, err)
	}

	// convert from SQLxPersonModel to domain.Person to
	retrievedPerson, err := toDomainPerson(&sqlxModel)
	if err != nil {
		return nil, fmt.Errorf("%w, ID: %s, %w", ErrConvertPerson, personId, err)
	}

	return retrievedPerson, nil
}

func (spr *SQLxPersonRepo) GetAll(ctx context.Context) ([]domain.Person, error) {
	query := `
		SELECT id, first_name, last_name, email, dob
		FROM person;`

	var dbModels []SQLxModelPerson

	err := spr.db.SelectContext(ctx, &dbModels, query)
	if err != nil {
		return nil, fmt.Errorf("GetAllPersons: %w: %w", ErrDBQueryFailed, err)
	}

	allPersons := make([]domain.Person, 0, len(dbModels))

	for i := range dbModels {
		domainPerson, err := toDomainPerson(&dbModels[i])

		switch {
		case err != nil:
			log.Printf("GetAllPersons: %v (ID: %s): %v. Skipping.", ErrConvertPerson, dbModels[i].ID.String(), err)
		case domainPerson == nil:
			log.Printf("GetAllPersons: %v (ID: %s). Skipping.", ErrNilDomainPerson, dbModels[i].ID.String())
		default:
			allPersons = append(allPersons, *domainPerson)
		}
	}

	return allPersons, nil
}
