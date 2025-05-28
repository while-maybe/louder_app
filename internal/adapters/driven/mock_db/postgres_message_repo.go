// Concrete implementations of the ItemRepository (driven port). Example: postgres_item_repo.go

package mockdbadapter

import (
	"fmt"
	"log"
	"louder/internal/core/domain"
	"time"
)

type MockDBMessageRepository struct {
	// db *sql.DB

	name   string
	mockDB *domain.MsgWithTime
}

func NewMockDBMessageRepository(startMessage string) *MockDBMessageRepository {
	fakeMsgObj := newMsgWithTime(startMessage)

	log.Println("Talking to mockDB message repo:", startMessage)

	return &MockDBMessageRepository{
		mockDB: fakeMsgObj,
		name:   "MockDBMessageRepository",
	}
}

func (r *MockDBMessageRepository) GetMessageFromRepo() domain.MsgWithTime {
	output := fmt.Sprintf("[%s] -> %s", r.name, r.mockDB.Content)
	return domain.MsgWithTime{
		Content:   output,
		Timestamp: time.Now(),
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
		Timestamp: time.Now(),
		Content:   msgContents,
	}
}
