package movies

type MovieServiceErr string

const (
	MovieNotFoundErr = MovieServiceErr("Movie not found")
)

func (e MovieServiceErr) Error() string {
	return string(e)
}

var data Data

func Init() {
	data = emptyMovies()
	CurrView.initView(&data)
}

func FillForTests() {
	mock := []Movie{
		{Id: 3, Title: "Interstellar", Year: 2014, Rate: 8},
		{Id: 1, Title: "The Godfather", Year: 1972, Rate: 7},
		{Id: 2, Title: "Pulp Fiction", Year: 1994, Rate: 9},
		{Id: 4, Title: "Fight Club", Year: 1999, Rate: 9},
		{Id: 5, Title: "The Shawshank Redemption", Year: 1994, Rate: 10},
		{Id: 6, Title: "The Dark Knight", Year: 2008, Rate: 9},
		{Id: 7, Title: "Inception", Year: 2010, Rate: 8},
		{Id: 8, Title: "The Matrix", Year: 1999, Rate: 7},
		{Id: 9, Title: "The Lord of the Rings: The Return of the King", Year: 2003, Rate: 8},
		{Id: 10, Title: "The Lord of the Rings: The Fellowship of the Ring", Year: 2001, Rate: 8},
		{Id: 11, Title: "The Lord of the Rings: The Two Towers", Year: 2002, Rate: 8},
		{Id: 12, Title: "Star Wars: Episode IV - A New Hope", Year: 1977, Rate: 7},
		{Id: 13, Title: "Star Wars: Episode V - The Empire Strikes Back", Year: 1980, Rate: 8},
		{Id: 14, Title: "Star Wars: Episode VI - Return of the Jedi", Year: 1983, Rate: 7},
		{Id: 15, Title: "Star Wars: Episode I - The Phantom Menace", Year: 1999, Rate: 6},
		{Id: 16, Title: "Star Wars: Episode II - Attack of the Clones", Year: 2002, Rate: 5},
		{Id: 17, Title: "Star Wars: Episode III - Revenge of the Sith", Year: 2005, Rate: 4},
		{Id: 18, Title: "Star Wars: Episode VII - The Force Awakens", Year: 2015, Rate: 6},
		{Id: 19, Title: "Star Wars: Episode VIII - The Last Jedi", Year: 2017, Rate: 5},
		{Id: 20, Title: "Star Wars: Episode IX - The Rise of Skywalker", Year: 2019, Rate: 4},
		{Id: 21, Title: "The Avengers", Year: 2012, Rate: 8},
		{Id: 22, Title: "Avengers: Age of Ultron", Year: 2015, Rate: 7},
		{Id: 23, Title: "Avengers: Infinity War", Year: 2018, Rate: 8},
		{Id: 24, Title: "Avengers: Endgame", Year: 2019, Rate: 8},
		{Id: 25, Title: "The Dark Knight Rises", Year: 2012, Rate: 8},
		{Id: 26, Title: "The Hobbit: An Unexpected Journey", Year: 2012, Rate: 7},
		{Id: 27, Title: "The Hobbit: The Desolation of Smaug", Year: 2013, Rate: 6},
		{Id: 28, Title: "The Hobbit: The Battle of the Five Armies", Year: 2014, Rate: 5},
		{Id: 29, Title: "The Hunger Games", Year: 2012, Rate: 6},
		{Id: 30, Title: "The Hunger Games: Catching Fire", Year: 2013, Rate: 5},
		{Id: 31, Title: "The Hunger Games: Mockingjay - Part 1", Year: 2014, Rate: 4},
		{Id: 32, Title: "The Hunger Games: Mockingjay - Part 2", Year: 2015, Rate: 3},
		{Id: 33, Title: "The Twilight Saga: Breaking Dawn - Part 1", Year: 2011, Rate: 4},
		{Id: 34, Title: "The Twilight Saga: Breaking Dawn - Part 2", Year: 2012, Rate: 3},
		{Id: 35, Title: "The Twilight Saga: Eclipse", Year: 2010, Rate: 5},
		{Id: 36, Title: "The Twilight Saga: New Moon", Year: 2009, Rate: 6},
		{Id: 37, Title: "The Twilight Saga: Twilight", Year: 2008, Rate: 7},
		{Id: 38, Title: "Iron Man", Year: 2008, Rate: 8},
		{Id: 39, Title: "Iron Man 2", Year: 2010, Rate: 7},
		{Id: 40, Title: "Iron Man 3", Year: 2013, Rate: 6},
		{Id: 41, Title: "Thor", Year: 2011, Rate: 7},
		{Id: 42, Title: "Thor: The Dark World", Year: 2013, Rate: 6},
		{Id: 43, Title: "Thor: Ragnarok", Year: 2017, Rate: 8},
	}
	for _, m := range mock {
		data.addMovie(m)
	}
}

func FindById(id int) (Movie, error) {
	movie, exists := data.MoviesMap[id]
	if !exists {
		return Movie{}, MovieNotFoundErr
	}
	return *movie, nil
}

func Save(m Movie) Movie {
	saved := data.addMovie(m)
	CurrView.refresh()
	return saved
}

func Update(id int, new Movie) {
	old := data.MoviesMap[id]
	old.Rate = new.Rate
	old.Title = new.Title
	old.Year = new.Year
	CurrView.refresh()
}

func Delete(id int) bool {
	remove := -1
	for i, m := range data.Movies {
		if m.Id == id {
			remove = i
			break
		}
	}

	if remove == -1 {
		return false
	}

	data.Movies = append(data.Movies[:remove], data.Movies[remove+1:]...)
	delete(data.MoviesMap, id)

	CurrView.refresh()
	return true
}
