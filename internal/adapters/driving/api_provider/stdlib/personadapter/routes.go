package personadapter

import "net/http"

func (h *PersonHandler) RegisterRoutes(mux *http.ServeMux) {
	const (
		GetPersonRoute = "/person/"
		NewPersonRoute = "/person"
	)
	mux.HandleFunc(http.MethodGet+" "+GetPersonRoute, h.HandleGetPersonByID)
	mux.HandleFunc(http.MethodPost+" "+NewPersonRoute, h.HandleCreatePerson)
}
