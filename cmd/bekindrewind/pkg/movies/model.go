package movies

import (
	"errors"
	"strings"
)

type Movie struct {
	Id    int
	Title string
	Year  uint16
	Rate  uint8
}
type MovieField int

const (
	MovieId MovieField = iota
	MovieTitle
	MovieYear
	MovieRate
)

var ErrInvalidMovieField = errors.New("invalid MovieField value")

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
