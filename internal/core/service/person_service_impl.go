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
	"log"
	"louder/internal/core/domain"
	drivenports "louder/internal/core/ports/driven"
)

type personServiceImpl struct {
	personRepo drivenports.PersonRepository
}

func NewPersonService(db drivenports.PersonRepository) *personServiceImpl {
	return &personServiceImpl{
		personRepo: db,
	}
}

// GetAllPersons receives the request from driving port and processes, dispatches the data
func (ps *personServiceImpl) GetAllPersons(ctx context.Context) ([]domain.Person, error) {
	log.Println("Getting all persons from db...")

	persons, err := ps.personRepo.GetAllPersonsFromRepo(ctx)
	if err != nil {
		return nil, err
	}
	return persons, nil
}
