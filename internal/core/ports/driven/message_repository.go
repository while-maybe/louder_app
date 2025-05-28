// Interfaces that the core application service will use to interact with external systems (things the application "drives"). Example: item_repository.go

package drivenports

import "louder/internal/core/domain"

// MessageWithTimeRepository could be a port for database interactions
type MessageWithTimeRepository interface {
	GetMessageFromRepo() domain.MsgWithTime
}
