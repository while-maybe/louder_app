package stdlibapiadapter

import "net/http"

const (
	NewRandomNumberRoute = "/random"
	NewDiceRollRoute     = "/diceroll"

	NewMessageRoute = "/message"
	NewPersonRoute  = "/person"
	GetPersonRoute  = "/person/"
)

func NewRouter(
	randomNumberHandler *RandomNumberHandler,
	diceRollHandler *DiceRollHandler,
	messageHandler *MessageHandler,
	personHandler *PersonHandler,
) *http.ServeMux {
	r := http.NewServeMux()
	// Go 1.22+ handles the METHOD /path pattern
	// RandomNumber, DiceRoll, Message
	r.HandleFunc(http.MethodGet+" "+NewRandomNumberRoute, randomNumberHandler.HandleGetRandomNumber)
	r.HandleFunc(http.MethodPost+" "+NewDiceRollRoute, diceRollHandler.HandleGetDiceRoll)
	r.HandleFunc(http.MethodGet+" "+NewMessageRoute, messageHandler.HandleGetMessage)
	// Person
	r.HandleFunc(http.MethodPost+" "+NewPersonRoute, personHandler.HandleCreatePerson)
	r.HandleFunc(http.MethodGet+" "+GetPersonRoute, personHandler.HandleGetPersonByID)
	// TODO research timout on individual routes
	return r
}
