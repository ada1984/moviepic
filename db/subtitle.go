package db

//Subtitle ...
type Subtitle struct {
	ID      int
	Text    string
	Lang    int
	Format  string
	Start   int
	End     int
	MovieID int
	Pic     Pic
	Movie   Movie
}

//Insert ...
func (subtitle *Subtitle) Save() {
	db.Where(subtitle).First(subtitle)
	if db.NewRecord(subtitle) {
		db.Create(subtitle)
	}
}
