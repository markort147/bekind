package movies

type Movies struct {
	nextId    int
	Movies    []*Movie
	MoviesMap map[int]*Movie
}

func (ms *Movies) AddMovie(m Movie) Movie {
	m.Id = ms.nextId
	ms.Movies = append(ms.Movies, &m)
	ms.MoviesMap[ms.nextId] = &m
	ms.nextId++
	return m
}

func (ms *Movies) HasId(id int) bool {
	_, exists := ms.MoviesMap[id]
	return exists
}

func EmptyMovies() Movies {
	ms := Movies{}
	SortedBy = "id"
	Desc = false
	ms.Movies = make([]*Movie, 0)
	ms.MoviesMap = make(map[int]*Movie, 0)
	return ms
}
