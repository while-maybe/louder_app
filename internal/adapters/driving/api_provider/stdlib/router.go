package stdlibapiadapter

import "net/http"

const (
	NewMessageRoute      = "/message"
	NewRandomNumberRoute = "/random"
)

func NewRouter(messageHandler *MessageHandler, randomNumberHandler *RandomNumberHandler) *http.ServeMux {
	r := http.NewServeMux()
	// Go 1.22+ handles the METHOD /path pattern
	r.HandleFunc(http.MethodGet+" "+NewMessageRoute, messageHandler.HandleGetMessage)
	r.HandleFunc(http.MethodGet+" "+NewRandomNumberRoute, randomNumberHandler.HandleGetRandomNumber)
	return r
}
