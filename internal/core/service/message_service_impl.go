// Implements the driving ports. Contains the core application logic. It depends on driven port interfaces, not concrete implementations. This is key for DI.
package service

import (
	"log"
	"louder/internal/core/domain"
	drivenports "louder/internal/core/ports/driven" // depends on driven port - dependencies
	drivingports "louder/internal/core/ports/driving"
)

type messageServiceImpl struct {
	messageRepo drivenports.MessageWithTimeRepository // Injected dependency
}

// NewMessageService is the constructor for messageServiceImpl
func NewMessageService(db drivenports.MessageWithTimeRepository) drivingports.MessageService { // returns the driving port interface
	return &messageServiceImpl{
		messageRepo: db,
	}
}

func (m *messageServiceImpl) GetMessage() domain.MsgWithTime {
	log.Println("Getting a message from db...")
	return m.messageRepo.GetMessageFromRepo()
}
