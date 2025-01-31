package movies

import (
	"fmt"
	"strings"
)

type SortInfo struct {
	SortedBy MovieField
	Desc     bool
}

var CurrentSorting = SortInfo{
	SortedBy: MovieId,
	Desc:     false,
}

var defaultSorting = SortInfo{
	SortedBy: MovieId,
	Desc:     false,
}

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
