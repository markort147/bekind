package main

import (
	"github.com/bekind/bekindfrontend/log"
	"github.com/bekind/bekindfrontend/movies"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	log.Init()	
	movies.InitWithRandoms(500)
	e := initEcho()
	e.Logger.Fatal(e.Start(":8080"))
}

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
	e.GET("/movies", GetAllMovieHandler)
	e.GET("/search_movie", SimpleViewHandler("search_movie"))
	e.GET("/add_movie", SimpleViewHandler("add_movie"))

	// movie handlers
	e.DELETE("/movie/:id", DeleteMovieHandler)
	e.GET("/movie", GetMovieHandler)
	e.POST("/movie", PostMovieHandler)
}
