package personadapter

import (
	"louder/internal/core/domain"
	"time"
)

// toPersonResponse converts a domain.Person (from the service layer) to a PersonResponse DTO.
func toPersonResponse(p domain.Person) *PersonResponse {
	return &PersonResponse{
		ID:        p.ID().String(),
		FirstName: p.FirstName(),
		LastName:  p.LastName(),
		Email:     p.Email(),
		DOB:       p.DOB().UTC().Format(time.RFC3339Nano),
	}
}
