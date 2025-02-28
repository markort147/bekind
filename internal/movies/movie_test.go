package movies

import (
	"testing"

	"github.com/markort147/bekind/internal/log"
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
		{MovieDirector, "Director", nil},
		{MovieField(-1), "", ErrInvalidMovieField},
	}

	for _, tc := range testCases {
		label, err := MFToStr(tc.field)
		if err != tc.err {
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
		{"Director", MovieDirector, nil},
		{"invalid", -1, ErrInvalidMovieField},
	}

	for _, tc := range testCases {
		field, err := StrToMF(tc.label)
		if err != tc.err {
			t.Fatalf("StrToMF(%s) returned error %v, want %v", tc.label, err, tc.err)
		}
		if field != tc.field {
			t.Errorf("StrToMF(%s) = %d, want %d", tc.label, field, tc.field)
		}
	}
}

func TestNewMovie(t *testing.T) {
	log.Test()
	movie := Movie{Title: "The Matrix", Year: "1999", Director: "Lana Wachowski, Lilly Wachowski"}
	if movie.Id != -1 {
		t.Errorf("newMovie().Id = %d, want -1", movie.Id)
	}
	if movie.Title != "The Matrix" {
		t.Errorf("newMovie().Title = %s, want 'The Matrix'", movie.Title)
	}
	if movie.Year != "1999" {
		t.Errorf("newMovie().Year = %s, want '1999'", movie.Year)
	}
	if movie.Director != "Lana Wachowski, Lilly Wachowski" {
		t.Errorf("newMovie().Director = %s, want 'Lana Wachowski, Lilly Wachowski'", movie.Director)
	}
}
