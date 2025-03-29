package main

import (
	"fmt"
)

type MemoryData struct {
	nextId    int
	Movies    []*Movie
	MoviesMap map[int]*Movie
}

func newMemoryData() *MemoryData {
	md := &MemoryData{}
	md.nextId = 0
	md.Movies = make([]*Movie, 0)
	md.MoviesMap = make(map[int]*Movie)
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
	m.Id = md.nextId
	md.Movies = append(md.Movies, &m)
	md.MoviesMap[md.nextId] = &m
	md.nextId++
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
	md.nextId = 0
	md.Movies = make([]*Movie, 0)
	md.MoviesMap = make(map[int]*Movie)
	return nil
}
