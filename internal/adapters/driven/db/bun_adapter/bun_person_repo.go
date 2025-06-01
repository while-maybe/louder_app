package bunadapter

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	dbcommon "louder/internal/adapters/driven/db/db_common"
	"louder/internal/core/domain"
	drivenports "louder/internal/core/ports/driven"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
)

type BunPersonRepo struct {
	db *bun.DB
}

const (
	ErrBunSavePerson = dbcommon.Error("error Bun save person")
)

// ensure BunPersonRepo implements the drivenports.PersonRepository interface with (won't compile otherwise)
var _ drivenports.PersonRepository = (*BunPersonRepo)(nil)

func NewBunPersonRepo(sqldb *sql.DB) (*BunPersonRepo, error) {
	db := bun.NewDB(sqldb, sqlitedialect.New())
	return &BunPersonRepo{
		db: db,
	}, nil
}

func (bpr *BunPersonRepo) Save(ctx context.Context, person *domain.Person) error {

	// convert from domain.Person to BunPersonModel first
	bunModel := toBunModelPerson(person)

	if bunModel == nil {
		return dbcommon.ErrConvertNilPerson
	}

	result, err := bpr.db.NewInsert().Model(bunModel).On("CONFLICT (id) DO UPDATE").Exec(ctx)
	if err != nil {
		return fmt.Errorf("%w (ID:%s): %w", ErrBunSavePerson, person.ID().String(), err)
	}

	rowsAffected, err := result.RowsAffected()

	switch {
	case err != nil:
		log.Printf("Bun: Warning - couldn't get rows affected for ID: %s, %v", person.ID().String(), err)

	case rowsAffected == 0:
		log.Printf("Bun: Info - 0 rows affected for ID: %s. Identical to record?\n", person.ID().String())

	default:
		log.Printf("Bun: Successfully saved/updated person ID %s. Fetching current state.", person.ID().String())
		_, err = bpr.GetByID(ctx, person.ID().String())
		if err != nil {
			return fmt.Errorf("%w, ID: %s, %w", dbcommon.ErrSavedButNotInDB, person.ID().String(), err)
		}
	}
	return nil
}

func (bpr *BunPersonRepo) GetByID(ctx context.Context, personId string) (*domain.Person, error) {
	if personId == "" {
		return nil, dbcommon.ErrEmptyID
	}

	// validate format
	_, err := domain.PersonIDFromString(personId)
	if err != nil {
		return nil, fmt.Errorf("%w ID: %s, %w", dbcommon.ErrInvalidID, personId, err)
	}

	bunModel := new(BunModelPerson)

	err = bpr.db.NewSelect().Model(bunModel).Where("id = ?", personId).Scan(ctx)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%w for ID '%s'", dbcommon.ErrNotFoundInDB, personId)
		}
		return nil, fmt.Errorf("%w ID: %s, %w", dbcommon.ErrDBQueryFailed, personId, err)
	}

	retrievedPerson, err := bunModel.toDomainPerson()
	if err != nil {
		log.Printf("error BunPersonRepo.GetByID - Failed to convert BunModelPerson to domain.Person for ID '%s'. Model: %+v. Error: %v", personId, bunModel, err)
		return nil, fmt.Errorf("%w, ID: %s, %w", dbcommon.ErrConvertPerson, personId, err)
	}

	return retrievedPerson, nil
}

func (bpr *BunPersonRepo) GetAll(ctx context.Context) ([]domain.Person, error) {
	var dbModels []BunModelPerson

	err := bpr.db.NewSelect().Model(&dbModels).Scan(ctx)

	if err != nil {
		return nil, fmt.Errorf("GetAllPersons: %w: %w", dbcommon.ErrDBQueryFailed, err)
	}

	allPersons := make([]domain.Person, 0, len(dbModels))

	for i := range dbModels {
		domainPerson, err := dbModels[i].toDomainPerson()

		switch {
		case err != nil:
			log.Printf("GetAllPersons: %v (ID: %s): %v. Skipping.", dbcommon.ErrConvertPerson, dbModels[i].ID.String(), err)
		case domainPerson == nil:
			log.Printf("GetAllPersons: %v (ID: %s). Skipping.", dbcommon.ErrNilDomainPerson, dbModels[i].ID.String())
		default:
			allPersons = append(allPersons, *domainPerson)
		}
	}

	return allPersons, nil
}
