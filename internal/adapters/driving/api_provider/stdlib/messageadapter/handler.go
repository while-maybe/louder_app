package messageadapter

import (
	"encoding/json"
	"log"
	"louder/internal/core/domain"
	drivingports "louder/internal/core/ports/driving"
	"net/http"
)

type MessageResponse struct {
	Message domain.MsgWithTime `json:"message_data,omitempty"`
}

// Message handler

type MessageHandler struct {
	MessageService drivingports.MessageService // injected core service
}

func (h *MessageHandler) RegisterRoutes(mux *http.ServeMux) {
	const (
		NewMessageRoute = "/message"
	)
	mux.HandleFunc(http.MethodGet+" "+NewMessageRoute, h.HandleGetMessage)
}

func NewMessageHandler(service drivingports.MessageService) *MessageHandler {
	return &MessageHandler{MessageService: service}
}

// HandleGetMessage is an http.HandlerFunc for the /message route
func (mh *MessageHandler) HandleGetMessage(w http.ResponseWriter, r *http.Request) {
	log.Println("stdlib API adapter: Got GET request for /message")

	msgData := mh.MessageService.GetMessage()
	response := MessageResponse{Message: msgData}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Failed to encode message response %v", err)
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}
