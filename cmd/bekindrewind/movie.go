package main

type Movie struct {
	Id        int
	Title     string
	Sagas     []string
	Year      int
	SeenYear  int
	Rate      int
	People    map[*string][]string
	Studios   []string
	Countries []string
	Genres    []string
}

type PersonStats struct {
	AvgRate float32
}

type MovieField = string

const (
	MovieId    MovieField = "id"
	MovieTitle            = "title"
	MovieYear             = "year"
	MovieRate             = "rate"
)
