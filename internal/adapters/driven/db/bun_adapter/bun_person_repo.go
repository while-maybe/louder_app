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

	"github.com/gofrs/uuid/v5"
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

func (bpr *BunPersonRepo) Save(ctx context.Context, person *domain.Person) (*domain.Person, error) {

	// convert from domain.Person to BunPersonModel first
	bunModel := toBunModelPerson(person)

	fmt.Println(bunModel)
	if bunModel == nil {
		return nil, dbcommon.ErrConvertNilPerson
	}
	result, err := bpr.db.NewInsert().Model(bunModel).On("CONFLICT (id) DO UPDATE").Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w (ID:%s): %w", ErrBunSavePerson, person.ID().String(), err)
	}

	rowsAffected, err := result.RowsAffected()
	var createdPerson *domain.Person

	switch {
	case err != nil:
		log.Printf("Bun: Warning - couldn't get rows affected for ID: %s, %v", person.ID().String(), err)

	case rowsAffected == 0:
		log.Printf("Bun: Info - 0 rows affected for ID: %s. Identical to record?\n", person.ID().String())

	default:
		log.Printf("Bun: Successfully saved/updated person ID %s. Fetching current state.", person.ID().String())
		createdPerson, err = bpr.GetByID(ctx, person.ID())
		if err != nil {
			return nil, fmt.Errorf("%w, ID: %s, %w", dbcommon.ErrSavedButNotInDB, person.ID().String(), err)
		}
	}
	return createdPerson, nil
}

func (bpr *BunPersonRepo) GetByID(ctx context.Context, pid domain.PersonID) (*domain.Person, error) {
	// check if pid is empty
	if uuid.UUID(pid).IsNil() {
		return nil, dbcommon.ErrEmptyID
	}

	// check if pid is properly formatted
	if _, err := uuid.FromBytes(pid.Bytes()); err != nil {
		return nil, dbcommon.ErrInvalidID
	}

	// check if pid is the right uuid version (version 7)
	if uuid.UUID(pid).Version() != 7 {
		return nil, dbcommon.ErrInvalidID
	}

	bunModel := new(BunModelPerson)

	err := bpr.db.NewSelect().Model(bunModel).Where("id = ?", pid).Scan(ctx)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%w for ID '%s'", dbcommon.ErrNotFound, pid.String())
		}
		return nil, fmt.Errorf("%w ID: %s, %w", dbcommon.ErrDBQueryFailed, pid.String(), err)
	}

	retrievedPerson, err := bunModel.toDomainPerson()
	if err != nil {
		log.Printf("error BunPersonRepo.GetByID - Failed to convert BunModelPerson to domain.Person for ID '%s'. Model: %+v. Error: %v", pid.String(), bunModel, err)
		return nil, fmt.Errorf("%w, ID: %s, %w", dbcommon.ErrConvertPerson, pid.String(), err)
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
