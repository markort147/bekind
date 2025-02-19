package movies

import (
	"sort"

	"github.com/markort147/bekindfrontend/internal/log"
)

/*
=== Movie Service ===
This is the service layer of the application.
It is responsible for exposing the CRUD operations on movies.
It is the only layer that has access to the data storage.
======================
*/

type MovieServiceErr string

const (
	MovieNotFoundErr = MovieServiceErr("Movie not found")
)

func (e MovieServiceErr) Error() string {
	return string(e)
}

// movies is the in-memory storage of movies
var movies Movies

// Init initializes the data with an empty list of movies
func Init() {
	movies = EmptyMovies()
}

// FillForTests fills the in-memory storage with some mock data
func FillForTests() {
	mock := []Movie{
		{Id: 3, Title: "Interstellar", Year: "2014", Director: "Christopher Nolan"},
		{Id: 1, Title: "The Godfather", Year: "1972", Director: "Francis Ford Coppola"},
		{Id: 2, Title: "Pulp Fiction", Year: "1994", Director: "Quentin Tarantino"},
		{Id: 4, Title: "Fight Club", Year: "1999", Director: "David Fincher"},
	}
	for _, m := range mock {
		movies.AddMovie(m)
	}
}

// FindAll returns all movies sorted according to the given sortInfo
func FindAll(sortInfo *SortInfo) []Movie {
	if sortInfo == nil {
		sortInfo = &defaultSorting
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

// FindByIds returns the movies with the given ids sorted according to the given sortInfo
func FindByIds(ids []int, sortInfo *SortInfo) []Movie {
	filteredMovies := make([]Movie, 0, len(ids))
	seen := make(map[int]bool)

	for _, id := range ids {
		if _, alreadyAdded := seen[id]; alreadyAdded {
			continue
		}
		if movie, exists := movies.MoviesMap[id]; exists {
			filteredMovies = append(filteredMovies, *movie)
		}
		seen[id] = true
	}

	if sortInfo == nil {
		sortInfo = &defaultSorting
	}

	sort.Sort(MovieSorter{SortInfo: *sortInfo, Movies: filteredMovies})
	CurrentSorting = *sortInfo
	log.Logger.Debugf("FindByIds. Sorted according to %v: %v", *sortInfo, filteredMovies)
	return filteredMovies
}

// Save adds the given movie to the collection
func Save(m Movie) Movie {
	return movies.AddMovie(m)
}

// Update updates the movie with the given id
func Update(id int, new Movie) {
	old := movies.MoviesMap[id]
	old.Director = new.Director
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
