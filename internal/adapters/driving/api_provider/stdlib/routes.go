package stdlibapiadapter

import "net/http"

// Resource is an interface that a feature handler (like personadapter) must implement so the main router can register its routes.
type Resource interface {
	RegisterRoutes(mux *http.ServeMux)
}
