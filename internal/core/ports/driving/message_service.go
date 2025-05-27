package drivingports

// MessageService defines the primary use case for Messages
type MessageService interface {
	FetchMessage() string
	FetchRandomNumber() uint
}
