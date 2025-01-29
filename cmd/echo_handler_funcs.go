package main

import (
	"strconv"
	"strings"

	"github.com/bekind/bekindfrontend/log"
	"github.com/bekind/bekindfrontend/movies"
	"github.com/labstack/echo/v4"
)

func GetAllMovieHandler(c echo.Context) error {
	ms := movies.FindAll()
	log.Logger.Debugf("movies.FindAll(): %+v\n", ms)
	return c.Render(200, "movie-list", ms)
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
	ids := strings.FieldsFunc(value, func(r rune) bool { return r == ',' })

	res := make(map[int]movies.Movie)
	for _, stringId := range ids {
		id, err := strconv.Atoi(stringId)
		if err != nil {
			continue
		}

		temp, err := movies.FindById(id)
		if err != nil {
			continue
		}

		res[id] = temp[id]
	}

	return c.Render(200, "movie-list", res)
}

func PostMovieHandler(c echo.Context) error {
	movies.Save(movies.NewMovie(
		c.FormValue("title"),
		c.FormValue("year"),
		c.FormValue("director"),
	))
	return c.NoContent(204)
}
