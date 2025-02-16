package movies

import (
	"testing"

	"github.com/bekind/bekindfrontend/log"
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
		label, err := GetMovieFieldLabel(tc.field)
		if err != tc.err {
			t.Fatalf("GetMovieFieldLabel(%d) returned error %v, want %v", tc.field, err, tc.err)
		}
		if label != tc.label {
			t.Errorf("GetMovieFieldLabel(%d) = %s, want %s", tc.field, label, tc.label)
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
		field, err := ParseMovieField(tc.label)
		if err != tc.err {
			t.Fatalf("ParseMovieField(%s) returned error %v, want %v", tc.label, err, tc.err)
		}
		if field != tc.field {
			t.Errorf("ParseMovieField(%s) = %d, want %d", tc.label, field, tc.field)
		}
	}
}

func TestNewMovie(t *testing.T) {
	log.Test()
	movie := NewMovie("The Matrix", "1999", "Lana Wachowski, Lilly Wachowski")
	if movie.Id != -1 {
		t.Errorf("NewMovie().Id = %d, want -1", movie.Id)
	}
	if movie.Title != "The Matrix" {
		t.Errorf("NewMovie().Title = %s, want 'The Matrix'", movie.Title)
	}
	if movie.Year != "1999" {
		t.Errorf("NewMovie().Year = %s, want '1999'", movie.Year)
	}
	if movie.Director != "Lana Wachowski, Lilly Wachowski" {
		t.Errorf("NewMovie().Director = %s, want 'Lana Wachowski, Lilly Wachowski'", movie.Director)
	}
}
