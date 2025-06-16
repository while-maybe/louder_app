package sqlxadapter

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	dbcommon "louder/internal/adapters/driven/db/dbcommon"
	"louder/internal/core/domain"
	"louder/internal/core/service/personcore"

	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
)

var (
	ErrSqlxSavePerson = errors.New("error SQLx save person")

	// ErrSaveNoRowsAffected = errors.New("error SQLx can't get rows affected")
)

type PersonRepo struct {
	db *sqlx.DB
}

// ensure SQLxPersonRepo implements the drivenports.PersonRepository interface with (won't compile otherwise)
var _ personcore.PersonRepository = (*PersonRepo)(nil)

func NewSQLxPersonRepo(sqldb *sql.DB) (*PersonRepo, error) {

	db := sqlx.NewDb(sqldb, "sqlite3")
	return &PersonRepo{db: db}, nil
}

func (spr *PersonRepo) Save(ctx context.Context, person *domain.Person) (*domain.Person, error) {

	// convert from domain.Person to SQLxPersonModel first
	sqlxModel := toSQLxModelPerson(person)

	if sqlxModel == nil {
		return nil, dbcommon.ErrConvertNilPerson
	}

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

	var createdPerson *domain.Person
	switch {
	case err != nil:
		log.Printf("SQLx: Warning - couldn't get rows affected for ID: %s, %v", person.ID().String(), err)

	case rowsAffected == 0:
		log.Printf("SQLx: Info - 0 rows affected for ID: %s. Existing record?\n", person.ID().String())

	default:
		log.Printf("SQLx: Successfully saved/updated person ID %s. Fetching current state.", person.ID().String())
		createdPerson, err = spr.GetByID(ctx, person.ID())
		if err != nil {
			return nil, fmt.Errorf("%w, ID: %s, %w", dbcommon.ErrSavedButNotInDB, person.ID().String(), err)
		}
	}
	return createdPerson, nil
}

func (spr *PersonRepo) GetByID(ctx context.Context, pid domain.PersonID) (*domain.Person, error) {

	if uuid.UUID(pid).IsNil() {
		return nil, dbcommon.ErrEmptyID
	}

	// validate format
	// _, err := domain.PersonIDFromString(personId)
	// if err != nil {
	// 	return nil, fmt.Errorf("%w ID: %s, %w", dbcommon.ErrInvalidID, personId, err)
	// }

	query := `
		SELECT id, first_name, last_name, email, dob
		FROM person
		WHERE id = ?;`

	var sqlxModel SQLxModelPerson

	err := spr.db.GetContext(ctx, &sqlxModel, query, pid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%w ID: %s", dbcommon.ErrNotFound, pid.String())
		}
		return nil, fmt.Errorf("%w ID: %s, %w", dbcommon.ErrDBQueryFailed, pid.String(), err)
	}

	// convert from SQLxPersonModel to domain.Person to
	retrievedPerson, err := sqlxModel.toDomainPerson()
	if err != nil {
		return nil, fmt.Errorf("%w, ID: %s, %w", dbcommon.ErrConvertToPerson, pid.String(), err)
	}

	return retrievedPerson, nil
}

func (spr *PersonRepo) GetAll(ctx context.Context) ([]domain.Person, error) {
	query := `
		SELECT id, first_name, last_name, email, dob
		FROM person;`

	var dbModels []SQLxModelPerson

	err := spr.db.SelectContext(ctx, &dbModels, query)
	if err != nil {
		return nil, fmt.Errorf("GetAllPersons: %w: %w", dbcommon.ErrDBQueryFailed, err)
	}

	allPersons := make([]domain.Person, 0, len(dbModels))

	for i := range dbModels {
		domainPerson, err := dbModels[i].toDomainPerson()

		switch {
		case err != nil:
			log.Printf("GetAllPersons: %v (ID: %s): %v. Skipping.", dbcommon.ErrConvertToPerson, dbModels[i].ID.String(), err)
		case domainPerson == nil:
			log.Printf("GetAllPersons: %v (ID: %s). Skipping.", dbcommon.ErrNilDomainPerson, dbModels[i].ID.String())
		default:
			allPersons = append(allPersons, *domainPerson)
		}
	}

	return allPersons, nil
}
