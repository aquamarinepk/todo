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
	*Router
	fs embed.FS
}

// NewFileServer creates a new FileServer.
func NewFileServer(fs embed.FS, opts ...Option) *FileServer {
	routerName := fmt.Sprintf("%s-router", fileServerName)

	r := NewRouter(routerName, opts...)
	return &FileServer{
		Router: r,
		fs:     fs,
	}
}

func (f *FileServer) SetupRoutes() error {
	cfg := f.Router.Cfg()
	if cfg.BoolVal(Key.ServerIndexEnabled, false) {
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

	f.Router.HandleFunc(staticPath+"/*", func(w http.ResponseWriter, r *http.Request) {
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

	f.Router.HandleFunc(staticPath+"/*", func(w http.ResponseWriter, r *http.Request) {
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
//func (f *FileServer) Router() *Router {
//	return f.router
//}

// Setup is the default implementation for the Setup method in FileServer.
func (f *FileServer) Setup(ctx context.Context) error {
	err := f.Router.Setup(ctx)
	if err != nil {
		return err
	}

	return f.SetupRoutes()
}
