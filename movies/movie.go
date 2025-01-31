package movies

import (
	"fmt"
	"strings"
)

type MovieField int

const (
	MovieId MovieField = iota
	MovieTitle
	MovieYear
	MovieDirector
)

func GetMovieFieldLabel(field MovieField) (string, error) {
	switch field {
	case MovieId:
		return "Id", nil
	case MovieTitle:
		return "Title", nil
	case MovieYear:
		return "Year", nil
	case MovieDirector:
		return "Director", nil
	}
	return "", fmt.Errorf("invalid MovieField value: %d", field)
}

func ParseMovieField(label string) (MovieField, error) {
	switch strings.ToLower(label) {
	case "id":
		return MovieId, nil
	case "title":
		return MovieTitle, nil
	case "year":
		return MovieYear, nil
	case "director":
		return MovieDirector, nil
	}
	return -1, fmt.Errorf("invalid MovieField value: %s", label)
}

type Movie struct {
	Id       int
	Title    string
	Year     string
	Director string
}

func NewMovie(title, year, director string) Movie {
	return Movie{
		Id:       -1,
		Title:    title,
		Year:     year,
		Director: director,
	}
}

func EmptyMovie() Movie {
	return Movie{}
}
