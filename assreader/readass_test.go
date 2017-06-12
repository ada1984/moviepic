package assreader

import (
	"fmt"
	"os/exec"
	"testing"
)

func TestParseSub(t *testing.T) {
	subtitle := "Dialogue: Marked=0,2:01:37.10,0:00:45.30,*Default,NTP,0000,0000,0000,,{\\3c&HFF8000&}{\\fs18}多米尼克共和国"
	sub := NewSub(subtitle)
	fmt.Println(sub)
}

func TestNewAssReader(t *testing.T) {
	assReader := NewAssReader("D:/unsystem/movie/res/75913c586447de9f2140bad20bad287f.ass")
	fmt.Println(assReader.Subs[0])
}

func TestFuckme(t *testing.T) {
	fmt.Println(exec.Command("C:\\Go\\project\\src\\moviepic\\ffmpeg", "-i 1.mkv 1.ass").Output())
}
