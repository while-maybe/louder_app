// Interfaces that the core application service will use to interact with external systems (things the application "drives"). Example: item_repository.go

package drivenports

import "louder/internal/core/domain"

// MessageRepository could be a port for database interactions
type DataRepository interface {
	GetDBMessage() domain.MsgWithTime
	GetDBRandomNumber() domain.RandomNumber
}
