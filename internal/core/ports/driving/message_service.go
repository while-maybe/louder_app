package drivingports

import "louder/internal/core/domain"

// MessageService defines the primary use case for Messages - What do we do with Messages?
type MessageService interface {
	GetMessage() domain.MsgWithTime
}

type RandomNumberService interface {
	GetRandomNumber() domain.RandomNumber
}
