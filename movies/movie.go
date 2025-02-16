package movies

import (
	"fmt"
	"strings"
)

/*
=== Movie ===
This file contains the definition of the Movie type and related functions.
================
*/

type MovieField int

const (
	MovieId MovieField = iota
	MovieTitle
	MovieYear
	MovieDirector
)

// GetMovieFieldLabel returns the label for a MovieField value.
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

// Movie represents a movie entity.
type Movie struct {
	Id       int
	Title    string
	Year     string
	Director string
}

// ParseMovieField returns the MovieField value for a label.
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

// NewMovie creates a new Movie instance.
func NewMovie(title, year, director string) Movie {
	return Movie{
		Id:       -1,
		Title:    title,
		Year:     year,
		Director: director,
	}
}

// EmptyMovie returns an empty Movie instance.
func EmptyMovie() Movie {
	return Movie{}
}
