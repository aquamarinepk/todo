package am

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

const (
	resPath = "/res"
)

type App struct {
	Core
	opts          []Option
	Version       string
	Router        *Router
	APIRouter     *Router
	APIRouters    map[string]*Router
	ResRouter     *Router
	ResAPIRouter  *Router
	ResAPIRouters map[string]*Router
	deps          map[string]*Dep
	depOrder      []string
	depsMutex     sync.Mutex
	fs            embed.FS
}

func NewApp(name, version string, fs embed.FS, opts ...Option) *App {
	core := NewCore(name, opts...)
	core.SetName(name)
	for _, opt := range opts {
		opt(core)
	}

	if core.Log() == nil {
		core.SetLog(NewLogger("info"))
		opts = append(opts, WithLog(core.Log()))
	}

	if core.Cfg() == nil {
		core.SetCfg(NewConfig())
		opts = append(opts, WithCfg(core.Cfg()))
	}

	app := &App{
		opts:          opts,
		Core:          core,
		Router:        NewWebRouter("web-router", opts...),
		APIRouter:     NewWebRouter("api-router", opts...),
		APIRouters:    make(map[string]*Router),
		ResRouter:     NewWebRouter("res-router", opts...),
		ResAPIRouter:  NewWebRouter("res-api-router", opts...),
		ResAPIRouters: make(map[string]*Router),
		fs:            fs,
		deps:          make(map[string]*Dep),
	}

	resPath := app.Cfg().StrValOrDef(Key.ServerResPath, resPath)

	app.Router.Mount("/api", app.APIRouter)
	app.Router.Mount(resPath, app.ResRouter)
	app.ResRouter.Mount(resPath, app.ResAPIRouter)

	return app
}

func (a *App) Add(dep Core) {
	err := a.checkSetup()
	if err != nil {
		a.Log().Errorf("cannot add dependency: %v", err)
		return
	}

	if dep.Name() == "" {
		dep.SetName(genName())
	}

	dep.SetOpts(a.opts...)

	dep.SetLog(a.Log())
	dep.SetCfg(a.Cfg())

	a.Log().Infof("Adding dependency: %s", dep.Name())

	a.depsMutex.Lock()
	defer a.depsMutex.Unlock()

	a.deps[dep.Name()] = &Dep{
		Core:   dep,
		Status: Stopped,
	}
	a.depOrder = append(a.depOrder, dep.Name())
}

func (a *App) Dep(name string) (*Dep, bool) {
	a.depsMutex.Lock()
	defer a.depsMutex.Unlock()
	dep, ok := a.deps[name]
	return dep, ok
}

func (a *App) Setup(ctx context.Context) error {
	var errs []string

	// Debug the content of deps
	a.depsMutex.Lock()
	order := make([]string, len(a.depOrder))
	copy(order, a.depOrder)
	depsCopy := make(map[string]*Dep, len(a.deps))
	for k, v := range a.deps {
		depsCopy[k] = v
	}
	a.depsMutex.Unlock()

	for _, name := range order {
		dep, ok := depsCopy[name]
		if !ok {
			continue
		}
		if _, ok := dep.Core.(Core); !ok {
			msg := fmt.Sprintf("cannot setup %s: not a core dep", dep.Core.Name())
			a.Log().Info(msg)
			continue
		}
		a.Log().Info("setting up ", dep.Name())
		err := dep.Setup(ctx)
		if err != nil {
			msg := fmt.Sprintf("canot setup %s: %v", dep.Name(), err)
			errs = append(errs, msg)
		}
	}

	if a.Log() == nil {
		errs = append(errs, "logging services not available")
	}

	if a.Cfg() == nil {
		errs = append(errs, "config not available")
	}

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "; "))
	}

	return nil
}

func (a *App) Start(ctx context.Context) error {
	// Start all dependencies
	var errs []string

	a.depsMutex.Lock()
	order := make([]string, len(a.depOrder))
	copy(order, a.depOrder)
	depsCopy := make(map[string]*Dep, len(a.deps))
	for k, v := range a.deps {
		depsCopy[k] = v
	}
	a.depsMutex.Unlock()

	for _, name := range order {
		dep, ok := depsCopy[name]
		if !ok {
			continue
		}
		if coreDep, ok := dep.Core.(Core); ok {
			err := coreDep.Start(ctx)
			if err != nil {
				msg := fmt.Sprintf("failed to start %s: %v", coreDep.Name(), err)
				errs = append(errs, msg)
			}
		}
	}

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "; "))
	}

	// Start the servers
	webAddr := a.Cfg().WebAddr()
	apiAddr := a.Cfg().APIAddr()

	if a.Cfg().BoolVal(Key.ServerWebEnabled, true) {
		webServer := &http.Server{
			Addr:    webAddr,
			Handler: a.Router,
		}
		go a.StartServer(webServer, webServer.Addr)
	}

	if a.Cfg().BoolVal(Key.ServerAPIEnabled, true) {
		apiServer := &http.Server{
			Addr:    apiAddr,
			Handler: a.APIRouter,
		}
		go a.StartServer(apiServer, apiServer.Addr)
	}

	return nil
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

func (a *App) Mount(path string, handler http.Handler) {
	a.Router.Mount(path, handler)
}

func (a *App) MountAPI(version, path string, handler http.Handler) {
	version = fmt.Sprintf("/%s", version)
	versionPath := fmt.Sprintf("%s%s", path, version)
	router, exists := a.APIRouters[version]
	if !exists {
		name := fmt.Sprintf("api-router-%s", versionPath)
		router = NewWebRouter(name, a.opts...)
		router.Mount(path, handler)
		a.APIRouters[versionPath] = router
	}
	a.APIRouter.Mount(version, router)
}

func (a *App) MountRes(path string, handler http.Handler) {
	a.ResRouter.Mount(path, handler)
}

func (a *App) MountResAPI(version, path string, handler http.Handler) {
	version = fmt.Sprintf("/%s", version)
	versionPath := fmt.Sprintf("%s%s", path, version)
	router, exists := a.ResAPIRouters[version]
	if !exists {
		name := fmt.Sprintf("res-api-router-%s", versionPath)
		router = NewWebRouter(name, a.opts...)
		router.Mount(path, handler)
		a.ResAPIRouters[versionPath] = router
	}
	a.ResAPIRouter.Mount(version, router)
}

func (a *App) MountFileServer(path string, fs *FileServer) {
	a.Mount(path, fs.Router)
}

func (app *App) SetWebRouter(router *Router) {
	app.Router = router
}

func (app *App) SetAPIRouter(router *Router) {
	app.APIRouter = router
}

func (app *App) MountWeb(path string, router *Router) {
	app.Mount(path, router)
}

func (app *App) MountResWeb(path string, router *Router) {
	app.MountRes(path, router)
}

func (a *App) checkSetup() error {
	if a.Log() == nil {
		return errors.New("logging services not available")
	}

	if a.Cfg() == nil {
		return errors.New("config not available")
	}

	return nil
}

func genName() string {
	u := uuid.New()
	segments := strings.Split(u.String(), "-")
	rand.Seed(time.Now().UnixNano())
	firstPart := make([]rune, 8)
	for i := range firstPart {
		firstPart[i] = 'a' + rune(rand.Intn(26))
	}
	return fmt.Sprintf("%s-%s", string(firstPart), segments[len(segments)-1])
}
