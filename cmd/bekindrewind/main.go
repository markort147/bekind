package main

import (
	"embed"
	"fmt"
	"os"
	"strconv"
)

//go:embed assets/*
var assetsFS embed.FS

// data global data storage
var data Data

// moviesView global movies view
var moviesView *MoviesView

func main() {
	// get environment variables
	port, _ := strconv.Atoi(os.Getenv("BEKIND_PORT"))
	logLevel := os.Getenv("BEKIND_LOG_LEVEL")
	logOutput := os.Getenv("BEKIND_LOG_OUTPUT")
	mode := os.Getenv("BEKIND_MODE")

	// log configuration
	parsedLogLevel := parseLogLevel(logLevel)
	parsedLogOutput, closeFunc := parseLogOutput(logOutput)
	if closeFunc != nil {
		defer closeFunc()
	}
	if err := InitLog(&LogConfig{
		Output: parsedLogOutput,
		Level:  parsedLogLevel,
	}); err != nil {
		_, err1 := fmt.Fprintf(os.Stderr, "Error init logger: %v", err)
		if err1 != nil {
			panic(err1)
		}
	}
	Logger.SetHeader("${time_rfc3339} ${short_file}:${line} ${level} ${message}")

	// init data storage
	switch mode {
	//case "psql":
	//	psqlData, err := newPsqlData(os.Getenv("BEKIND_DB_CONN_STRING"))
	//	if err != nil {
	//		Logger.Fatalf("Error init psql: %v", err)
	//	}
	//	data = psqlData
	default:
		data = newMemoryData()
	}

	// init moviesView
	moviesView = newMoviesView()

	// echo server initialization
	wgServer, err := StartServer(&SrvConfig{
		Port:       port,
		Logger:     Logger,
		FileSystem: assetsFS,
		RoutesRegister: func(e *Echo) {
			e.GET("/views/movies", getViewsMovies)
			e.PUT("/views/movies/filter", putViewsMoviesFilter)
			e.PUT("/views/movies/sort", putViewsMoviesSort)
			e.GET("/views/add-movie", getViewsAddMovie)
			e.GET("/views/edit-movie/:id", getViewsEditMovie)
			e.GET("/views/data", getViewsData)
			e.GET("/views/people", getViewsPeople)

			e.GET("/movie/:id/details", getMovieDetails)
			e.GET("/movie/:id/row", getMovieRow)
			e.POST("/movie", postMovie)
			e.PUT("/movie/:id", putMovie)
			e.DELETE("/movie/:id", deleteMovie)

			e.GET("people/table-body", getPeopleTableBody)
			e.GET("/person/:id/avgrate", getPersonAvgRate)
			e.GET("/person/:id/name", getPersonName)
			e.GET("/person/:id/table-row", getPersonTableRow)

			e.POST("/validate/title", postValidateTitle)
			e.POST("/validate/year", postValidateYear)
			e.POST("/validate/rate", postValidateRate)

			e.POST("/upload", postUpload)
			e.GET("/download", getDownload)
		},
		CustomFuncs: map[string]any{
			"IndexPair": func(i, j int) map[string]int {
				return map[string]int{
					"Curr": i,
					"Next": j,
				}
			},
		},
	})
	if err != nil {
		Logger.Fatalf("Error starting server: %v", err)
	}
	defer Logger.Info("Server exited")

	// wait for the server to exit
	wgServer.Wait()
}
