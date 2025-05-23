package service

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"time"
)

const DefaultServiceTimeout = 1 * time.Minute

// FIXME: Simple CORS middleware. Improve this
func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*") // Allow all origins
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-custom_type, Authorization")

		// Handle preflight (OPTIONS) request
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Proceed with the next handler if not preflight
		next.ServeHTTP(w, r)
	})
}

// Used for testing
type HttpService struct {
	Name    string
	Address string
	Handler http.Handler
}

func (s *HttpService) String() string {
	return s.Name
}

func (s *HttpService) Start(ctx context.Context, ready chan<- struct{}, logger *slog.Logger) error {
	server := http.Server{
		Addr:     s.Address,
		Handler:  CorsMiddleware(s.Handler),
		ErrorLog: slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	listener, err := net.Listen("tcp", s.Address)
	if err != nil {
		return err
	}

	logger.Info("HTTP server started listening", "service", s, "port", listener.Addr())
	ready <- struct{}{}

	done := make(chan error, 1)
	go func() {
		err := server.Serve(listener)
		if !errors.Is(err, http.ErrServerClosed) {
			logger.Warn("Service exited with error", "service", s, "error", err)
		}
		done <- err
	}()

	select {
	case err = <-done:
		return err
	case <-ctx.Done():
		ctx, cancel := context.WithTimeout(context.Background(), DefaultServiceTimeout)
		defer cancel()
		return server.Shutdown(ctx)
	}
}
