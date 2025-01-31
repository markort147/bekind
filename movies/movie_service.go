package movies

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/bekind/bekindfrontend/log"
)


var SortedBy = "id"
var Desc = false
var movies = EmptyMovies()

func getMovieSortingFunc(sortedBy string) func(m1, m2 Movie) int {
	if sortedBy == "id" {
		return func(m1, m2 Movie) int {
			return m1.Id - m2.Id
		}
	}
	if sortedBy == "title" {
		return func(m1, m2 Movie) int {
			return strings.Compare(m1.Title, m2.Title)
		}
	}
	if sortedBy == "director" {
		return func(m1, m2 Movie) int {
			return strings.Compare(m1.Director, m2.Director)
		}
	}
	if sortedBy == "year" {
		return func(m1, m2 Movie) int {
			return strings.Compare(m1.Year, m2.Year)
		}
	}
	panic(fmt.Errorf("value of sortedBy=%s not recognized", sortedBy))
}

func InitWithRandoms(n int) {
	AddRandomMovies(&movies, n)
	log.Logger.Debugf("Initialized with %d movies: %+v\n", n, movies)
}

func FindAll() []Movie {
	return FindAllSorted("id", false)
}

func FindAllSorted(sortedBy string, desc bool) []Movie {
	sort(sortedBy, desc)
	res := make([]Movie, 0)
	for _, m := range movies.Movies {
		res = append(res, *m)
	}
	return res
}

func sort(sortedBy string, desc bool) {
	slices.SortFunc(movies.Movies, func(m1, m2 *Movie) int {
		out := getMovieSortingFunc(sortedBy)(*m1, *m2)
		if desc {
			return -out
		}
		return out
	})
	SortedBy = sortedBy
	Desc = desc
}

func FindById(id int) (Movie, error) {
	movie, exists := movies.MoviesMap[id]
	if !exists {
		return Movie{}, errors.New("Requested id=" + strconv.Itoa(id) + " does not exist")
	}
	return *movie, nil
}

func FindByIds(ids ...int) []Movie {
	return FindByIdsSorted("id", false, ids...)
}

func FindByIdsSorted(sortedBy string, desc bool, ids ...int) []Movie {
	sort(sortedBy, desc)
	res := make([]Movie, 0)
	seen := make(map[int]bool)
	for _, id := range ids {
		if _, alreadyAdded := seen[id]; alreadyAdded {
			continue
		}
		if movie, exists := movies.MoviesMap[id]; exists {
			res = append(res, *movie)
		}
		seen[id] = true
	}
	return res
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
