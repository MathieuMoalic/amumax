package zarr

import (
	"os"
	"strings"

	"github.com/fatih/color"
	"golang.org/x/term"
)

type ProgressBar struct {
	start float64
	stop  float64
	last  int
	out   *color.Color
}

func (bar *ProgressBar) New(start float64, stop float64) {
	bar.start = start
	bar.stop = stop
	bar.last = 0
	bar.out = color.New(color.FgGreen)
}

func (bar *ProgressBar) GetTermWidth() int {
	width, _, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		width = 30
	}
	return width
}

func (bar *ProgressBar) Update(current_time float64) {
	percentage := int((current_time-bar.start)/(bar.stop-bar.start)*100) + 1
	if percentage > 100 {
		percentage = 100
	}
	if percentage > bar.last {
		width := bar.GetTermWidth()
		bar.last = percentage
		a := int((width - 4) / 2 * percentage / 100)
		b := int((width - 4) * (100 - percentage) / 100)
		bar.out.Print("\r//[" + strings.Repeat("🧲", a) + strings.Repeat(" ", b) + "]")
	}

}

func (bar *ProgressBar) Finish() {
	bar.out.Println("\r//[" + strings.Repeat("🧲", int(bar.GetTermWidth()-4)/2) + "]")
}
