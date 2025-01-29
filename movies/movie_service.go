package movies

import (
	"errors"
	"strconv"

	"github.com/bekind/bekindfrontend/log"
)

var movies = EmptyMovies()

func InitWithRandoms(n int) {
	AddRandomMovies(&movies, n)
	log.Logger.Debugf("Initialized with %d movies: %+v\n", n, movies)
}

func FindAll() map[int]Movie {
	return movies.Movies
}

func FindById(id int) (map[int]Movie, error) {
	movie, exists := movies.Movies[id]
	if !exists {
		return nil, errors.New("Requested id=" + strconv.Itoa(id) + " does not exist")
	}
	return map[int]Movie{id: movie}, nil
}

func FindByIds(ids ...int) map[int]Movie {
	res := make(map[int]Movie, 0)
	for _, id := range ids {
		if movie, exists := movies.Movies[id]; exists {
			res[id] = movie
		}
	}
	return res
}

func Save(m Movie) int {
	return movies.AddMovie(m)
}

func DeleteById(id int) bool {
	if !movies.HasId(id) {
		return false
	}
	delete(movies.Movies, id)
	return true
}
