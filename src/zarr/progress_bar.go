package zarr

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"golang.org/x/term"
)

type ProgressBar struct {
	start     float64
	stop      float64
	last      int
	out       *color.Color
	enabled   bool
	symbol    string
	minLength int
}

func NewProgressBar(start float64, stop float64, symbol string, enabled bool) *ProgressBar {
	return &ProgressBar{
		start:     start,
		stop:      stop,
		last:      -1,
		out:       color.New(color.FgGreen),
		enabled:   enabled,
		symbol:    symbol,
		minLength: 10,
	}
}

func (bar *ProgressBar) getTermWidth() int {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || width <= 0 {
		width = bar.minLength // Default width if unable to get terminal size
	}
	return width
}

func (bar *ProgressBar) Update(currentTime float64) {
	if !bar.enabled {
		return
	}

	percentage := int((currentTime - bar.start) / (bar.stop - bar.start) * 100)
	if percentage > 100 {
		percentage = 100
	}
	if percentage < 0 {
		percentage = 0
	}

	if percentage > bar.last {
		width := bar.getTermWidth()
		bar.last = percentage

		// Adjust the bar length to account for fixed characters
		barLength := width - bar.minLength // For "[", "]", percentage, and spaces
		if barLength < bar.minLength {
			barLength = bar.minLength // Minimum bar length
		}

		filledLength := int(float64(barLength) * float64(percentage) / 100.0)
		emptyLength := barLength - filledLength

		progressBar := fmt.Sprintf("\r[%s%s] %3d%%",
			strings.Repeat(bar.symbol, filledLength),
			strings.Repeat(" ", emptyLength),
			percentage)

		// Clear the line
		fmt.Print("\r" + strings.Repeat(" ", width))
		// Write the progress bar
		bar.out.Print(progressBar)
	}
}

func (bar *ProgressBar) Finish() {
	if bar.enabled {
		width := bar.getTermWidth()

		barLength := width - bar.minLength
		if barLength < bar.minLength {
			barLength = bar.minLength
		}

		progressBar := fmt.Sprintf("\r[%s] 100%%\n",
			strings.Repeat(bar.symbol, barLength))

		// Clear the line
		fmt.Print("\r" + strings.Repeat(" ", width))
		// Write the progress bar
		bar.out.Print(progressBar)
	}
}
