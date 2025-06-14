package stdlibapiadapter

import (
	"net/http"
)

// NewRouter now takes a slice of Resource
func NewRouter(resources ...Resource) *http.ServeMux {
	mux := http.NewServeMux()

	for _, r := range resources {
		r.RegisterRoutes(mux)
	}

	return mux
}
