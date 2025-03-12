package movies

import (
	"errors"
	"testing"

	"github.com/markort147/gopkg/log"
)

func TestGetMovieFieldLabel(t *testing.T) {
	log.Test()
	testCases := []struct {
		field MovieField
		label string
		err   error
	}{
		{MovieId, "Id", nil},
		{MovieTitle, "Title", nil},
		{MovieYear, "Year", nil},
		{MovieRate, "Rate", nil},
		{MovieField(-1), "", ErrInvalidMovieField},
	}

	for _, tc := range testCases {
		label, err := MFToStr(tc.field)
		if !errors.Is(err, tc.err) {
			t.Fatalf("MFToStr(%d) returned error %v, want %v", tc.field, err, tc.err)
		}
		if label != tc.label {
			t.Errorf("MFToStr(%d) = %s, want %s", tc.field, label, tc.label)
		}
	}
}

func TestParseMovieField(t *testing.T) {
	log.Test()
	testCases := []struct {
		label string
		field MovieField
		err   error
	}{
		{"Id", MovieId, nil},
		{"Title", MovieTitle, nil},
		{"Year", MovieYear, nil},
		{"Rate", MovieRate, nil},
		{"invalid", MovieId, ErrInvalidMovieField},
	}

	for _, tc := range testCases {
		field := StrToMF(tc.label)
		if field != tc.field {
			t.Errorf("StrToMF(%s) = %d, want %d", tc.label, field, tc.field)
		}
	}
}

func TestNewMovie(t *testing.T) {
	log.Test()
	movie := Movie{Title: "The Matrix", Year: 1999, Rate: 8}
	if movie.Title != "The Matrix" {
		t.Errorf("newMovie().Title = %s, want 'The Matrix'", movie.Title)
	}
	if movie.Year != 1999 {
		t.Errorf("newMovie().Year = %d, want '1999'", movie.Year)
	}
	if movie.Rate != 8 {
		t.Errorf("newMovie().Rate = %d, want 8", movie.Rate)
	}
}
