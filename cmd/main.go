package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	hdls "github.com/bekind/bekindfrontend/handlers"
	cfg "github.com/bekind/bekindfrontend/internal/config"
	"github.com/bekind/bekindfrontend/internal/log"
	ms "github.com/bekind/bekindfrontend/movies"

	"github.com/labstack/echo/v4"
	mdw "github.com/labstack/echo/v4/middleware"
)

func main() {
	cfg.FromFile("config.yml")

	log.Init()
	defer log.Close()

	ms.Init()
	ms.FillForTests()

	e := initEcho()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
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
	}()

	// Start the server in a separate goroutine
	go func() {
		if err := e.Start(fmt.Sprintf(":%d", cfg.GetConfig().Server.Port)); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("Error starting the server: ", err)
		}
	}()

	wg.Wait()
	log.Logger.Info("Server exited")
}

// initEcho initializes the echo server
// and registers all the endpoints
func initEcho() *echo.Echo {
	e := echo.New()
	e.Use(mdw.Logger())
	e.Debug = true
	e.Renderer = newTemplate("static/templates/*")
	registerEndpoints(e)
	return e
}

func registerEndpoints(e *echo.Echo) {
	// home page
	e.File("/", "static/index.html")

	// side panels views
	e.GET("/views/movies", hdls.GetAllMovie)
	e.GET("/views/search_movie", hdls.SimpleView("search_movie"))
	e.GET("/views/add_movie", hdls.SimpleView("add_movie"))
	e.GET("/views/edit-movie/:id", hdls.EditMovie)
	e.GET("/views/movies/sort", hdls.SortMovie)

	// movie handlers
	e.DELETE("/movie/:id", hdls.DeleteMovie)
	e.GET("/movie", hdls.GetMovie)
	e.POST("/movie", hdls.PostMovie)
	e.PUT("/movie/:id", hdls.PutMovie)
}
