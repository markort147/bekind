package main

import (
	"context"
	"embed"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/markort147/bekindfrontend/internal/api"
	"github.com/markort147/bekindfrontend/internal/config"
	"github.com/markort147/bekindfrontend/internal/log"
	ms "github.com/markort147/bekindfrontend/internal/movies"
	"github.com/markort147/bekindfrontend/internal/utils"

	"github.com/labstack/echo/v4"
	mdw "github.com/labstack/echo/v4/middleware"
)

//go:embed assets/*
var assets embed.FS

func main() {

	var configPath string
	flag.StringVar(&configPath, "config", "", "Path to the configuration file")
	flag.Parse()
	if configPath == "" {
		flag.Usage()
		os.Exit(1)
	}

	cfg, err := config.FromFile("config.yml")
	if err != nil {
		log.Logger.Fatalf("Error loading config: %v", err)
	}

	// log configuration
	log.Init(cfg)
	defer log.Close()

	// movies service initialization
	ms.Init()
	ms.FillForTests()

	// echo server initialization
	e, wg := initEcho(cfg)
	go startServer(e, cfg)
	defer log.Logger.Info("Server exited")

	// wait for the server to exit
	wg.Wait()
}

// initEcho initializes the echo server
// and registers all the endpoints
func initEcho(cfg *config.Config) (*echo.Echo, *sync.WaitGroup) {
	// create and configure the echo server
	e := echo.New()
	e.Use(mdw.LoggerWithConfig(mdw.LoggerConfig{
		Output: log.ParseOutput(cfg.Log.Output),
	}))
	e.Logger.SetLevel(log.ParseLevel(cfg.Log.Level))
	e.Use(mdw.Recover())
	e.Renderer = utils.NewTemplateRendererFromFS(assets, "assets/templates/*")
	registerEndpoints(e)

	// intercept shutdown signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	var wg sync.WaitGroup
	wg.Add(1)
	go interceptShutdown(e, quit, &wg)

	return e, &wg
}

func startServer(e *echo.Echo, cfg *config.Config) {
	if err := e.Start(fmt.Sprintf(":%d", cfg.Server.Port)); err != nil && err != http.ErrServerClosed {
		e.Logger.Fatal("Error starting the server: ", err)
	}
}

func interceptShutdown(e *echo.Echo, quit chan os.Signal, wg *sync.WaitGroup) {
	defer wg.Done()
	<-quit
	log.Logger.Info("Shutting down the server")

	// Create a context with a timeout to allow for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Attempt to gracefully shut down the server
	if err := e.Shutdown(ctx); err != nil {
		log.Logger.Error("Server forced to shutdown: ", err)
	}

	log.Logger.Info("Server exiting")
}

func registerEndpoints(e *echo.Echo) {
	// home page
	e.FileFS("/", "index.html", echo.MustSubFS(assets, "assets"))

	// side panels views
	e.GET("/views/movies", api.GetAllMovie)
	e.GET("/views/search_movie", api.SimpleView("search_movie"))
	e.GET("/views/add_movie", api.SimpleView("add_movie"))
	e.GET("/views/edit-movie/:id", api.EditMovie)
	e.GET("/views/movies/sort", api.SortMovie)

	// movie handlers
	e.DELETE("/movie/:id", api.DeleteMovie)
	e.GET("/movie", api.GetMovie)
	e.POST("/movie", api.PostMovie)
	e.PUT("/movie/:id", api.PutMovie)
}
