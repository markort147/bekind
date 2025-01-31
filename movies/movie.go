package movies

type Movie struct {
	Id       int
	Title    string
	Year     string
	Director string
}

func NewMovie(title, year, director string) Movie {
	return Movie{
		Id:       -1,
		Title:    title,
		Year:     year,
		Director: director,
	}
}

func EmptyMovie() Movie {
	return Movie{}
}
