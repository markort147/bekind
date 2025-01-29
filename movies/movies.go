package movies

type Movies struct {
	nextId int
	Movies map[int]Movie
}

func (ms *Movies) AddMovie(m Movie) int {
	saved := ms.nextId
	ms.Movies[ms.nextId] = m
	ms.nextId++
	return saved
}

func (ms *Movies) HasId(id int) bool {
	_, exists := ms.Movies[id]
	return exists
}

func EmptyMovies() Movies {
	ms := Movies{}
	ms.Movies = make(map[int]Movie, 0)
	return ms
}
