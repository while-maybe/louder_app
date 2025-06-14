package stdlibapiadapter

import (
	"context"
	"log"
	"net/http"
	"time"
)

type StdAPIServer struct {
	httpServer *http.Server
}

// NewStdAPIServer creates and configures a new server and we configure timeouts here
func NewStdAPIServer(addr string, router http.Handler) *StdAPIServer {
	httpServer := &http.Server{
		Addr:    addr,
		Handler: router,
		// timeouts to prevent slow-client attacks
		ReadTimeout:  5 * time.Second,   // Max time entire request including body
		WriteTimeout: 10 * time.Second,  // Max time before writes time out
		IdleTimeout:  120 * time.Second, // Max time to wait for new request on a keep-alive connection
	}
	return &StdAPIServer{
		httpServer: httpServer,
	}
}

// ListenAndServe starts the HTTP server. It's a blocking call.
func (s *StdAPIServer) ListenAndServe() error {
	log.Printf("Starting net/http server on %s", s.httpServer.Addr)
	if err := s.httpServer.ListenAndServe(); err != nil {
		return err // return any error that is not graceful shutdown http.ErrServerClosed
	}

	return nil
}

// Shutdown gracefully shuts down the HTTP server
func (s *StdAPIServer) Shutdown(ctx context.Context) error {
	log.Printf("Shutting down net/http gracefully...")
	return s.httpServer.Shutdown(ctx)
}
