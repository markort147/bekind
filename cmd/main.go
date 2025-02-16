package main

import (
	"fmt"

	"github.com/bekind/bekindfrontend/config"
	"github.com/bekind/bekindfrontend/log"
	"github.com/bekind/bekindfrontend/movies"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	config.FromFile("config.yml")

	log.Init()

	movies.Init()
	movies.FillForTests()

	e := initEcho()
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", config.GetConfig().Server.Port)))
}

// initEcho initializes the echo server
// and registers all the endpoints
func initEcho() *echo.Echo {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Debug = true
	e.Renderer = newTemplate("static/templates/*")
	registerEndpoints(e)
	return e
}

func registerEndpoints(e *echo.Echo) {
	// home page
	e.File("/", "static/index.html")

	// side panels views
	e.GET("/views/movies", GetAllMovieHandler)
	e.GET("/views/search_movie", SimpleViewHandler("search_movie"))
	e.GET("/views/add_movie", SimpleViewHandler("add_movie"))
	e.GET("/views/edit-movie/:id", EditMovieHandler)
	e.GET("/views/movies/sort", SortMovieHandler)

	// movie handlers
	e.DELETE("/movie/:id", DeleteMovieHandler)
	e.GET("/movie", GetMovieHandler)
	e.POST("/movie", PostMovieHandler)
	e.PUT("/movie/:id", PutMovieHandler)
}
