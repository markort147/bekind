package movies

import (
	"reflect"
	"sort"
	"testing"

	"github.com/bekind/bekindfrontend/log"
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

// ✅ **1. Test Sorting Logic with `MovieSorter`**
func TestMovieSorter(t *testing.T) {
	log.Test()
	testCases := []struct {
		field  MovieField
		desc   bool
		output []int // Expected order of IDs
	}{
		{MovieId, false, []int{1, 2, 3, 4}},       // Sort by ID Ascending
		{MovieId, true, []int{4, 3, 2, 1}},        // Sort by ID Descending
		{MovieTitle, false, []int{4, 3, 2, 1}},    // Sort by Title Ascending
		{MovieTitle, true, []int{1, 2, 3, 4}},     // Sort by Title Descending
		{MovieYear, false, []int{1, 2, 4, 3}},     // Sort by Year Ascending
		{MovieYear, true, []int{3, 4, 2, 1}},      // Sort by Year Descending
		{MovieDirector, false, []int{3, 4, 1, 2}}, // Sort by Director Ascending
		{MovieDirector, true, []int{2, 1, 4, 3}},  // Sort by Director Descending
	}

	for _, tc := range testCases {
		// Make a copy of test movies
		movies := getTestMovies()
		sort.Sort(MovieSorter{SortInfo: SortInfo{SortedBy: tc.field, Desc: tc.desc}, Movies: movies})

		// Extract sorted IDs
		sortedIDs := []int{}
		for _, m := range movies {
			sortedIDs = append(sortedIDs, m.Id)
		}

		// Check if sorted correctly
		if !reflect.DeepEqual(sortedIDs, tc.output) {
			sortingField, _ := GetMovieFieldLabel(tc.field)
			t.Fatalf("Sorting by %v (Desc: %v) failed. Got %v, expected %v",
				sortingField, tc.desc, sortedIDs, tc.output)
		}
	}
}

// ✅ **2. Test `FindAll()`**
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

// ✅ **3. Test `FindByIds()`**
func TestFindByIds(t *testing.T) {
	log.Test()
	// Initialize test data
	Init()
	for _, m := range getTestMovies() {
		Save(m)
	}

	// Fetch only two movies, sorted by Year Descending
	ids := []int{1, 2}
	sortedMovies := FindByIds(ids, &SortInfo{SortedBy: MovieYear, Desc: true})

	expectedOrder := []string{"Pulp Fiction", "The Godfather"}

	for i, movie := range sortedMovies {
		if movie.Title != expectedOrder[i] {
			t.Errorf("FindByIds failed. Expected title %s, got %s", expectedOrder[i], movie.Title)
		}
	}
}
