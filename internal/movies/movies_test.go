package movies

import (
	"testing"
)

func TestAddMovie(t *testing.T) {
	// Create a new Movies struct
	ms := EmptyMovies()

	// Add a movie to the Movies struct
	m := NewMovie("The Matrix", "1999", "The Wachowskis")
	ms.AddMovie(m)

	// Check if the movie was added
	if len(ms.Movies) != 1 {
		t.Errorf("expected 1 movie, got %d", len(ms.Movies))
	}

	// Check if the movie was added correctly
	got := ms.Movies[0]
	if got.Title != m.Title || got.Year != m.Year || got.Director != m.Director {
		t.Errorf("expected %v, got %v", m, got)
	}
}

func TestHasId(t *testing.T) {
	// Create a new Movies struct
	ms := EmptyMovies()

	// Add a movie to the Movies struct
	m := NewMovie("The Matrix", "1999", "The Wachowskis")
	ms.AddMovie(m)

	// Check if the movie exists
	if !ms.ExistsById(0) {
		t.Errorf("expected movie with id %d to exist", m.Id)
	}

	// Check if a non-existent movie does not exist
	if ms.ExistsById(1) {
		t.Errorf("expected movie with id %d to not exist", m.Id+1)
	}
}
