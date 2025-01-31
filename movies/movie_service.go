package movies

import (
	"errors"
	"sort"
	"strconv"

	"github.com/bekind/bekindfrontend/log"
)

var movies Movies

func Init() {
	movies = EmptyMovies()
}

func FillWithRandoms(n int) {
	AddRandomMovies(&movies, n)
	log.Logger.Debugf("Initialized with %d movies: %+v\n", n, movies)
}

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

func FindById(id int) (Movie, error) {
	movie, exists := movies.MoviesMap[id]
	if !exists {
		return Movie{}, errors.New("Requested id=" + strconv.Itoa(id) + " does not exist")
	}
	return *movie, nil
}

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

func Save(m Movie) Movie {
	return movies.AddMovie(m)
}

func Update(id int, new Movie) {
	old := movies.MoviesMap[id]
	old.Director = new.Director
	old.Title = new.Title
	old.Year = new.Year
}

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
