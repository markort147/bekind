package main

import (
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

func getViewsMovies(c echo.Context) error {
	return c.Render(200, "movie-list", moviesView)
}

func putViewsMoviesSort(c echo.Context) error {
	moviesView.setSortingBy(c.QueryParam("by"))
	moviesView.refresh()
	return getViewsMovies(c)
}

func putViewsMoviesFilter(c echo.Context) error {
	criteria := MoviesViewFilter{}
	criteria.Title = strings.TrimSpace(c.FormValue("title"))
	criteria.Rate = strings.ReplaceAll(c.FormValue("rate"), " ", "")
	criteria.Year = strings.ReplaceAll(c.FormValue("year"), " ", "")
	moviesView.setFilterCriteria(criteria)
	moviesView.refresh()
	return getViewsMovies(c)
}

func getViewsAddMovie(c echo.Context) error {
	return c.Render(200, "add_movie", nil)
}

func deleteMovie(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.NoContent(echo.ErrBadRequest.Code)
	}

	deleted := data.deleteMovie(id)
	if !deleted {
		return c.NoContent(echo.ErrNotFound.Code)
	}

	moviesView.refresh()
	return c.NoContent(200)
}

func postMovie(c echo.Context) error {
	rate, _ := strconv.Atoi(c.FormValue("rate"))
	year, _ := strconv.Atoi(c.FormValue("year"))
	data.addMovie(Movie{
		Title: c.FormValue("title"),
		Year:  year,
		Rate:  rate,
	})
	return c.Render(200, "movie-list", moviesView)
}

func getViewsEditMovie(c echo.Context) error {
	strId := c.Param("id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		Logger.Fatal(err)
	}

	movie, exists := data.movie(id)
	if !exists {
		Logger.Fatal(fmt.Errorf("movie with ID %d not found", id))
	}
	return c.Render(200, "edit_movie", movie)
}

func putMovie(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	rate, _ := strconv.Atoi(c.FormValue("rate"))
	year, _ := strconv.Atoi(c.FormValue("year"))
	if err := data.updateMovie(id, Movie{
		Title: c.FormValue("title"),
		Year:  year,
		Rate:  rate,
	}); err != nil {
		Logger.Error(fmt.Errorf("failed to update movie: %w", err))
		return c.NoContent(echo.ErrBadRequest.Code)
	}
	moviesView.refresh()
	return c.Render(200, "movie-list", moviesView)
}

type fieldValidation struct {
	Name    string
	Label   string
	Value   string
	Valid   bool
	Message string
}

func postValidateYear(c echo.Context) error {
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

func postValidateRate(c echo.Context) error {
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

func postValidateTitle(c echo.Context) error {
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

	movie, exists := data.movie(id)
	if !exists {
		return c.NoContent(echo.ErrBadRequest.Code)
	}

	return c.Render(200, "movie_details", movie)
}

func getMovieRow(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.NoContent(echo.ErrBadRequest.Code)
	}

	movie, exists := data.movie(id)
	if !exists {
		return c.NoContent(echo.ErrBadRequest.Code)
	}

	return c.Render(200, "movie_row", movie)
}

func getViewsData(c echo.Context) error {
	return c.Render(http.StatusOK, "data", nil)
}

func postUpload(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}

	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	movies, err := csvToMovies(src)
	if err != nil {
		return err
	}

	err = data.purge()
	if err != nil {
		return err
	}

	for _, movie := range movies {
		data.addMovie(*movie)
	}
	moviesView.refresh()

	return getViewsMovies(c)
}

func getPersonAvgRate(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	stats, err := data.personStats(id)
	if err != nil {
		return err
	}

	return c.String(http.StatusOK, strconv.FormatFloat(float64(stats.AvgRate), 'f', 2, 32))
}

func getViewsPeople(c echo.Context) error {
	return c.Render(http.StatusOK, "people-view", nil)
}

func getDownload(c echo.Context) error {
	stream, err := moviesToCSV(data.movies())
	if err != nil {
		return err
	}
	c.Response().Header().Add("Content-Disposition", "attachment")
	c.Response().Header().Add("HX-Download", "movies.csv")
	c.Response().Header().Add("Content-Type", "text/csv")
	return c.String(http.StatusOK, stream)
}

func getPersonName(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	name, exists := data.person(id)
	if !exists {
		return c.NoContent(http.StatusNotFound)
	}
	return c.String(http.StatusOK, *name)
}

func getPersonTableRow(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	if _, exists := data.person(id); !exists {
		return c.NoContent(http.StatusNotFound)
	}

	next := id + 1
	for {
		if _, exists := data.person(next); exists || next > len(data.people()) {
			break
		}
		next += 1
	}

	return c.Render(http.StatusOK, "person-row", map[string]any{
		"Curr": id,
		"Next": next,
	})
}

func getPeopleTableBody(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	people := data.people()

	lowerBound := page * limit
	if lowerBound > len(people) {
		return c.NoContent(http.StatusNoContent)
	}

	upperBound := lowerBound + limit
	if upperBound > len(people) {
		upperBound = len(people)
	}

	slices.SortFunc(people, func(i, j int) int {
		return j - i
		//if i == j {
		//	return 0
		//}
		//p1, _ := data.person(i)
		//p2, _ := data.person(j)
		//return strings.Compare(strings.ToLower(*p1), strings.ToLower(*p2))
	})

	return c.Render(http.StatusOK, "people-table-body", map[string]any{
		"People": people[lowerBound:upperBound],
		"Page":   page + 1,
		"Limit":  limit,
	})
}
