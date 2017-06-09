package assreader

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

//...
const (
	DIALOGUE = "Dialogue"
)

//Sub ...
type Sub struct {
	Marked  int
	Start   float64
	End     float64
	Style   string
	Name    string
	MarginL string
	MarginR string
	MarginV string
	Effect  string
	Text    SubText
}

//NewSub ....
func NewSub(subLine string) *Sub {
	sub := &Sub{}
	sub.ParseSub(subLine)
	return sub
}

//AssReader ...
type AssReader struct {
	Path string
	Subs []Sub
}

//NewAssReader ..
func NewAssReader(path string) *AssReader {
	assReader := &AssReader{Path: path}
	file, err := os.Open(path)
	if err != nil {
		return nil
	}

	isValidAss := false
	fscanner := bufio.NewScanner(file)
	for fscanner.Scan() {
		textLine := fscanner.Text()
		if len(textLine) > len(DIALOGUE) && textLine[:len(DIALOGUE)] == DIALOGUE {
			assReader.Subs = append(assReader.Subs, *NewSub(textLine))
			isValidAss = true
		}
	}
	if !isValidAss {
		return nil
	}
	return assReader
}

//ParseSub ...
func (sub *Sub) ParseSub(subLine string) {
	lineSplit := strings.SplitN(subLine, ",", 10)
	// if len(lineSplit) != 10 {
	// 	return
	// }
	sub.Marked, _ = strconv.Atoi(lineSplit[0])
	sub.Start = timeToSeconds(lineSplit[1])
	sub.End = timeToSeconds(lineSplit[2])
	sub.Style = lineSplit[3]
	sub.Name = lineSplit[4]
	sub.MarginL = lineSplit[5]
	sub.MarginR = lineSplit[6]
	sub.MarginV = lineSplit[7]
	sub.Effect = lineSplit[8]
	sub.Text = *NewSubText(lineSplit[9])
}
