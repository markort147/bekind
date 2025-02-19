package movies

import (
	"fmt"
	"strings"
)

/*
=== Movie Sorter ===
This file contains the sorting logic for the movies.
The sorting is done based on the field specified in the SortInfo struct.
======================
*/

// SortInfo represents the sorting information.
// It contains the field to sort by and the sorting order.
type SortInfo struct {
	SortedBy MovieField
	Desc     bool
}

var defaultSorting = SortInfo{
	SortedBy: MovieId,
	Desc:     false,
}

// CurrentSorting is the current sorting information.
// It is updated every time a new sorting is requested.
var CurrentSorting = SortInfo{
	SortedBy: MovieId,
	Desc:     false,
}

// MovieSorter is a struct that contains the sorting information and the movies to sort.
// It implements the sort.Interface interface.
type MovieSorter struct {
	SortInfo
	Movies []Movie
}

func (ms MovieSorter) Len() int      { return len(ms.Movies) }
func (ms MovieSorter) Swap(i, j int) { ms.Movies[i], ms.Movies[j] = ms.Movies[j], ms.Movies[i] }
func (ms MovieSorter) Less(i, j int) bool {
	cmp := getMovieSortingFunc(ms.SortedBy)(ms.Movies[i], ms.Movies[j])
	if ms.Desc {
		return cmp > 0
	}
	return cmp < 0
}

func getMovieSortingFunc(sortedBy MovieField) func(m1, m2 Movie) int {
	switch sortedBy {
	case MovieId:
		return func(m1, m2 Movie) int {
			return m1.Id - m2.Id
		}

	case MovieTitle:
		return func(m1, m2 Movie) int {
			return strings.Compare(m1.Title, m2.Title)
		}
	case MovieDirector:
		return func(m1, m2 Movie) int {
			return strings.Compare(m1.Director, m2.Director)
		}
	case MovieYear:
		return func(m1, m2 Movie) int {
			return strings.Compare(m1.Year, m2.Year)
		}
	default:
		panic(fmt.Errorf("unknown state: %d", sortedBy))
	}
}
