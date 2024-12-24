package am

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type App struct {
	Core Core
}

func NewApp(opts ...Option) *App {
	core := NewCore(opts...)
	return &App{
		Core: core,
	}
}

func (a *App) StartServer(server *http.Server, addr string) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	go func() {
		a.Log().Info("Starting server on ", addr)
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.Log().Errorf("Could not listen on %s: %v\n", addr, err)
		}
	}()

	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	a.Log().Info("Shutting down server on ", addr)
	err := server.Shutdown(ctx)
	if err != nil {
		a.Log().Errorf("Server forced to shutdown: %v", err)
	}

	a.Log().Info("Server stopped gracefully")
}

func (a *App) Log() Logger {
	return a.Core.Log()
}

func (a *App) Cfg() *Config {
	return a.Core.Cfg()
}
