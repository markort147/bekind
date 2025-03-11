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

func moviesIds(c echo.Context, criteria *movies.FindCriteria, sorting *movies.SortInfo) error {
	movieList := movies.Find(criteria, sorting)
	ids := make([]int, len(movieList))
	for i := range len(movieList) {
		ids[i] = movieList[i].Id
	}
	return c.Render(200, "movie-list", ids)
}

func sortMovies(c echo.Context) error {
	sortedBy := movies.StrToMF(c.QueryParam("by"))
	newSorting := movies.SortInfo{
		SortedBy: sortedBy,
		Desc:     movies.CurrentSorting.SortedBy == sortedBy && !movies.CurrentSorting.Desc,
	}

	return moviesIds(c, nil, &newSorting)
}

func filterMovies(c echo.Context) error {
	criteria := movies.FindCriteria{}

	// title
	//criteria.Title = strings.ReplaceAll(c.FormValue("title"), " ", "")
	criteria.Title = strings.TrimSpace(c.FormValue("title"))

	// rate range
	rateRange := strings.ReplaceAll(c.FormValue("rate"), " ", "")
	if rateRange != "" {
		minRate := uint8(0)
		maxRate := uint8(10)
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

	// year range
	yearRange := strings.ReplaceAll(c.FormValue("year"), " ", "")
	if yearRange != "" {
		minYear := uint16(0)
		maxYear := ^uint16(0)
		years := strings.Split(yearRange, "-")
		if len(years) > 0 {
			value, err := strconv.Atoi(years[0])
			if err == nil && value >= int(minYear) {
				minYear = uint16(value)
			}
		}
		if strings.ContainsRune(yearRange, '-') {
			if len(years) == 2 {
				value, err := strconv.Atoi(years[1])
				if err == nil && value <= int(maxYear) {
					maxYear = uint16(value)
				}
			}
		} else {
			maxYear = minYear
		}
		criteria.Year = []uint16{minYear, maxYear}
	}

	log.Logger.Debugf("filterMovies. criteria: %+v", criteria)
	return moviesIds(c, &criteria, nil)
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

	deleted := movies.DeleteById(id)
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
	return c.Render(200, "movies", nil)
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
	return c.Render(200, "movies", nil)
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
