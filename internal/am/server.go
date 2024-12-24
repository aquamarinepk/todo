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
	Core   Core
	server *http.Server
}

// NewServer creates a new Server instance with the provided opts.
// TODO: port should be removed from the arguments and be set in the config.
func NewServer(port string, handler http.Handler, opts ...Option) *Server {
	core := NewCore(opts...)

	server := &http.Server{
		Addr:    ":" + port,
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

func (s *Server) Log() Logger {
	return s.Core.Log()
}

func (s *Server) Cfg() *Config {
	return s.Core.Cfg()
}
