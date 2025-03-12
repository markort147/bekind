package movies

import (
	"slices"
	"strconv"
	"strings"
)

type ViewHeader struct {
	Id    string
	Title string
	Year  string
	Rate  string
}

type FilterCriteria struct {
	Title string
	Rate  string
	Year  string
}

type SortCriteria struct {
	SortedBy MovieField
	Desc     bool
}

type View struct {
	SortInfo       SortCriteria
	Header         ViewHeader
	MovieIds       []int
	data           *Data
	FilterCriteria FilterCriteria
}

var CurrView View

func (v *View) initView(data *Data) {
	CurrView.data = data
	for id, _ := range data.MoviesMap {
		CurrView.MovieIds = append(CurrView.MovieIds, id)
	}
	field, _ := MFToStr(MovieId)
	CurrView.Sort(field)
}

func (v *View) refresh() {
	v.refreshFilter()
	v.refreshSorting()
}

func (v *View) Sort(by string) {
	field := StrToMF(by)
	if v.SortInfo.SortedBy == field {
		v.SortInfo.Desc = !v.SortInfo.Desc
	}
	v.SortInfo.SortedBy = field
	v.refreshSorting()
}

func (v *View) refreshSorting() {
	slices.SortFunc(v.MovieIds, func(i, j int) int {
		var res int
		switch v.SortInfo.SortedBy {
		case MovieId:
			res = i - j
		case MovieRate:
			res = int(v.data.MoviesMap[i].Rate) - int(v.data.MoviesMap[j].Rate)
		case MovieYear:
			res = int(v.data.MoviesMap[i].Year) - int(v.data.MoviesMap[j].Year)
		case MovieTitle:
			res = strings.Compare(strings.ToLower(v.data.MoviesMap[i].Title), strings.ToLower(v.data.MoviesMap[j].Title))
		}
		if v.SortInfo.Desc {
			res = -res
		}
		return res
	})

	v.Header.Id = sortLabel(MovieId, v.SortInfo)
	v.Header.Title = sortLabel(MovieTitle, v.SortInfo)
	v.Header.Year = sortLabel(MovieYear, v.SortInfo)
	v.Header.Rate = sortLabel(MovieRate, v.SortInfo)
}

func sortLabel(field MovieField, info SortCriteria) string {
	res, _ := MFToStr(field)
	if field == info.SortedBy {
		if info.Desc {
			res = res + " ↑"
		} else {
			res = res + " ↓"
		}
	}
	return res
}

func (v *View) Filter(criteria FilterCriteria) {
	v.FilterCriteria = criteria
	v.refreshFilter()
}

func (v *View) refreshFilter() { // title
	title := v.FilterCriteria.Title

	// rate range
	var rate []uint8
	if v.FilterCriteria.Rate != "" {
		minRate := uint8(0)
		maxRate := uint8(10)
		rates := strings.Split(v.FilterCriteria.Rate, "-")
		if len(rates) > 0 {
			value, err := strconv.Atoi(rates[0])
			if err == nil && value >= 0 {
				minRate = uint8(value)
			}
		}
		if strings.ContainsRune(v.FilterCriteria.Rate, '-') {
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
	if v.FilterCriteria.Year != "" {
		minYear := uint16(0)
		maxYear := ^uint16(0)
		years := strings.Split(v.FilterCriteria.Year, "-")
		if len(years) > 0 {
			value, err := strconv.Atoi(years[0])
			if err == nil && value >= int(minYear) {
				minYear = uint16(value)
			}
		}
		if strings.ContainsRune(v.FilterCriteria.Year, '-') {
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
	for i, movie := range v.data.MoviesMap {
		if (title == "" || strings.Contains(strings.ToLower(movie.Title), strings.ToLower(title))) &&
			(rate == nil || (rate[0] <= movie.Rate && movie.Rate <= rate[1])) &&
			(year == nil || (year[0] <= movie.Year && movie.Year <= year[1])) {
			final = append(final, i)
		}
	}
	v.MovieIds = final
}
