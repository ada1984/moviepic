package utils

type Movie struct {
	FilePath string
	AssPath  string
	Ext      string
}

func NewMovie(path string, ext string) *Movie {
	movie := &Movie{FilePath: path, Ext: ext}
	return movie
}
