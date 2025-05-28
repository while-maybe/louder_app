package stdlibapiadapter

import "net/http"

const (
	NewMessageRoute      = "/message"
	NewRandomNumberRoute = "/random"
)

func NewRouter(handler *MessageHandler) *http.ServeMux {
	r := http.NewServeMux()
	// Go 1.22+ handles the METHOD /path pattern
	r.HandleFunc(http.MethodGet+" "+NewMessageRoute, handler.HandleGetMessage)
	r.HandleFunc(http.MethodGet+" "+NewRandomNumberRoute, handler.HandleGetRandomNumber)
	return r
}
