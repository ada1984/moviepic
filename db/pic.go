package db

//Pic ...
type Pic struct {
	ID         int
	Name       int
	Time       int
	SubtitleID int
}

//Save ...
func (pic *Pic) Save() {
	db.Create(pic)
}

//Exist ...
func (pic *Pic) Exist() bool {
	newPic := &Pic{Time: pic.Time, SubtitleID: pic.SubtitleID}
	db.Where(newPic).First(newPic)
	if db.NewRecord(newPic) { //cannot find
		return false
	}
	return true
}
