package assreader

import (
	"strings"
	"time"
)

func timeToSeconds(v string) float64 {
	vSplit := strings.Split(v, ":")
	newTimeFormat := vSplit[0] + "h" + vSplit[1] + "m" + vSplit[2] + "s"
	duration, _ := time.ParseDuration(newTimeFormat)
	return duration.Seconds()
}

type SubText struct {
	Formats []string
	Text    string
}

func NewSubText(text string) *SubText {
	subText := &SubText{}
	subText.Parse(text)
	return subText
}

func (subText *SubText) Parse(text string) {
	splitedSubText := strings.SplitAfter(text, "}")
	for _, v := range splitedSubText {
		if len(v) < 1 {
			subText.Text = ""
		} else if v[len(v)-1] == '}' {
			subText.Formats = append(subText.Formats, v)
		} else {
			subText.Text = v
		}
	}
}
