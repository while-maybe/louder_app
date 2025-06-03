package stdlibapiadapter

import "net/http"

const (
	NewMessageRoute      = "/message"
	NewRandomNumberRoute = "/random"
	NewPersonRoute       = "/person"
	GetPersonRoute       = "/person/"
)

func NewRouter(
	messageHandler *MessageHandler,
	randomNumberHandler *RandomNumberHandler,
	personHandler *PersonHandler,
) *http.ServeMux {
	r := http.NewServeMux()
	// Go 1.22+ handles the METHOD /path pattern
	// Message, RandomNumber
	r.HandleFunc(http.MethodGet+" "+NewMessageRoute, messageHandler.HandleGetMessage)
	r.HandleFunc(http.MethodGet+" "+NewRandomNumberRoute, randomNumberHandler.HandleGetRandomNumber)
	// Person
	r.HandleFunc(http.MethodPost+" "+NewPersonRoute, personHandler.HandleCreatePerson)
	// r.HandleFunc(http.MethodGet+" "+GetPersonRoute, personHandler.HandleGetPerson)
	// TODO research timout on individual routes
	return r
}
