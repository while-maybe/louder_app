// Concrete implementations of the ItemRepository (driven port). Example: postgres_item_repo.go

package mockdbadapter

import (
	"log"
	"louder/internal/core/domain"
	"math/rand/v2"
)

type MockDBRandRepository struct {
	// db *sql.DB

	name   string
	mockDB *domain.RandomNumber
}

func NewMockDBRandRepository(startMsg string) *MockDBRandRepository {
	// this should be a db

	fakeRandObj := newRandomNumber()

	log.Println("Talking to mockDB random number repo:", startMsg)

	return &MockDBRandRepository{
		name:   "MockDBRandRepository",
		mockDB: &fakeRandObj,
	}
}

func (r *MockDBRandRepository) GetRandomNumberFromRepo() domain.RandomNumber {
	return newRandomNumber()
}

func newRandomNumber() domain.RandomNumber {
	return domain.RandomNumber(rand.Uint32())
}
