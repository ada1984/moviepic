package db

func init() {

}

//Movie ...
type Movie struct {
	ID        int
	Name      string
	FileName  string
	Md5       string
	AssStream int
}

//NewMovie ...
func NewMovie() *Movie {
	movie := &Movie{}
	movie.Init()
	return movie
}

func (movie *Movie) Init() {
	db.Where(movie).First(movie)
}

func (movie *Movie) Save() {
	db.Where(movie).First(movie)
	if db.NewRecord(movie) {
		db.Create(movie)
	} else {
		db.Save(&movie)
	}
}
