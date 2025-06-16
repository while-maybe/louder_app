package personadapter

import (
	"encoding/json"
	"errors"
	"log"
	dbcommon "louder/internal/adapters/driven/db/dbcommon"
	"louder/internal/core/domain"
	"net/http"
	"strings"

	"github.com/gofrs/uuid/v5"
)

// HandleGetPersonByID handles get requests to /person
func (h *PersonHandler) HandleGetPersonByID(w http.ResponseWriter, r *http.Request) {
	// do not forget to pass the context!
	ctx := r.Context()

	// check the method
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract the id from the path in the request
	path := strings.Trim(r.URL.Path, "/")
	parts := strings.Split(path, "/")

	if len(parts) < 2 || parts[0] != "person" {
		log.Printf("warning HandleGetPersonByID - Invalid path format: %s", r.URL.Path)
		http.Error(w, "Bad Request: Invalid URL path. Expected /person/{id}", http.StatusBadRequest)
		return
	}

	idStr := parts[1]

	// validate the extracted string

	// check for empty string
	if idStr == "" {
		http.Error(w, "id cannot be blank", http.StatusBadRequest)
		return
	}

	// convert the string we extracted into a uuid (also check for a valid uuid)
	personUUID, err := uuid.FromString(idStr)
	if err != nil {
		http.Error(w, "Invalid id format", http.StatusBadRequest)
		log.Printf("warning HandleGetPersonByID - Invalid UUID format '%s': %v", idStr, err)
		return
	}

	// check if uuid is V7
	if personUUID.Version() != 7 {
		http.Error(w, "Invalid id version", http.StatusBadRequest)
		return
	}
	// log.Println("\n----->", personUUID.Version())

	// convert the uuid into a domain.PersonID
	personID := domain.PersonID(personUUID)

	// get this personID from service layer
	retrievedPerson, err := h.service.GetPersonByID(ctx, personID)
	if err != nil {
		log.Printf("error HandleGetPersonByID - service.GetPersonByID for ID %s: %v", idStr, err)

		if errors.Is(err, dbcommon.ErrNotFound) { // use error from core/driven
			http.Error(w, "Not found: person with the specified ID does not exist", http.StatusNotFound)

		} else { // this for all other possible error that are not NotFound
			http.Error(w, "Internal Server Error: Failed to retrieve person.", http.StatusInternalServerError)
		}
		return
	}

	// convert the domain.Person returned by the DB to a PersonResponseDTO
	responseDTO := toPersonResponse(*retrievedPerson)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(responseDTO); err != nil {
		log.Printf("error HandleGetPersonByID - encoding response: %v", err)
	}

	log.Printf("Info HandleGetPersonByID - Successfully retrieved person ID: %s", retrievedPerson.ID().String())
}
