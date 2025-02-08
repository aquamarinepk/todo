/*
Package am provides functionality for serving static files from an embedded filesystem.

NOTE: The path to the static files is currently fixed.
In the future, we can consider delegating the path (mounting point) for the FileServer to the app setup in main, as is done for other handlers.
*/
package am

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"strings"
)

const (
	fileServerName = "file-server"
	staticPath     = "/static"
	assetsFilePath = "assets/static"
)

// FileServer serves static files from an embedded filesystem.
type FileServer struct {
	Core   Core
	router *Router
	fs     embed.FS
}

// NewFileServer creates a new FileServer.
func NewFileServer(fs embed.FS, opts ...Option) *FileServer {
	core := NewCore(fileServerName, opts...)

	routerName := fmt.Sprintf("%s-router", fileServerName)

	r := NewRouter(routerName, opts...)
	return &FileServer{
		Core:   core,
		router: r,
		fs:     fs,
	}
}

func (f *FileServer) SetupRoutes() error {
	if f.Cfg().BoolVal(Key.ServerIndexEnabled, false) {
		return f.SetupRoutesIndex()
	}

	return f.SetupRoutesNoIndex()
}

// SetupRoutesIndex sets up the routes to serve static files index listing.
func (f *FileServer) SetupRoutesIndex() error {
	staticFS, err := fs.Sub(f.fs, assetsFilePath)
	if err != nil {
		return fmt.Errorf("failed to create sub filesystem: %w", err)
	}

	server := http.FileServer(http.FS(staticFS))

	f.router.HandleFunc(staticPath+"/*", func(w http.ResponseWriter, r *http.Request) {
		http.StripPrefix(staticPath, server).ServeHTTP(w, r)
	})

	return nil
}

// SetupRoutesNoIndex sets up the routes to serve static files without index listing.
func (f *FileServer) SetupRoutesNoIndex() error {
	staticFS, err := fs.Sub(f.fs, assetsFilePath)
	if err != nil {
		return fmt.Errorf("failed to create sub filesystem: %w", err)
	}

	fileServer := http.FileServer(http.FS(staticFS))

	f.router.HandleFunc(staticPath+"/*", func(w http.ResponseWriter, r *http.Request) {
		requestedFile := strings.TrimPrefix(r.URL.Path, staticPath+"/")

		f, err := staticFS.Open(requestedFile)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		defer f.Close()

		stat, err := f.Stat()
		if err != nil || stat.IsDir() {
			http.NotFound(w, r)
			return
		}

		http.StripPrefix(staticPath, fileServer).ServeHTTP(w, r)
	})

	return nil
}

// Router returns the underlying chi.Router.
func (f *FileServer) Router() *Router {
	return f.router
}

// Name returns the name in FileServer.
func (f *FileServer) Name() string {
	return f.Core.Name()
}

// SetName sets the name in FileServer.
func (f *FileServer) SetName(name string) {
	f.Core.SetName(name)
}

// Log returns the Logger in FileServer.
func (f *FileServer) Log() Logger {
	return f.Core.Log()
}

// SetLog sets the Logger in FileServer.
func (f *FileServer) SetLog(log Logger) {
	f.Core.SetLog(log)
}

// Cfg returns the Config in FileServer.
func (f *FileServer) Cfg() *Config {
	return f.Core.Cfg()
}

// SetCfg sets the Config in FileServer.
func (f *FileServer) SetCfg(cfg *Config) {
	f.Core.SetCfg(cfg)
}

// Setup is the default implementation for the Setup method in FileServer.
func (f *FileServer) Setup(ctx context.Context) error {
	err := f.Core.Setup(ctx)
	if err != nil {
		return err
	}

	return f.SetupRoutes()
}

// Start is the default implementation for the Start method in FileServer.
func (f *FileServer) Start(ctx context.Context) error {
	return f.Core.Start(ctx)
}

// Stop is the default implementation for the Stop method in FileServer.
func (f *FileServer) Stop(ctx context.Context) error {
	return f.Core.Stop(ctx)
}
