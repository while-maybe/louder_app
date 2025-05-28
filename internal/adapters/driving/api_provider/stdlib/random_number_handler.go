package stdlibapiadapter

import (
	"encoding/json"
	"log"
	"louder/internal/core/domain"
	drivingports "louder/internal/core/ports/driving"
	"net/http"
)

type RandomNumberResponse struct {
	RandomNumber domain.RandomNumber `json:"random_number,omitempty"`
}

// RandomNumber Handlers

type RandomNumberHandler struct {
	RandomNumberService drivingports.RandomNumberService // inject core service
}

func NewRandomNumberHandler(service drivingports.RandomNumberService) *RandomNumberHandler {
	return &RandomNumberHandler{RandomNumberService: service}
}

// HandleGetRandomNumber is an http.HandlerFunc for the /random route
func (mh *RandomNumberHandler) HandleGetRandomNumber(w http.ResponseWriter, r *http.Request) {
	log.Println("stdlib API adapter: Got GET request for /random")

	randomNumberData := mh.RandomNumberService.GetRandomNumber()
	response := RandomNumberResponse{RandomNumber: randomNumberData}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Failed to encode random number response %v", err)
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}
