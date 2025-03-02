package movies

import (
	"errors"
	"strings"
)

/*
=== Movie ===
This file contains the definition of the Movie type and related functions.
================
*/

// Movie represents a movie entity.
type Movie struct {
	Id    int
	Title string
	Year  string
	Rate  uint8
}

// newMovie creates a new Movie instance.
func newMovie(title, year string, rate uint8) Movie {
	return Movie{
		Title: title,
		Year:  year,
		Rate:  rate,
	}
}

// MovieField represents a field of a Movie entity.
type MovieField int

const (
	MovieId MovieField = iota
	MovieTitle
	MovieYear
	MovieRate
)

var ErrInvalidMovieField = errors.New("invalid MovieField value")

// MFToStr returns the label for a MovieField value.
func MFToStr(field MovieField) (string, error) {
	switch field {
	case MovieId:
		return "Id", nil
	case MovieTitle:
		return "Title", nil
	case MovieYear:
		return "Year", nil
	case MovieRate:
		return "Rate", nil
	}
	return "", ErrInvalidMovieField
}

// StrToMF returns the MovieField value for a label.
// The default value is MovieId.
func StrToMF(label string) MovieField {
	switch strings.ToLower(label) {
	case "title":
		return MovieTitle
	case "year":
		return MovieYear
	case "rate":
		return MovieRate
	}
	return MovieId
}
