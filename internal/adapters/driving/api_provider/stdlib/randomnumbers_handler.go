package stdlibapiadapter

import (
	"encoding/json"
	"log"
	"louder/internal/core/domain"
	randomnumber "louder/internal/core/service/randomnumbers"

	"net/http"
)

type RandomNumberResponse struct {
	RandomNumber domain.RandomNumber `json:"random_number,omitempty"`
}

// RandomNumber Handler

type RandomNumberHandler struct {
	RandomNumberService randomnumber.Port // inject core service
}

func NewRandomNumberHandler(service randomnumber.Port) *RandomNumberHandler {
	return &RandomNumberHandler{RandomNumberService: service}
}

// HandleGetRandomNumber is an http.HandlerFunc for the /random route
func (h *RandomNumberHandler) HandleGetRandomNumber(w http.ResponseWriter, r *http.Request) {
	log.Println("stdlib API adapter: Got GET request for /random")

	randomNumber := h.RandomNumberService.GetRandomNumber()
	response := RandomNumberResponse{RandomNumber: randomNumber}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Failed to encode random number response %v", err)
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}
