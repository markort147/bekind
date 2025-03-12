package movies

type Movie struct {
	Id    int
	Title string
	Year  uint16
	Rate  uint8
}
type MovieField int

const (
	MovieId MovieField = iota
	MovieTitle
	MovieYear
	MovieRate
)
