// Concrete implementations of the ItemRepository (driven port). Example: postgres_item_repo.go

package mockdb

import (
	"log"
	"louder/internal/core/domain"
	"math/rand/v2"
)

type MockDBRandRepository struct {
	// db *sql.DB

	mockDB *domain.RandomNumber
}

func NewMockDBRandRepository() *MockDBRandRepository {
	fakeMsgObj := newRandomNumber()

	log.Println("Talking to mockDB message repo:")

	return &MockDBRandRepository{
		mockDB: &fakeMsgObj,
	}
}

func newRandomNumber() domain.RandomNumber {
	return domain.RandomNumber(rand.Uint32())
}
