package db

func init() {

}

//Movie ...
type Movie struct {
	ID       int
	Name     string
	FileName string
	Md5      string
}

//NewMovie ...
func NewMovie() *Movie {
	movie := &Movie{}
	return movie
}

func (movie *Movie) Save() {
	db.Where(movie).First(movie)
	if db.NewRecord(movie) {
		db.Create(movie)
	}
}
