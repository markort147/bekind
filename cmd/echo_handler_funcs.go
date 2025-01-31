package main

import (
	"strconv"
	"strings"

	"github.com/bekind/bekindfrontend/log"
	"github.com/bekind/bekindfrontend/movies"
	"github.com/labstack/echo/v4"
)

func GetAllMovieHandler(c echo.Context) error {

	sortedBy := c.QueryParam("sorted-by")
	if len(sortedBy) == 0 {
		sortedBy = "id"
	}

	desc := c.QueryParam("desc") == "true"

	// log.Logger.Debugf("movies.FindAll(): %+v\n", ms)
	ms := movies.FindAllSorted(sortedBy, desc)

	return c.Render(200, "movie-list", struct {
		SortedBy string
		Desc     bool
		Body     []movies.Movie
	}{
		SortedBy: sortedBy,
		Desc:     desc,
		Body:     ms,
	})
}

func SortMovieHandler(c echo.Context) error {

	sortedBy := c.QueryParam("by")

	desc := movies.SortedBy != sortedBy || (movies.SortedBy == sortedBy && !movies.Desc)

	// log.Logger.Debugf("movies.FindAll(): %+v\n", ms)
	ms := movies.FindAllSorted(sortedBy, desc)

	return c.Render(200, "movie-list", struct {
		SortedBy string
		Desc     bool
		Body     []movies.Movie
	}{
		SortedBy: sortedBy,
		Desc:     desc,
		Body:     ms,
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

	return c.Render(200, "movie-list", struct {
		SortedBy string
		Desc     bool
		Body     []movies.Movie
	}{
		SortedBy: "id",
		Desc:     false,
		Body:     movies.FindByIds(ids...),
	})
}

func PostMovieHandler(c echo.Context) error {
	movies.Save(movies.NewMovie(
		c.FormValue("title"),
		c.FormValue("year"),
		c.FormValue("director"),
	))

	return c.Render(200, "movie-list", struct {
		SortedBy string
		Desc     bool
		Body     []movies.Movie
	}{
		SortedBy: "id",
		Desc:     false,
		Body:     movies.FindAllSorted("id", false),
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
	return c.Render(200, "movie-list", struct {
		SortedBy string
		Desc     bool
		Body     []movies.Movie
	}{
		SortedBy: "id",
		Desc:     false,
		Body:     movies.FindAllSorted("id", false),
	})
}
