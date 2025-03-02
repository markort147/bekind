package main

import (
	"github.com/markort147/bekind/cmd/bekindrewind/pkg/movies"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/markort147/bekind/internal/log"
)

/*
==================== HANDLER FUNCTIONS ====================
These functions are used to handle the requests from the client.
They are called by the router when a request is made to the server.
The functions are responsible for processing the request and returning a response.
The response consists of a status code and a body (HTML template).
The body is rendered using the Go template engine.
===========================================================
*/

// MovieList struct is used to pass data to the movie-list template.
// The struct contains the list of ms to display, the field to sort by, and the sorting order.
type MovieList struct {
	SortedBy string
	Desc     bool
	Body     []movies.Movie
}

// sortMovies is a handler function that returns the "movie-list" template with the list of ms sorted by the given field.
// The field to sort by is specified in the "by" query parameter.
func sortMovies(c echo.Context) error {

	sortedBy := movies.StrToMF(c.QueryParam("by"))
	newSorting := movies.SortInfo{
		SortedBy: sortedBy,
		Desc:     movies.CurrentSorting.SortedBy == sortedBy && !movies.CurrentSorting.Desc,
	}

	movieFieldLabel, err := movies.MFToStr(newSorting.SortedBy)
	if err != nil {
		panic(err) // it should be impossible to get an error here
	}

	return c.Render(200, "movie-list", MovieList{
		SortedBy: movieFieldLabel,
		Desc:     newSorting.Desc,
		Body:     movies.Find(nil, &newSorting),
	})
}

// staticView is a handler function that returns a simple view with the given template name.
// The function is used to render static HTML pages that do not require any data from the server.
func staticView(templateName string) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.Render(200, templateName, nil)
	}
}

// deleteMovie is a handler function that deletes a movie from the database.
// The id of the movie to delete is specified in the URL path.
// The function returns a 200 status code if the movie was deleted successfully, or a 404 status code if the movie was not found.
func deleteMovie(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.NoContent(echo.ErrBadRequest.Code)
	}

	deleted := movies.DeleteById(id)
	if !deleted {
		return c.NoContent(echo.ErrNotFound.Code)
	}
	return c.NoContent(200)
}

// findMovie is a handler function that returns the "movie-list" template with the list of ms that match the given ids.
// The ids are specified by the "id" form value, which is a comma-separated list of movie ids.
func findMovie(c echo.Context) error {
	criteria := movies.FindCriteria{}

	// title
	criteria.Title = strings.ReplaceAll(c.FormValue("title"), " ", "")

	// rate range
	rateRange := strings.ReplaceAll(c.FormValue("rate"), " ", "")
	if rateRange != "" {
		var minRate uint8 = 0
		var maxRate uint8 = 10
		rates := strings.Split(rateRange, "-")
		if len(rates) > 0 {
			value, err := strconv.Atoi(rates[0])
			if err == nil && value >= 0 {
				minRate = uint8(value)
			}
		}
		if strings.ContainsRune(rateRange, '-') {
			if len(rates) == 2 {
				value, err := strconv.Atoi(rates[1])
				if err == nil && value <= 10 {
					maxRate = uint8(value)
				}
			}
		} else {
			maxRate = minRate
		}
		criteria.Rate = []uint8{minRate, maxRate}
	}

	return c.Render(200, "movie-list", MovieList{
		SortedBy: "id",
		Desc:     false,
		Body:     movies.Find(&criteria, nil),
	})
}

// postMovie is a handler function that creates a new movie and adds it to the database.
// The movie data is specified in the form data of the request.
// The function returns the "movie-list" template with the updated list of ms.
func postMovie(c echo.Context) error {
	rate, _ := strconv.Atoi(c.FormValue("rate"))
	movies.Save(movies.Movie{
		Title: c.FormValue("title"),
		Year:  c.FormValue("year"),
		Rate:  uint8(rate),
	})
	/*return c.Render(200, "movie-list", MovieList{
		SortedBy: "id",
		Desc:     false,
		Body:     ms.FindAll(nil),
	})*/
	return c.Render(200, "search_movie", nil)
}

// editMovieView is a handler function that returns the "edit_movie" template with the movie data to edit.
// The id of the movie to edit is specified in the URL path.
func editMovieView(c echo.Context) error {
	strId := c.Param("id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		log.Logger.Fatal(err)
	}

	movie, err := movies.FindById(id)
	if err != nil {
		log.Logger.Fatal(err)
	}

	return c.Render(200, "edit_movie", movie)
}

// putMovie is a handler function that updates the movie data in the database.
// The id of the movie to update is specified in the URL path.
// The new movie data is specified in the form data of the request.
// The function returns the "movie-list" template with the updated list of ms.
func putMovie(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	rate, _ := strconv.Atoi(c.FormValue("rate"))
	movies.Update(id, movies.Movie{
		Title: c.FormValue("title"),
		Year:  c.FormValue("year"),
		Rate:  uint8(rate),
	})
	/*return c.Render(200, "movie-list", MovieList{
		SortedBy: "id",
		Desc:     false,
		Body:     ms.FindAll(nil),
	})*/
	return c.Render(200, "search_movie", nil)
}

// validateYear is a helper function that validates the year format.
func validateYear(c echo.Context) error {
	var message string
	valid := true
	value := c.FormValue("year")

	if value == "" {
		valid = false
		message = "It cannot be empty."
	} else {
		year, err := strconv.Atoi(value)
		if err != nil || year < 0 {
			valid = false
			message = "Invalid year format."
		}
	}

	return c.Render(200, "year_input", struct {
		Value   string
		Valid   bool
		Message string
	}{
		Value:   value,
		Valid:   valid,
		Message: message,
	})
}

// validateRate is a helper function that validates the rate format.
func validateRate(c echo.Context) error {
	var message string
	valid := true
	value := c.FormValue("rate")

	if value == "" {
		valid = false
		message = "It cannot be empty."
	} else {
		rate, err := strconv.Atoi(value)
		if err != nil || rate < 0 || rate > 10 {
			valid = false
			message = "Invalid rate format."
		}
	}

	return c.Render(200, "rate_input", struct {
		Value   string
		Valid   bool
		Message string
	}{
		Value:   value,
		Valid:   valid,
		Message: message,
	})
}

// validateTitle is a helper function that validates the title format.
func validateTitle(c echo.Context) error {
	var message string
	title := c.FormValue("title")
	valid := true

	if strings.ReplaceAll(title, " ", "") == "" {
		valid = false
		message = "It cannot be empty."
	}

	return c.Render(200, "title_input", struct {
		Value   string
		Valid   bool
		Message string
	}{
		Value:   title,
		Valid:   valid,
		Message: message,
	})
}

func getMovie(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.NoContent(echo.ErrBadRequest.Code)
	}

	movie, err := movies.FindById(id)
	if err != nil {
		return c.NoContent(echo.ErrBadRequest.Code)
	}

	return c.Render(200, "movie", movie)
}
