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

const (
	SQLFindAllPicNamesByMovie string = "select distinct pic.name from moviepic.pic,moviepic.subtitle where subtitle.movie_id = ? and pic.subtitle_id = subtitle.id ;"
)

//NewMovie ...
func NewMovie() *Movie {
	movie := &Movie{}
	movie.Init()
	return movie
}

func (movie *Movie) Init() {
	if movie.Md5 != "" {
		db.Where("md5 = ?", movie.Md5).First(movie)
	}
}

func (movie *Movie) Save() {
	db.Where(movie).First(movie)
	if db.NewRecord(movie) {
		db.Create(movie)
	} else {
		db.Save(&movie)
	}
}

func (movie *Movie) FindAll() []Movie {
	movies := []Movie{}
	db.Where(movie).Find(&movies)
	return movies
}

func (movie *Movie) FindAllPicNamesByMovie() []int {
	picNames := []int{}
	rows, _ := db.Raw(SQLFindAllPicNamesByMovie, movie.ID).Rows()
	defer rows.Close()
	for rows.Next() {
		picName := 0
		// fmt.Println(rows.Scan(&pic))
		rows.Scan(&picName)
		picNames = append(picNames, picName)
	}
	return picNames
}
