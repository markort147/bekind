package movies

import (
	"math/rand/v2"

	"github.com/bekind/bekindfrontend/log"
)

var runes = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
var numbers = []rune("0123456789")

func AddRandomMovies(ms *Movies, n int) {
	for range n {
		m := RandomMovie()
        log.Logger.Debugf("Adding movie: %+v\n", m)
		ms.AddMovie(m)
	}
}

func RandomMovie() Movie {
	return Movie{
		Title:    randomString(5),
		Year:     randomNumericString(4),
		Director: randomString(10),
	}
}

func randomString(n int) string {
	s := make([]rune, n)
	for i := range n {
		s[i] = runes[rand.IntN(len(runes))]
	}
	return string(s)
}

func randomNumericString(n int) string {
	s := make([]rune, n)
	for i := range n {
		s[i] = numbers[rand.IntN(len(numbers))]
	}
	return string(s)
}
