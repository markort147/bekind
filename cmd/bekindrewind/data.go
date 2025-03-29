package main

type Data interface {
	// movie return a movie by id and a boolean indicating if it was found
	movie(int) (*Movie, bool)
	// movies return all movies
	movies() []*Movie
	// addMovie adds a movie to the data store
	addMovie(Movie) *Movie
	// deleteMovie deletes a movie from the data store
	deleteMovie(int) bool
	// updateMovie updates a movie in the data store
	updateMovie(int, Movie) error
	// purge deletes all data from the data store
	purge() error
}
