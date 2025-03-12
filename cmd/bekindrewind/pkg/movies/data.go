package movies

import "github.com/markort147/gopkg/log"

type Data struct {
	nextId    int
	Movies    []*Movie
	MoviesMap map[int]*Movie
}

func (ms *Data) addMovie(m Movie) Movie {
	m.Id = ms.nextId
	ms.Movies = append(ms.Movies, &m)
	ms.MoviesMap[ms.nextId] = &m
	ms.nextId++
	log.Logger.Debugf("added movie %v", m)
	return m
}

func (ms *Data) existsById(id int) bool {
	_, exists := ms.MoviesMap[id]
	return exists
}

func emptyMovies() Data {
	ms := Data{}
	ms.Movies = make([]*Movie, 0)
	ms.MoviesMap = make(map[int]*Movie)
	return ms
}
