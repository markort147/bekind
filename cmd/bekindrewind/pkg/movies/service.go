package movies

type MovieServiceErr string

const (
	MovieNotFoundErr = MovieServiceErr("Movie not found")
)

func (e MovieServiceErr) Error() string {
	return string(e)
}

var data Data
var view View

func Init() {
	initData()
	initView()
}

func GetView() View {
	return view
}

func initData() {
	data = Data{}
	data.Movies = make([]*Movie, 0)
	data.MoviesMap = make(map[int]*Movie)
}

func initView() {
	view.data = &data
	for id, _ := range data.MoviesMap {
		view.MovieIds = append(view.MovieIds, id)
	}
	view.SortInfo = SortCriteria{
		SortedBy: MovieId,
		Desc:     false,
	}
	view.refresh()
}

func FindById(id int) (Movie, error) {
	movie, exists := data.MoviesMap[id]
	if !exists {
		return Movie{}, MovieNotFoundErr
	}
	return *movie, nil
}

func Save(m Movie) Movie {
	saved := data.addMovie(m)
	view.refresh()
	return saved
}

func Update(id int, new Movie) {
	old := data.MoviesMap[id]
	old.Rate = new.Rate
	old.Title = new.Title
	old.Year = new.Year
	view.refresh()
}

func Delete(id int) bool {
	remove := -1
	for i, m := range data.Movies {
		if m.Id == id {
			remove = i
			break
		}
	}

	if remove == -1 {
		return false
	}

	data.Movies = append(data.Movies[:remove], data.Movies[remove+1:]...)
	delete(data.MoviesMap, id)

	view.refresh()
	return true
}

func SortView(by string) {
	field := strToMF(by)
	if view.SortInfo.SortedBy == field {
		view.SortInfo.Desc = !view.SortInfo.Desc
	}
	view.SortInfo.SortedBy = field
	view.refreshSorting()
}

func FilterView(criteria FilterCriteria) {
	view.FilterCriteria = criteria
	view.refreshFilter()
}
