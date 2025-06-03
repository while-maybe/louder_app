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

package service

import (
	"context"
	"fmt"
	"log"
	"louder/internal/core/domain"
	drivenports "louder/internal/core/ports/driven"
	drivingports "louder/internal/core/ports/driving"
)

type personServiceImpl struct {
	personRepo drivenports.PersonRepository
}

func NewPersonService(db drivenports.PersonRepository) *personServiceImpl {
	return &personServiceImpl{
		personRepo: db,
	}
}

var _ drivingports.PersonService = (*personServiceImpl)(nil)

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
		return nil, fmt.Errorf("%w: first name cannot be empty", ErrInvalidPersonData)
	}
	if lastName == "" {
		return nil, fmt.Errorf("%w: last name cannot be empty", ErrInvalidPersonData)
	}
	if email == "" { // More robust email validation is usually needed
		return nil, fmt.Errorf("%w: email cannot be empty", ErrInvalidPersonData)
	}

	// Business Rules (Example: Check for duplicate email before attempting to create)
	//    This might involve a repository call.
	//    Note: Some prefer to let the DB handle unique constraints and catch the error,
	//    others prefer a proactive check. A proactive check can give a cleaner error.
	//    existingPerson, err := ps.personRepo.FindByEmail(ctx, email) // Assume FindByEmail exists
	//    if err != nil && !errors.Is(err, driven.ErrPersonNotFound) { // driven.ErrPersonNotFound or your repo's equivalent
	//        log.Printf("ERROR CreatePerson - checking for existing email '%s': %v", email, err)
	//        return nil, fmt.Errorf("service error checking email: %w", err) // Generic service error
	//    }
	//    if existingPerson != nil {
	//        return nil, fmt.Errorf("%w: email '%s'", driving.ErrPersonEmailConflict, email)
	//    }

	// create the domain object
	newPersonDOB := domain.NewRandomDOB()
	newPerson, err := domain.NewPerson(firstName, lastName, email, newPersonDOB)
	if err != nil {
		// This error likely means the data failed domain-level validation within NewPerson
		log.Printf("error CreatePerson - domain.NewPerson: %v", err)
		return nil, fmt.Errorf("%w: %w", ErrInvalidPersonData, err) // Wrap domain error
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

	// 5. [Optional] Publish Domain Event (e.g., PersonCreatedEvent) ps.eventPublisher.Publish(ctx, domain.NewPersonCreatedEvent(savedPerson.ID(), ...))

	log.Printf("INFO CreatePerson: Successfully created person ID %s", savedPerson.ID().String())
	return savedPerson, nil
}
