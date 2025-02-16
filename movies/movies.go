package movies

import "github.com/bekind/bekindfrontend/log"

/*
=== Movies ===
This file contains the definition of the Movies struct and related functions.
The Movies struct is a collection of Movie structs. It is used to store all the movies in the system.
================
*/

// Movies defines the data structure used to store all the movies in the system.
// It contains a slice of Movie pointers and a map of Movie pointers.
// The slice is used to store the order of the movies, while the map is used to quickly access a movie by its id.
// The nextId field is used to assign a unique id to each movie.
type Movies struct {
	nextId    int
	Movies    []*Movie
	MoviesMap map[int]*Movie
}

// AddMovie adds a movie to the Movies struct.
// It assigns a unique id to the movie and adds it to the slice and map.
// It returns the added movie.
func (ms *Movies) AddMovie(m Movie) Movie {
	m.Id = ms.nextId
	ms.Movies = append(ms.Movies, &m)
	ms.MoviesMap[ms.nextId] = &m
	ms.nextId++
	log.Logger.Infof("added movie %v", m)
	return m
}

// ExistsById checks if a movie with the given id exists in the Movies struct.
// It returns true if the movie exists, false otherwise.
func (ms *Movies) ExistsById(id int) bool {
	_, exists := ms.MoviesMap[id]
	return exists
}

// EmptyMovies creates an empty Movies struct and returns it.
// It initializes the Movies struct with an empty slice and map.
func EmptyMovies() Movies {
	ms := Movies{}
	ms.Movies = make([]*Movie, 0)
	ms.MoviesMap = make(map[int]*Movie, 0)
	return ms
}
