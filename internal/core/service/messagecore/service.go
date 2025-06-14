// Implements the driving ports. Contains the core application logic. It depends on driven port interfaces, not concrete implementations. This is key for DI.
package messagecore

import (
	"log"
	"louder/internal/core/domain"
)

type messageServiceImpl struct {
	messageRepo MessageWithTimeRepository // Injected dependency
}

// NewMessageService is the constructor for messageServiceImpl
func NewMessageService(db MessageWithTimeRepository) MessageService { // returns the driving port interface
	return &messageServiceImpl{
		messageRepo: db,
	}
}

func (m *messageServiceImpl) GetMessage() domain.MsgWithTime {
	log.Println("Getting a message from db...")
	return m.messageRepo.GetMessageFromRepo()
}
