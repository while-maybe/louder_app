package randomnumberadapter

import "net/http"

func (h *RandomNumberHandler) RegisterRoutes(mux *http.ServeMux) {
	const (
		NewRandomNumberRoute = "/random"
	)

	mux.HandleFunc(http.MethodGet+" "+NewRandomNumberRoute, h.HandleGetRandomNumber)
}

func (h *DiceRollHandler) RegisterRoutes(mux *http.ServeMux) {
	const (
		NewDiceRollRoute = "/diceroll"
	)
	mux.HandleFunc(http.MethodPost+" "+NewDiceRollRoute, h.HandleGetDiceRoll)
}
