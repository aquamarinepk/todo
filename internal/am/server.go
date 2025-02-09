package am

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type Server struct {
	Core
	server *http.Server
}

// NewServer creates a new Server instance with the provided options.
//
// Parameters:
// - hostKey: The configuration key for the server host (e.g., "server.web.host").
// - portKey: The configuration key for the server port (e.g., "server.web.port").
// - handler: The HTTP handler to use for the server.
// - opts: Variadic options to configure the server.
//
// Example usage:
//
//	webServer := am.NewServer("server.web.host", "server.web.port", webRouter, opts...)
//	apiServer := am.NewServer("server.api.host", "server.api.port", apiRouter, opts...)
//
// Note: This is a WIP and an improved way of configuring the server will be provided in future updates.
func NewServer(hostKey, portKey string, handler http.Handler, opts ...Option) *Server {
	core := NewCore("", opts...)
	cfg := core.Cfg()

	host := cfg.StrValOrDef(hostKey, "localhost")
	port := cfg.StrValOrDef(portKey, "8080")

	server := &http.Server{
		Addr:    host + ":" + port,
		Handler: handler,
	}
	return &Server{
		Core:   core,
		server: server,
	}
}

func (s *Server) Start() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	go func() {
		s.Log().Info("Starting server on ", s.server.Addr)
		err := s.server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.Log().Errorf("Could not listen on %s: %v\n", s.server.Addr, err)
		}
	}()

	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	s.Log().Info("Shutting down server on ", s.server.Addr)
	err := s.server.Shutdown(ctx)
	if err != nil {
		s.Log().Errorf("Server forced to shutdown: %v", err)
	}

	s.Log().Info("Server stopped gracefully")
}
