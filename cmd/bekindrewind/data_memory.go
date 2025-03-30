package main

import (
	"fmt"
)

type MemoryData struct {
	nextMovieId int
	Movies      []*Movie
	MoviesMap   map[int]*Movie

	nextPersonId int
	People       map[string]struct{}
	PeopleMap    map[int]*string
}

func newMemoryData() *MemoryData {
	md := &MemoryData{}
	md.nextMovieId = 0
	md.Movies = make([]*Movie, 0)
	md.MoviesMap = make(map[int]*Movie)

	md.nextPersonId = 0
	md.People = make(map[string]struct{})
	md.PeopleMap = make(map[int]*string)
	return md
}

func (md *MemoryData) movie(id int) (*Movie, bool) {
	movie, exists := md.MoviesMap[id]
	return movie, exists
}

func (md *MemoryData) movies() []*Movie {
	return md.Movies
}

func (md *MemoryData) addMovie(m Movie) *Movie {
	m.Id = md.nextMovieId
	md.Movies = append(md.Movies, &m)
	md.MoviesMap[md.nextMovieId] = &m
	md.nextMovieId++

	for name := range m.People {
		if _, exists := md.People[*name]; !exists {
			md.People[*name] = struct{}{}
			md.PeopleMap[md.nextPersonId] = name
			md.nextPersonId++
		}
	}
	return &m
}

func (md *MemoryData) deleteMovie(id int) bool {
	remove := -1
	for i, m := range md.Movies {
		if m.Id == id {
			remove = i
			break
		}
	}

	if remove == -1 {
		return false
	}

	md.Movies = append(md.Movies[:remove], md.Movies[remove+1:]...)
	delete(md.MoviesMap, id)

	return true
}

func (md *MemoryData) updateMovie(id int, new Movie) error {
	old, exists := md.MoviesMap[id]
	if !exists {
		return fmt.Errorf("movie to update with id %d not found", id)
	}
	old.Rate = new.Rate
	old.Title = new.Title
	old.Year = new.Year
	return nil
}

func (md *MemoryData) purge() error {
	md.nextMovieId = 0
	md.Movies = make([]*Movie, 0)
	md.MoviesMap = make(map[int]*Movie)
	return nil
}

func (md *MemoryData) people() []int {
	peopleList := make([]int, 0)
	for p := range md.PeopleMap {
		peopleList = append(peopleList, p)
	}
	return peopleList
}

func (md *MemoryData) personStats(id int) (*PersonStats, error) {
	person := md.PeopleMap[id]
	movies := make([]*Movie, 0)
	for _, m := range md.Movies {
		for p := range m.People {
			if *person == *p {
				movies = append(movies, m)
				break
			}
		}
	}

	rate := float32(0)
	if len(movies) > 0 {
		for _, m := range movies {
			rate += float32(m.Rate)
		}
		rate /= float32(len(movies))
	}

	return &PersonStats{
		AvgRate: rate,
	}, nil
}

func (md *MemoryData) person(id int) (*string, bool) {
	if person, exists := md.PeopleMap[id]; exists {
		return person, true
	}
	return nil, false
}
