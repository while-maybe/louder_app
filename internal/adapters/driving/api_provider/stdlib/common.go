package stdlibapiadapter

import (
	"encoding/json"
	"log"
	"net/http"
)

func RespondWithJSON(w http.ResponseWriter, code int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if payload != nil {
		if err := json.NewEncoder(w).Encode(payload); err != nil {
			log.Printf("Failed to encode JSON response: %v", err)
		}
	}
}

func RespondWithError(w http.ResponseWriter, code int, message string) {
	// We wrap the single message in our standard ErrorResponse struct
	RespondWithJSON(w, code, ErrorResponse{ErrorMsgs: []string{message}})
}
