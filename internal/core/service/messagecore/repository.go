package messagecore

import "louder/internal/core/domain"

// MessageWithTimeRepository could be a port for database interactions
type MessageWithTimeRepository interface {
	GetMessageFromRepo() domain.MsgWithTime
}
