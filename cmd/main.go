package main

import (
	"fmt"

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

	ms.Init()
	ms.FillForTests()

	e := initEcho()
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", cfg.GetConfig().Server.Port)))
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
