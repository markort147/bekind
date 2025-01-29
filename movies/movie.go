package movies

type Movie struct {
	Title    string
	Year     string
	Director string
}

func NewMovie(title, year, director string) Movie {
	return Movie{
		Title:    title,
		Year:     year,
		Director: director,
	}
}

func EmptyMovie() Movie {
	return Movie{}
}