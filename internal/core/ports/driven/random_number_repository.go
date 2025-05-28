// Interfaces that the core application service will use to interact with external systems (things the application "drives"). Example: item_repository.go

package drivenports

import "louder/internal/core/domain"

// RandomNumberRepository could be a port for database interactions
type RandomNumberRepository interface {
	GetRandomNumberFromRepo() domain.RandomNumber
}
