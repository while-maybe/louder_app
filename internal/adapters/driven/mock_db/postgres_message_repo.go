// Concrete implementations of the ItemRepository (driven port). Example: postgres_item_repo.go

package mockdb

import (
	"log"
	"louder/internal/core/domain"
	"time"
)

type MockDBMessageRepository struct {
	// db *sql.DB
	mockDB *domain.MsgWithTime
}

func NewMockDBMessageRepository(dataSourceName string) *MockDBMessageRepository {
	fakeMsgObj := newMsgWithTime(dataSourceName)

	log.Println("Talking to mockDB message repo:", dataSourceName)

	return &MockDBMessageRepository{
		mockDB: fakeMsgObj,
	}
}

func newMsgWithTime(msg string) *domain.MsgWithTime {
	var msgContents string
	if msg == "" {
		msgContents = "no message"
	} else {
		msgContents = msg
	}

	return &domain.MsgWithTime{
		CurrentLocalTime: time.Now(),
		Message:          msgContents,
	}
}

// func newRandomNumber() RandomNumber {
// 	return RandomNumber(rand.Uint32())
// }
