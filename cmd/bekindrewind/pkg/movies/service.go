package movies

import (
	"github.com/markort147/gopkg/log"
	"sort"
	"strings"
)

/*
=== Movie Service ===
This is the service layer of the application.
It is responsible for exposing the CRUD operations on movies.
It is the only layer that has access to the data storage.
======================
*/

/* ==== ERRORS ==== */

type MovieServiceErr string

const (
	MovieNotFoundErr = MovieServiceErr("Movie not found")
)

func (e MovieServiceErr) Error() string {
	return string(e)
}

// movies is the in-memory storage of movies
var movies Movies

var currMovies Movies

// CurrentSorting is the current sorting information.
// It is updated every time a new sorting is requested.
var CurrentSorting = SortInfo{
	SortedBy: MovieId,
	Desc:     false,
}

// CurrCriteria is the current search criteria.
var CurrCriteria = FindCriteria{}

// Init initializes the data with an empty list of movies
func Init() {
	movies = emptyMovies()
	currMovies = emptyMovies()
}

// FillForTests fills the in-memory storage with some mock data
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
		movies.addMovie(m)
		currMovies.addMovie(m)
	}
}

// FindAll returns all movies sorted according to the given sortInfo
func FindAll(sortInfo *SortInfo) []Movie {
	if sortInfo == nil {
		sortInfo = &SortInfo{
			SortedBy: MovieId,
			Desc:     false,
		}
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

// FindById returns the movie with the given id
func FindById(id int) (Movie, error) {
	movie, exists := movies.MoviesMap[id]
	if !exists {
		return Movie{}, MovieNotFoundErr
	}
	return *movie, nil
}

type FindCriteria struct {
	Title string
	Rate  []uint8
	Year  []uint16
}

func Find(criteria *FindCriteria, sortInfo *SortInfo) []Movie {
	if criteria != nil {
		CurrCriteria = *criteria
	}

	if sortInfo != nil {
		CurrentSorting = *sortInfo
	}

	log.Logger.Debugf("Find. Criteria: %+v, Sorting: %+v", criteria, CurrentSorting)

	// filter
	final := make([]Movie, 0)
	for _, movie := range movies.Movies {
		if (CurrCriteria.Title == "" || strings.Contains(strings.ToLower(movie.Title), strings.ToLower(CurrCriteria.Title))) &&
			(CurrCriteria.Rate == nil || (CurrCriteria.Rate[0] <= movie.Rate && movie.Rate <= CurrCriteria.Rate[1])) &&
			(CurrCriteria.Year == nil || (CurrCriteria.Year[0] <= movie.Year && movie.Year <= CurrCriteria.Year[1])) {
			final = append(final, *movie)
		}
	}

	sort.Sort(MovieSorter{SortInfo: CurrentSorting, Movies: final})
	return final
}

// Save adds the given movie to the collection
func Save(m Movie) Movie {
	return movies.addMovie(m)
}

// Update updates the movie with the given id
func Update(id int, new Movie) {
	old := movies.MoviesMap[id]
	old.Rate = new.Rate
	old.Title = new.Title
	old.Year = new.Year
}

// DeleteById deletes the movie with the given id
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
