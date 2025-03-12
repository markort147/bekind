package main

import (
	"github.com/markort147/bekind/cmd/bekindrewind/pkg/movies"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/markort147/gopkg/log"
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

func getMoviesView(c echo.Context) error {
	return c.Render(200, "movie-list", movies.CurrView)
}

func updateMoviesViewSorting(c echo.Context) error {
	movies.CurrView.Sort(c.QueryParam("by"))
	return c.Render(200, "movie-list", movies.CurrView)
}

func updateMoviesViewFilter(c echo.Context) error {
	criteria := movies.FilterCriteria{}
	criteria.Title = strings.TrimSpace(c.FormValue("title"))
	criteria.Rate = strings.ReplaceAll(c.FormValue("rate"), " ", "")
	criteria.Year = strings.ReplaceAll(c.FormValue("year"), " ", "")
	movies.CurrView.Filter(criteria)
	return c.Render(200, "movie-list", movies.CurrView)
}

func staticView(templateName string) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.Render(200, templateName, nil)
	}
}

func deleteMovie(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.NoContent(echo.ErrBadRequest.Code)
	}

	deleted := movies.Delete(id)
	if !deleted {
		return c.NoContent(echo.ErrNotFound.Code)
	}
	return c.NoContent(200)
}

func postMovie(c echo.Context) error {
	rate, _ := strconv.Atoi(c.FormValue("rate"))
	year, _ := strconv.Atoi(c.FormValue("year"))
	movies.Save(movies.Movie{
		Title: c.FormValue("title"),
		Year:  uint16(year),
		Rate:  uint8(rate),
	})
	return c.Render(200, "movie-list", movies.CurrView)
}

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

func updateMovie(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	rate, _ := strconv.Atoi(c.FormValue("rate"))
	year, _ := strconv.Atoi(c.FormValue("year"))
	movies.Update(id, movies.Movie{
		Title: c.FormValue("title"),
		Year:  uint16(year),
		Rate:  uint8(rate),
	})
	return c.Render(200, "movie-list", movies.CurrView)
}

type fieldValidation struct {
	Name    string
	Label   string
	Value   string
	Valid   bool
	Message string
}

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
		} else if year > int(^uint16(0)) {
			valid = false
			message = "Year exceeds the maximum value of " + strconv.Itoa(int(^uint16(0)))
		}
	}

	return c.Render(200, "movie_input", fieldValidation{
		Name:    "year",
		Label:   "Year:",
		Value:   value,
		Valid:   valid,
		Message: message,
	})
}

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

	return c.Render(200, "movie_input", fieldValidation{
		Name:    "rate",
		Label:   "Rate:",
		Value:   value,
		Valid:   valid,
		Message: message,
	})
}

func validateTitle(c echo.Context) error {
	var message string
	title := c.FormValue("title")
	valid := true

	if strings.ReplaceAll(title, " ", "") == "" {
		valid = false
		message = "It cannot be empty."
	}

	return c.Render(200, "movie_input", fieldValidation{
		Name:    "title",
		Label:   "Title:",
		Value:   title,
		Valid:   valid,
		Message: message,
	})
}

func getMovieDetails(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.NoContent(echo.ErrBadRequest.Code)
	}

	movie, err := movies.FindById(id)
	if err != nil {
		return c.NoContent(echo.ErrBadRequest.Code)
	}

	return c.Render(200, "movie_details", movie)
}

func getMovieRow(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.NoContent(echo.ErrBadRequest.Code)
	}

	movie, err := movies.FindById(id)
	if err != nil {
		return c.NoContent(echo.ErrBadRequest.Code)
	}

	return c.Render(200, "movie_row", movie)
}
