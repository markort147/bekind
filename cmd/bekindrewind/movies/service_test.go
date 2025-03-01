package movies

import (
	"testing"

	"github.com/markort147/bekind/internal/log"
)

// Helper function to generate test movies
func getTestMovies() []Movie {
	return []Movie{
		{Id: 3, Title: "Interstellar", Year: "2014", Director: "Christopher Nolan"},
		{Id: 1, Title: "The Godfather", Year: "1972", Director: "Francis Ford Coppola"},
		{Id: 2, Title: "Pulp Fiction", Year: "1994", Director: "Quentin Tarantino"},
		{Id: 4, Title: "Fight Club", Year: "1999", Director: "David Fincher"},
	}
}

func TestFindAll(t *testing.T) {
	log.Test()
	// Initialize test data
	Init()
	for _, m := range getTestMovies() {
		Save(m)
	}

	// Test sorting by Title Ascending
	sortedMovies := FindAll(&SortInfo{SortedBy: MovieTitle, Desc: false})

	if len(sortedMovies) == 0 {
		t.Fatalf("FindAll failed. Expected sortedMovies of size %d, got %d", len(getTestMovies()), len(sortedMovies))
	}

	expectedOrder := []string{"Fight Club", "Interstellar", "Pulp Fiction", "The Godfather"}

	for i, movie := range sortedMovies {
		if movie.Title != expectedOrder[i] {
			t.Errorf("FindAll failed. Expected title %s, got %s", expectedOrder[i], movie.Title)
		}
	}
}

func TestFindByIds(t *testing.T) {
	log.Test()
	// Initialize test data
	Init()
	for _, m := range getTestMovies() {
		Save(m)
	}

	// Fetch only two movies, sorted by Year Descending
	ids := []int{1, 2}
	sortedMovies := Find(FindCriteria{Id: ids})

	expectedOrder := []string{"The Godfather", "Pulp Fiction"}

	for i, movie := range sortedMovies {
		if movie.Title != expectedOrder[i] {
			t.Errorf("FindByIds failed. Expected title %s, got %s", expectedOrder[i], movie.Title)
		}
	}
}

func TestFindById(t *testing.T) {
	log.Test()
	// Initialize test data
	Init()
	for _, m := range getTestMovies() {
		Save(m)
	}

	id := 1
	movie, _ := FindById(id)
	want := "The Godfather"

	if movie.Title != want {
		t.Errorf("FindByIds failed. Expected title %s, got %s", want, movie.Title)
	}
}

func TestUpdate(t *testing.T) {
	log.Test()
	// Initialize test data
	Init()
	for _, m := range getTestMovies() {
		Save(m)
	}

	id := 1
	movie, _ := FindById(id)
	movie.Title = "The Godfather II"
	Update(id, movie)

	movie, _ = FindById(id)
	want := "The Godfather II"

	if movie.Title != want {
		t.Errorf("Update failed. Expected title %s, got %s", want, movie.Title)
	}
}

func TestDeleteById(t *testing.T) {
	t.Run("Delete existing movie", func(t *testing.T) {
		log.Test()
		// Initialize test data
		Init()
		for _, m := range getTestMovies() {
			Save(m)
		}

		id := 1
		DeleteById(id)

		_, err := FindById(id)
		want := "Requested id=1 does not exist"

		if err == nil || err != MovieNotFoundErr {
			t.Errorf("DeleteById failed. Expected error %v, got %v", want, err)
		}
	})
	t.Run("Delete non-existing movie", func(t *testing.T) {
		log.Test()
		// Initialize test data
		Init()
		for _, m := range getTestMovies() {
			Save(m)
		}

		id := 5
		ok := DeleteById(id)

		if ok {
			t.Errorf("DeleteById failed. Expected false, got true")
		}
	})

}
