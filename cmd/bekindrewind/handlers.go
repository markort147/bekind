package main

import (
	"github.com/labstack/echo/v4"
	ms "github.com/markort147/bekind/cmd/bekindrewind/movies"
	"github.com/markort147/bekind/internal/log"
	"strconv"
	"strings"
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
	Body     []ms.Movie
}

// allMoviesView is a handler function that returns the "movie-list" template with the list of all ms.
// The ms are sorted by the field specified in the "sorted-by" query parameter.
// The "desc" query parameter is used to specify the sorting order (true for descending, false for ascending).
func allMoviesView(c echo.Context) error {

	sorting := ms.SortInfo{
		SortedBy: ms.StrToMF(c.QueryParam("sorted-by")),
		Desc:     c.QueryParam("desc") == "true",
	}

	sortedByLabel, err := ms.MFToStr(sorting.SortedBy)
	if err != nil {
		return err
	}

	return c.Render(200, "movie-list", MovieList{
		SortedBy: sortedByLabel,
		Desc:     sorting.Desc,
		Body:     ms.FindAll(&sorting),
	})
}

// sortMovies is a handler function that returns the "movie-list" template with the list of ms sorted by the given field.
// The field to sort by is specified in the "by" query parameter.
func sortMovies(c echo.Context) error {

	sortedBy := ms.StrToMF(c.QueryParam("by"))
	newSorting := ms.SortInfo{
		SortedBy: sortedBy,
		Desc:     ms.CurrentSorting.SortedBy == sortedBy && !ms.CurrentSorting.Desc,
	}

	movieFieldLabel, err := ms.MFToStr(newSorting.SortedBy)
	if err != nil {
		panic(err) // it should be impossible to get an error here
	}

	return c.Render(200, "movie-list", MovieList{
		SortedBy: movieFieldLabel,
		Desc:     newSorting.Desc,
		Body:     ms.FindAll(&newSorting),
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

	deleted := ms.DeleteById(id)
	if !deleted {
		return c.NoContent(echo.ErrNotFound.Code)
	}
	return c.NoContent(200)
}

// getMovie is a handler function that returns the "movie-list" template with the list of ms that match the given ids.
// The ids are specified by the "id" form value, which is a comma-separated list of movie ids.
func getMovie(c echo.Context) error {
	value := strings.ReplaceAll(c.FormValue("id"), " ", "")
	stringIds := strings.FieldsFunc(value, func(r rune) bool { return r == ',' })

	ids := make([]int, 0)
	for _, stringId := range stringIds {
		id, err := strconv.Atoi(stringId)
		if err != nil {
			continue
		}

		ids = append(ids, id)
	}

	return c.Render(200, "movie-list", MovieList{
		SortedBy: "id",
		Desc:     false,
		Body:     ms.FindByIds(ids, nil),
	})
}

// postMovie is a handler function that creates a new movie and adds it to the database.
// The movie data is specified in the form data of the request.
// The function returns the "movie-list" template with the updated list of ms.
func postMovie(c echo.Context) error {
	ms.Save(ms.Movie{
		Title:    c.FormValue("title"),
		Year:     c.FormValue("year"),
		Director: c.FormValue("director"),
	})

	return c.Render(200, "movie-list", MovieList{
		SortedBy: "id",
		Desc:     false,
		Body:     ms.FindAll(nil),
	})
}

// editMovieView is a handler function that returns the "edit_movie" template with the movie data to edit.
// The id of the movie to edit is specified in the URL path.
func editMovieView(c echo.Context) error {
	strId := c.Param("id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		log.Logger.Fatal(err)
	}

	movie, err := ms.FindById(id)
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
	ms.Update(id, ms.Movie{
		Title:    c.FormValue("title"),
		Year:     c.FormValue("year"),
		Director: c.FormValue("director"),
	})
	return c.Render(200, "movie-list", MovieList{
		SortedBy: "id",
		Desc:     false,
		Body:     ms.FindAll(nil),
	})
}

// validateYear is a helper function that validates the year format.
func validateYear(c echo.Context) error {
	year := c.FormValue("year")

	valid := true

	if len(year) != 4 {
		valid = false
	}

	value, err := strconv.Atoi(year)
	if err != nil {
		valid = false
	}

	if value < 1900 {
		valid = false
	}

	return c.Render(200, "year_input", YearValidation{
		Year:  year,
		Valid: valid,
	})
}
