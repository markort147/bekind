package main

import (
	"slices"
	"strconv"
	"strings"
)

type MoviesViewHeader struct {
	Id    string
	Title string
	Year  string
	Rate  string
}

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
	Header         MoviesViewHeader
	MovieIds       []int
	data           *Data
	FilterCriteria MoviesViewFilter
}

func newMoviesView() *MoviesView {
	movieIds := make([]int, 0)
	for _, movie := range data.movies() {
		movieIds = append(movieIds, movie.Id)
	}

	view := &MoviesView{
		data: &data,
		SortInfo: MoviesViewSorting{
			SortedBy: MovieId,
			Desc:     false,
		},
		MovieIds: movieIds,
	}
	view.refresh()

	return view
}

func (mv *MoviesView) refresh() {
	mv.refreshFilter()
	mv.refreshSorting()
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
			res = int(first.Rate) - int(second.Rate)
		case MovieYear:
			first, _ := data.movie(i)
			second, _ := data.movie(j)
			res = int(first.Year) - int(second.Year)
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

	mv.Header.Id = moviesViewHeaderLabel(MovieId, mv.SortInfo)
	mv.Header.Title = moviesViewHeaderLabel(MovieTitle, mv.SortInfo)
	mv.Header.Year = moviesViewHeaderLabel(MovieYear, mv.SortInfo)
	mv.Header.Rate = moviesViewHeaderLabel(MovieRate, mv.SortInfo)
}

func (mv *MoviesView) refreshFilter() {
	title := mv.FilterCriteria.Title

	// rate range
	var rate []uint8
	if mv.FilterCriteria.Rate != "" {
		minRate := uint8(0)
		maxRate := uint8(10)
		rates := strings.Split(mv.FilterCriteria.Rate, "-")
		if len(rates) > 0 {
			value, err := strconv.Atoi(rates[0])
			if err == nil && value >= 0 {
				minRate = uint8(value)
			}
		}
		if strings.ContainsRune(mv.FilterCriteria.Rate, '-') {
			if len(rates) == 2 {
				value, err := strconv.Atoi(rates[1])
				if err == nil && value <= 10 {
					maxRate = uint8(value)
				}
			}
		} else {
			maxRate = minRate
		}
		rate = []uint8{minRate, maxRate}
	}

	// year range
	var year []uint16
	if mv.FilterCriteria.Year != "" {
		minYear := uint16(0)
		maxYear := ^uint16(0)
		years := strings.Split(mv.FilterCriteria.Year, "-")
		if len(years) > 0 {
			value, err := strconv.Atoi(years[0])
			if err == nil && value >= int(minYear) {
				minYear = uint16(value)
			}
		}
		if strings.ContainsRune(mv.FilterCriteria.Year, '-') {
			if len(years) == 2 {
				value, err := strconv.Atoi(years[1])
				if err == nil && value <= int(maxYear) {
					maxYear = uint16(value)
				}
			}
		} else {
			maxYear = minYear
		}
		year = []uint16{minYear, maxYear}
	}

	final := make([]int, 0)
	for _, movie := range data.movies() {
		if (title == "" || strings.Contains(strings.ToLower(movie.Title), strings.ToLower(title))) &&
			(rate == nil || (rate[0] <= movie.Rate && movie.Rate <= rate[1])) &&
			(year == nil || (year[0] <= movie.Year && movie.Year <= year[1])) {
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
