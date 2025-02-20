package api

import (
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/markort147/bekind/internal/log"
	"github.com/markort147/bekind/internal/movies"
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
// The struct contains the list of movies to display, the field to sort by, and the sorting order.
type MovieList struct {
	SortedBy string
	Desc     bool
	Body     []movies.Movie
}

// GetAllMovie is a handler function that returns the "movie-list" template with the list of all movies.
// The movies are sorted by the field specified in the "sorted-by" query parameter.
// The "desc" query parameter is used to specify the sorting order (true for descending, false for ascending).
func GetAllMovie(c echo.Context) error {

	var sortedBy movies.MovieField

	sortedBy, err := movies.ParseMovieField(c.QueryParam("sorted-by"))
	if err != nil {
		sortedBy = movies.MovieId
	}

	sorting := movies.SortInfo{
		SortedBy: sortedBy,
		Desc:     c.QueryParam("desc") == "true",
	}

	sortedByLabel, err := movies.GetMovieFieldLabel(sorting.SortedBy)
	if err != nil {
		return err
	}

	return c.Render(200, "movie-list", MovieList{
		SortedBy: sortedByLabel,
		Desc:     sorting.Desc,
		Body:     movies.FindAll(&sorting),
	})
}

// SortMovie is a handler function that returns the "movie-list" template with the list of movies sorted by the given field.
// The field to sort by is specified in the "by" query parameter.
func SortMovie(c echo.Context) error {

	sortedBy, err := movies.ParseMovieField(c.QueryParam("by"))
	if err != nil {
		return err
	}

	currSorting := movies.CurrentSorting

	newSorting := movies.SortInfo{
		SortedBy: sortedBy,
		Desc:     currSorting.SortedBy == sortedBy && !currSorting.Desc,
	}

	movieFieldLabel, err := movies.GetMovieFieldLabel(newSorting.SortedBy)
	if err != nil {
		return err
	}

	return c.Render(200, "movie-list", MovieList{
		SortedBy: movieFieldLabel,
		Desc:     newSorting.Desc,
		Body:     movies.FindAll(&newSorting),
	})
}

// SimpleView is a handler function that returns a simple view with the given template name.
// The function is used to render static HTML pages that do not require any data from the server.
func SimpleView(templateName string) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.Render(200, templateName, nil)
	}
}

// DeleteMovie is a handler function that deletes a movie from the database.
// The id of the movie to delete is specified in the URL path.
// The function returns a 200 status code if the movie was deleted successfully, or a 404 status code if the movie was not found.
func DeleteMovie(c echo.Context) error {
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

// GetMovie is a handler function that returns the "movie-list" template with the list of movies that match the given ids.
// The ids are specified by the "id" form value, which is a comma-separated list of movie ids.
func GetMovie(c echo.Context) error {
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
		Body:     movies.FindByIds(ids, nil),
	})
}

// PostMovie is a handler function that creates a new movie and adds it to the database.
// The movie data is specified in the form data of the request.
// The function returns the "movie-list" template with the updated list of movies.
func PostMovie(c echo.Context) error {
	movies.Save(movies.NewMovie(
		c.FormValue("title"),
		c.FormValue("year"),
		c.FormValue("director"),
	))

	return c.Render(200, "movie-list", MovieList{
		SortedBy: "id",
		Desc:     false,
		Body:     movies.FindAll(nil),
	})
}

// EditMovie is a handler function that returns the "edit_movie" template with the movie data to edit.
// The id of the movie to edit is specified in the URL path.
func EditMovie(c echo.Context) error {
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

// PutMovie is a handler function that updates the movie data in the database.
// The id of the movie to update is specified in the URL path.
// The new movie data is specified in the form data of the request.
// The function returns the "movie-list" template with the updated list of movies.
func PutMovie(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	movies.Update(id, movies.NewMovie(
		c.FormValue("title"),
		c.FormValue("year"),
		c.FormValue("director"),
	))
	return c.Render(200, "movie-list", MovieList{
		SortedBy: "id",
		Desc:     false,
		Body:     movies.FindAll(nil),
	})
}
