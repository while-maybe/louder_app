package stdlibapiadapter

import (
	"encoding/json"
	"fmt"
	"log"
	"louder/internal/core/domain"
	drivingports "louder/internal/core/ports/driving"
	"net/http"
	"time"
)

// PersonHandler handles HTTP requests related to person entities
type PersonHandler struct {
	service drivingports.PersonService // dependency on the Person Service Interface
}

// NewPersonHandler creates a new PersonHandler
func NewPersonHandler(srv drivingports.PersonService) *PersonHandler {
	return &PersonHandler{
		service: srv,
	}
}

// CreatePersonRequest defines the expected JSON payload for creating a person.
// This is a Data Transfer Object (DTO) for the HTTP layer.
type CreatePersonRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

// PersonResponse defines the JSON payload for returning a person just created (inc UUID).
// This is also DTO for the HTTP layer.
type PersonResponse struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	DOB       string `json:"dob"`

	// TODO - implement later
	// Pets             []domain.Pet `json:"pets,omitempty"`
	// BirthCountry     domain.Country `json:"birth_country,omitempty"`
	// ResidentCountry  domain.Country `json:"residence_country,omitempty"`
	// VisitedCountries []domain.Country `json:"visited_countries,omitempty"`
}

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

// HandleCreatePerson handles POST requests to /person
func (h *PersonHandler) HandleCreatePerson(w http.ResponseWriter, r *http.Request) {
	// do not forget to pass the context!
	ctx := r.Context()

	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreatePersonRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("error HandleCreatePerson - decoding request: %v", err)
		http.Error(w, fmt.Sprintf("invalid payload: %v", err), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// should user be able to submit person DOB:
	// Validate and parse DOB from request DTO
	// dob, err := time.Parse("2006-01-02", req.DOB)
	// if err != nil {
	// 	log.Printf("ERROR HandleCreatePerson - parsing DOB '%s': %v", req.DOB, err)
	// 	http.Error(w, fmt.Sprintf("Bad Request: Invalid DOB format (expected YYYY-MM-DD): %v", err), http.StatusBadRequest)
	// 	return
	// }

	// Now we call the service layer with the context and (validated) data
	createdPerson, err := h.service.CreatePerson(ctx, req.FirstName, req.LastName, req.Email)
	if err != nil {
		log.Printf("ERROR HandleCreatePerson - service.CreatePerson: %v", err)
		// IMPORTANT
		// check here for specific error types from the service layer
		// to return more granular HTTP status codes (e.g., 409 Conflict for duplicate email).
		// For now, a generic 500.
		http.Error(w, "failed to create person.", http.StatusInternalServerError)
		return
	}

	// Convert the domain.Person (from service) to PersonResponse DTO for the HTTP response.
	responseDTO := toPersonResponse(*createdPerson)

	w.Header().Set("Content-Type", "application/json")
	// TODO check my path below
	w.Header().Set("Location", fmt.Sprintf("/person/%s", responseDTO.ID))
	w.WriteHeader(http.StatusCreated)

	// encode the responseDTO into JSON
	if err := json.NewEncoder(w).Encode(responseDTO); err != nil {
		log.Printf("error HandleCreatePerson - encoding response: %v", err)
		http.Error(w, "could not encode JSON", http.StatusInternalServerError)
	}

	log.Printf("info: new person created succesfully: %s", createdPerson.ID().String())
}
