package main

import (
	"slices"
	"strconv"
	"strings"
)

type MoviesViewFilter struct {
	Title string
	Rate  string
	Year  string
}

type MoviesViewSorting struct {
	SortedBy MovieField
	Desc     bool
}

type MoviesView struct {
	SortInfo       MoviesViewSorting
	FilterCriteria MoviesViewFilter
	Header         map[MovieField]string
	MovieIds       []int
}

func newMoviesView() *MoviesView {
	movieIds := make([]int, 0)
	for _, movie := range data.movies() {
		movieIds = append(movieIds, movie.Id)
	}

	view := &MoviesView{
		SortInfo: MoviesViewSorting{
			SortedBy: MovieId,
			Desc:     false,
		},
		MovieIds: movieIds,
		Header:   make(map[MovieField]string),
	}
	view.refresh()

	return view
}

func (mv *MoviesView) refresh() {
	mv.refreshFilter()
	mv.refreshSorting()
	mv.refreshHeader()
}

func (mv *MoviesView) refreshHeader() {
	for _, mf := range []MovieField{MovieId, MovieTitle, MovieYear, MovieRate} {
		mv.Header[mf] = moviesViewHeaderLabel(mf, mv.SortInfo)
	}
}

func (mv *MoviesView) refreshSorting() {
	slices.SortFunc(mv.MovieIds, func(i, j int) int {
		var res int
		switch mv.SortInfo.SortedBy {
		case MovieId:
			res = i - j
		case MovieRate:
			first, _ := data.movie(i)
			second, _ := data.movie(j)
			res = first.Rate - second.Rate
		case MovieYear:
			first, _ := data.movie(i)
			second, _ := data.movie(j)
			res = first.Year - second.Year
		case MovieTitle:
			first, _ := data.movie(i)
			second, _ := data.movie(j)
			res = strings.Compare(strings.ToLower(first.Title), strings.ToLower(second.Title))
		}
		if mv.SortInfo.Desc {
			res = -res
		}
		return res
	})
}

func (mv *MoviesView) refreshFilter() {

	// rate range
	minRate := 0
	maxRate := 10
	if mv.FilterCriteria.Rate != "" {
		rates := strings.Split(mv.FilterCriteria.Rate, "-")
		if len(rates) > 0 {
			value, err := strconv.Atoi(rates[0])
			if err == nil && value >= 0 {
				minRate = value
			}
		}
		if strings.ContainsRune(mv.FilterCriteria.Rate, '-') {
			if len(rates) == 2 {
				value, err := strconv.Atoi(rates[1])
				if err == nil && value <= 10 {
					maxRate = value
				}
			}
		} else {
			maxRate = minRate
		}
	}

	// year range
	minYear := 0
	maxYear := int(^uint(0) >> 1)
	if mv.FilterCriteria.Year != "" {
		years := strings.Split(mv.FilterCriteria.Year, "-")
		if len(years) > 0 {
			value, err := strconv.Atoi(years[0])
			if err == nil && value >= minYear {
				minYear = value
			}
		}
		if strings.ContainsRune(mv.FilterCriteria.Year, '-') {
			if len(years) == 2 {
				value, err := strconv.Atoi(years[1])
				if err == nil && value <= maxYear {
					maxYear = value
				}
			}
		} else {
			maxYear = minYear
		}
	}

	final := make([]int, 0)
	for _, movie := range data.movies() {
		if (mv.FilterCriteria.Title == "" || strings.Contains(strings.ToLower(movie.Title), strings.ToLower(mv.FilterCriteria.Title))) &&
			(mv.FilterCriteria.Rate == "" || (minRate <= movie.Rate && movie.Rate <= maxRate)) &&
			(mv.FilterCriteria.Year == "" || (minYear <= movie.Year && movie.Year <= maxYear)) {
			final = append(final, movie.Id)
		}
	}
	mv.MovieIds = final
}

func (mv *MoviesView) setSortingBy(by string) {
	field := strToMF(by)
	if mv.SortInfo.SortedBy == field {
		mv.SortInfo.Desc = !mv.SortInfo.Desc
	}
	mv.SortInfo.SortedBy = field
}

func (mv *MoviesView) setFilterCriteria(criteria MoviesViewFilter) {
	mv.FilterCriteria = criteria
}

func moviesViewHeaderLabel(field MovieField, info MoviesViewSorting) string {
	res := mfToStr(field)
	if field == info.SortedBy {
		if info.Desc {
			res = res + " ↑"
		} else {
			res = res + " ↓"
		}
	}
	return res
}

func mfToStr(field MovieField) string {
	switch field {
	case MovieId:
		return "#"
	case MovieTitle:
		return "Title"
	case MovieYear:
		return "Year"
	case MovieRate:
		return "Rate"
	}
	return ""
}

func strToMF(label string) MovieField {
	switch strings.ToLower(label) {
	case "title":
		return MovieTitle
	case "year":
		return MovieYear
	case "rate":
		return MovieRate
	}
	return MovieId
}
