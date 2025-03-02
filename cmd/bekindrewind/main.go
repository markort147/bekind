package main

import (
	"embed"
	"flag"
	"fmt"
	ms "github.com/markort147/bekind/cmd/bekindrewind/pkg/movies"
	"os"

	"github.com/markort147/bekind/internal/echotmpl"
	"github.com/markort147/bekind/internal/log"
	"github.com/markort147/bekind/internal/ymlcfg"
)

//go:embed assets/*
var assetsFS embed.FS

func main() {

	// parse config file path
	var configPath string
	flag.StringVar(&configPath, "config", "", "Path to the configuration file")
	flag.Parse()
	if configPath == "" {
		flag.Usage()
		os.Exit(1)
	}

	// load configuration
	cfg, err := ymlcfg.ParseFile[Config](configPath)
	if err != nil {
		_, err1 := fmt.Fprintf(os.Stderr, "Error loading config: %v", err)
		if err1 != nil {
			panic(err1)
		}
		os.Exit(1)
	}
	logLevel := parseLogLevel(cfg.Log.Level)
	logOutput, closeFunc := parseLogOutput(cfg.Log.Output)
	if closeFunc != nil {
		defer closeFunc()
	}

	// log configuration
	if err = log.Init(&log.Config{
		Output: logOutput,
		Level:  logLevel,
	}); err != nil {
		_, err1 := fmt.Fprintf(os.Stderr, "Error init logger: %v", err)
		if err1 != nil {
			panic(err1)
		}
		os.Exit(1)
	}

	// movies service initialization
	ms.Init()
	ms.FillForTests()

	// echo server initialization
	wgServer, err := echotmpl.StartServer(&echotmpl.Config{
		Port:          cfg.Server.Port,
		FileSystem:    assetsFS,
		LogOutputPath: logOutput,
		LogLevel:      logLevel,
		RoutesRegister: func(e *echotmpl.Echo) {
			// views
			e.GET("/views/movies", staticView("movies"))
			e.PUT("/views/movies/filter", findMovie)
			e.PUT("/views/movies/sort", sortMovies)
			e.GET("/views/add-movie", staticView("add_movie"))
			e.GET("/views/edit-movie/:id", editMovieView)
			// movie
			e.GET("/movie/:id/details", getMovieDetails)
			e.GET("/movie/:id/row", getMovieRow)
			e.POST("/movie", postMovie)
			e.PUT("/movie/:id", putMovie)
			e.DELETE("/movie/:id", deleteMovie)
			// validate
			e.POST("/validate/title", validateTitle)
			e.POST("/validate/year", validateYear)
			e.POST("/validate/rate", validateRate)
		},
		CustomFuncs: echotmpl.FuncMap{
			"WrapStringValidation": WrapStringValidation,
			"WrapUint8Validation":  WrapUint8Validation,
		},
	})
	if err != nil {
		log.Logger.Fatalf("Error starting server: %v", err)
	}
	defer log.Logger.Info("Server exited")

	// wait for the server to exit
	wgServer.Wait()
}
