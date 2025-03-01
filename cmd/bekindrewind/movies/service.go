package movies

import (
	"github.com/markort147/bekind/internal/log"
	"sort"
	"strings"
)

/*
=== Movie Service ===
This is the service layer of the application.
It is responsible for exposing the CRUD operations on movies.
It is the only layer that has access to the data storage.
======================
*/

/* ==== ERRORS ==== */

type MovieServiceErr string

const (
	MovieNotFoundErr = MovieServiceErr("Movie not found")
)

func (e MovieServiceErr) Error() string {
	return string(e)
}

// movies is the in-memory storage of movies
var movies Movies

var currMovies Movies

// CurrentSorting is the current sorting information.
// It is updated every time a new sorting is requested.
var CurrentSorting = SortInfo{
	SortedBy: MovieId,
	Desc:     false,
}

// CurrCriteria is the current search criteria.
var CurrCriteria = FindCriteria{}

// Init initializes the data with an empty list of movies
func Init() {
	movies = emptyMovies()
	currMovies = emptyMovies()
}

// FillForTests fills the in-memory storage with some mock data
func FillForTests() {
	mock := []Movie{
		{Id: 3, Title: "Interstellar", Year: "2014", Rate: 8},
		{Id: 1, Title: "The Godfather", Year: "1972", Rate: 7},
		{Id: 2, Title: "Pulp Fiction", Year: "1994", Rate: 9},
		{Id: 4, Title: "Fight Club", Year: "1999", Rate: 9},
	}
	for _, m := range mock {
		movies.addMovie(m)
		currMovies.addMovie(m)
	}
}

// FindAll returns all movies sorted according to the given sortInfo
func FindAll(sortInfo *SortInfo) []Movie {
	if sortInfo == nil {
		sortInfo = &SortInfo{
			SortedBy: MovieId,
			Desc:     false,
		}
	}

	sortedMovies := make([]Movie, len(movies.Movies))
	for i, m := range movies.Movies {
		sortedMovies[i] = *m
	}

	sort.Sort(MovieSorter{SortInfo: *sortInfo, Movies: sortedMovies})
	CurrentSorting = *sortInfo
	log.Logger.Debugf("FindAll. Sorted according to %v: %v", *sortInfo, sortedMovies)
	return sortedMovies
}

// FindById returns the movie with the given id
func FindById(id int) (Movie, error) {
	movie, exists := movies.MoviesMap[id]
	if !exists {
		return Movie{}, MovieNotFoundErr
	}
	return *movie, nil
}

type FindCriteria struct {
	Title string
	Rate  []uint8
}

func Find(criteria *FindCriteria, sortInfo *SortInfo) []Movie {
	if criteria != nil {
		CurrCriteria = *criteria
	}

	if sortInfo != nil {
		CurrentSorting = *sortInfo
	}

	// filter
	final := make([]Movie, 0)
	for _, movie := range movies.Movies {
		if (CurrCriteria.Title == "" || strings.Contains(strings.ToLower(movie.Title), strings.ToLower(CurrCriteria.Title))) &&
			(CurrCriteria.Rate == nil || (CurrCriteria.Rate[0] <= movie.Rate && movie.Rate <= CurrCriteria.Rate[1])) {
			final = append(final, *movie)
		}
	}

	sort.Sort(MovieSorter{SortInfo: CurrentSorting, Movies: final})
	return final
}

// Save adds the given movie to the collection
func Save(m Movie) Movie {
	return movies.addMovie(m)
}

// Update updates the movie with the given id
func Update(id int, new Movie) {
	old := movies.MoviesMap[id]
	old.Rate = new.Rate
	old.Title = new.Title
	old.Year = new.Year
}

// DeleteById deletes the movie with the given id
func DeleteById(id int) bool {
	remove := -1
	for i, m := range movies.Movies {
		if m.Id == id {
			remove = i
			break
		}
	}

	if remove == -1 {
		return false
	}

	movies.Movies = append(movies.Movies[:remove], movies.Movies[remove+1:]...)
	delete(movies.MoviesMap, id)
	return true
}
