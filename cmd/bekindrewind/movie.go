package main

type Movie struct {
	Id        int
	Title     string
	Sagas     []string
	Year      uint16
	SeenYear  uint16
	Rate      uint8
	Directors []string
	Writers   []string
	Composers []string
	Dops      []string
	Editors   []string
	Producers []string
	Studios   []string
	Countries []string
	Genres    []string
}
type MovieField int

const (
	MovieId MovieField = iota
	MovieTitle
	MovieYear
	MovieRate
)
