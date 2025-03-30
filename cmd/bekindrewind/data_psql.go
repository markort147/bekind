package main

//
//import (
//	"database/sql"
//	"fmt"
//	_ "github.com/lib/pq"
//)
//
//type PsqlData struct {
//	db *sql.DB
//}
//
//
//func newPsqlData(connString string) (*PsqlData, error) {
//	pd := &PsqlData{}
//	db, err := sql.Open("postgres", connString)
//	if err != nil {
//		return nil, err
//	}
//	pd.db = db
//	return pd, nil
//}
//
//func (pd *PsqlData) movie(id int) (*Movie, bool) {
//	var exists bool
//	if err := pd.db.QueryRow("SELECT EXISTS(SELECT 1 FROM movies WHERE id=$1)", id).Scan(&exists); err != nil {
//		Logger.Errorf("Error checking movie existence: %v", err)
//	}
//	if !exists {
//		return nil, false
//	}
//
//	var title string
//	var rate int
//	var year int
//	var seen int
//	if err := pd.db.QueryRow("SELECT title, rate, year, seen FROM movies WHERE id=$1", id).Scan(id, title, year, rate, seen); err != nil {
//		Logger.Errorf("Error fetching movie: %v", err)
//		return nil, false
//	}
//
//	rows, err := pd.db.Query("SELECT person_id, role_id FROM people_movie WHERE movie_id=$1", id)
//	if err != nil {
//		Logger.Errorf("Error fetching people: %v", err)
//		return nil, false
//	}
//	defer rows.Close()
//
//	people := make([]*Person, 0)
//	for rows.Next() {
//		var personId int
//		var roleId int
//		if err := rows.Scan(&personId, &roleId); err != nil {
//			Logger.Errorf("Error scanning person: %v", err)
//			return nil, false
//		}
//
//		person := &Person{}
//		if err := pd.db.QueryRow("SELECT name FROM people WHERE id=$1", personId).Scan(&person.name); err != nil {
//			Logger.Errorf("Error fetching person name: %v", err)
//			return nil, false
//		}
//		if err := pd.db.QueryRow("SELECT name FROM roles WHERE id=$1", roleId).Scan(&person.role); err != nil {
//			Logger.Errorf("Error fetching person role: %v", err)
//			return nil, false
//		}
//
//		people = append(people, person)
//	}
//
//	movie := &Movie{
//		Id:        id,
//		Title:     title,
//		Year:      uint16(year),
//		SeenYear:  uint16(seen),
//		Rate:      uint8(rate),
//		People:    people,
//		Sagas:     nil,
//		Studios:   nil,
//		Countries: nil,
//		Genres:    nil,
//	}
//}
//
//func (pd *PsqlData) movies() []*Movie {
//	return pd.Movies
//}
//
//func (pd *PsqlData) addMovie(m Movie) *Movie {
//	m.Id = pd.nextMovieId
//	pd.Movies = append(pd.Movies, &m)
//	pd.MoviesMap[pd.nextMovieId] = &m
//	pd.nextMovieId++
//	return &m
//}
//
//func (pd *PsqlData) deleteMovie(id int) bool {
//	remove := -1
//	for i, m := range pd.Movies {
//		if m.Id == id {
//			remove = i
//			break
//		}
//	}
//
//	if remove == -1 {
//		return false
//	}
//
//	pd.Movies = append(pd.Movies[:remove], pd.Movies[remove+1:]...)
//	delete(pd.MoviesMap, id)
//
//	return true
//}
//
//func (pd *PsqlData) updateMovie(id int, new Movie) error {
//	old, exists := pd.MoviesMap[id]
//	if !exists {
//		return fmt.Errorf("movie to update with id %d not found", id)
//	}
//	old.Rate = new.Rate
//	old.Title = new.Title
//	old.Year = new.Year
//	return nil
//}
//
//func (pd *PsqlData) purge() error {
//	pd.nextMovieId = 0
//	pd.Movies = make([]*Movie, 0)
//	pd.MoviesMap = make(map[int]*Movie)
//	return nil
//}
