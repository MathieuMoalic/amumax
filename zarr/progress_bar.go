package zarr

import (
	"fmt"
	"strings"
)

type ProgressBar struct {
	start float64
	stop  float64
	last  int
}

func (bar *ProgressBar) New(start float64, stop float64) {
	bar.start = start
	bar.stop = stop
	bar.last = 0
}

func (bar *ProgressBar) Update(current float64) {
	p := int((current-bar.start)/(bar.stop-bar.start)*100) + 1
	if p > 100 {
		p = 100
	}
	if p > bar.last {
		bar.last = p
		fmt.Printf("\r// [%-50s]%3d%%", strings.Repeat("█", p/2), p)
	}

}

func (bar *ProgressBar) Finish() {
	fmt.Printf("\r// [%-50s]%3d%%\n", strings.Repeat("█", 50), 100)
}
