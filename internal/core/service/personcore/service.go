// [Driving Adapter (e.g., HTTP Handler)]
//        |
//        v  (calls method on...)
// [Driving Port (e.g., MessageService INTERFACE)]
//        ^
//        | (is implemented by...)
// [Application Service (e.g., messageServiceImpl.go)]  <-- This is the "bridge" or "orchestrator"
//        |
//        | (uses/calls methods on...)
//        v
// [Driven Port (e.g., MessageRepository INTERFACE)]
//        ^
//        | (is implemented by...)
// [Driven Adapter (e.g., PostgresMessageRepository)]

package personcore

import (
	"context"
	"errors"
	"fmt"
	"log"
	dbcommon "louder/internal/adapters/driven/db/dbcommon"
	"louder/internal/core/domain"
	"louder/internal/core/service"
	"louder/pkg/types"

	"github.com/gofrs/uuid/v5"
)

type personServiceImpl struct {
	personRepo PersonRepository
}

func NewPersonService(db PersonRepository) *personServiceImpl {
	return &personServiceImpl{
		personRepo: db,
	}
}

var _ PersonService = (*personServiceImpl)(nil)

// GetAllPersons receives the request from driving port and processes, dispatches the data
// func (ps *personServiceImpl) GetAllPersons(ctx context.Context) ([]domain.Person, error) {
// 	log.Println("Getting all persons from db...")

// 	persons, err := ps.personRepo.GetAll(ctx)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return persons, nil
// }

// CreatePerson implements the business logic for creating a new person
func (ps *personServiceImpl) CreatePerson(ctx context.Context, firstName, lastName, email string) (*domain.Person, error) {
	// some basic validation but more complex logic in domain if needed
	if firstName == "" {
		return nil, fmt.Errorf("%w: first name cannot be empty", service.ErrInvalidPersonData)
	}
	if lastName == "" {
		return nil, fmt.Errorf("%w: last name cannot be empty", service.ErrInvalidPersonData)
	}
	if email == "" { // More robust email validation is usually needed
		return nil, fmt.Errorf("%w: email cannot be empty", service.ErrInvalidPersonData)
	}

	// Business Rules (Example: Check for duplicate email before attempting to create)
	//    This might involve a repository call.
	//    Note: Some prefer to let the DB handle unique constraints and catch the error,  others prefer a proactive check. A proactive check can give a cleaner error.
	//    existingPerson, err := ps.personRepo.FindByEmail(ctx, email) // Assume FindByEmail exists
	//    if err != nil && !errors.Is(err, driven.ErrPersonNotFound) { // driven.ErrPersonNotFound or your repo's equivalent
	//        log.Printf("ERROR CreatePerson - checking for existing email '%s': %v", email, err)
	//        return nil, fmt.Errorf("service error checking email: %w", err) // Generic service error
	//    }
	//    if existingPerson != nil {
	//        return nil, fmt.Errorf("%w: email '%s'", driving.ErrPersonEmailConflict, email)
	//    }

	// create the domain object
	randomDOBTime := types.NewUTCTime(domain.NewRandomDOB())
	newPerson, err := domain.NewPerson(firstName, lastName, email, randomDOBTime)
	if err != nil {
		// This error likely means the data failed domain-level validation within NewPerson
		log.Printf("error CreatePerson - domain.NewPerson: %v", err)
		return nil, fmt.Errorf("%w: %w", service.ErrInvalidPersonData, err) // Wrap domain error
	}

	// save the newly created person to the repo
	// The repository's Save method should handle insert/upsert logic if needed. It returns the *persisted* person, which might have DB-generated fields (like ID if not pre-generated).

	savedPerson, err := ps.personRepo.Save(ctx, newPerson)
	if err != nil {
		log.Printf("error CreatePerson - personRepo.Save (ID: %s): %v", newPerson.ID().String(), err)
		// Check for specific repository errors (e.g., unique constraint violation if not checked before)
		// if errors.Is(err, driven.ErrDuplicateEntry) { // Example repo error
		//  return nil, fmt.Errorf("%w: %w", driving.ErrPersonEmailConflict, err)
		// }
		return nil, fmt.Errorf("failed to save person: %w", err) // Generic persistence error
	}

	// Publish Domain Event (e.g., PersonCreatedEvent) ps.eventPublisher.Publish(ctx, domain.NewPersonCreatedEvent(savedPerson.ID(), ...))

	log.Printf("INFO CreatePerson: Successfully created person ID %s\n", savedPerson.ID().String())
	return savedPerson, nil
}

// GetPersonByID implements the business logic for getting a person by ID from the DB
func (ps *personServiceImpl) GetPersonByID(ctx context.Context, pid domain.PersonID) (*domain.Person, error) {
	// TODO get some proper validation going lazy! (regex?)
	if uuid.UUID(pid).IsNil() {
		return nil, fmt.Errorf("%w: id cannot be nil", service.ErrInvalidPersonData)
	}

	savedPerson, err := ps.personRepo.GetByID(ctx, pid)
	if err != nil {
		log.Printf("warning GetPersonByID - personRepo (ID: %s): %v", pid.String(), err)

		if errors.Is(err, dbcommon.ErrNotFound) {
			return nil, fmt.Errorf("failed to get person: %w", err)
		}
		// For any other repository error, return a generic service/repository interaction error
		return nil, fmt.Errorf("service error: failed to retrieve person with ID %s from repository: %w", pid.String(), err)
	}

	// On success (err == nil), savedPerson variable holds the result. Defensive check: A well-behaved repository should not return (nil, nil).
	if savedPerson == nil {
		log.Printf("error GetPersonByID - repository returned (nil, nil) for ID %s, which is unexpected.", pid.String())
		// Return a generic service error as this indicates an issue with the repository implementation.
		return nil, fmt.Errorf("service error: inconsistent repository response for ID %s", pid.String())
	}

	log.Printf("INFO GetPersonByID: person with ID %s found\n", savedPerson.ID().String())
	return savedPerson, nil
}
