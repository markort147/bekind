package main

import (
	"strconv"
	"strings"

	"github.com/bekind/bekindfrontend/log"
	"github.com/bekind/bekindfrontend/movies"
	"github.com/labstack/echo/v4"
)

type MovieList struct {
	SortedBy string
	Desc     bool
	Body     []movies.Movie
}

func GetAllMovieHandler(c echo.Context) error {

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

func SortMovieHandler(c echo.Context) error {

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

func SimpleViewHandler(templateName string) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.Render(200, templateName, nil)
	}
}

func DeleteMovieHandler(c echo.Context) error {
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

func GetMovieHandler(c echo.Context) error {
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

func PostMovieHandler(c echo.Context) error {
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

func EditMovieHandler(c echo.Context) error {
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

func PutMovieHandler(c echo.Context) error {
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
