package zarr

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"golang.org/x/term"
)

type ProgressBar struct {
	start       float64
	stop        float64
	last        int
	out         *color.Color
	disabled    bool
	symbol      string
	symbolWidth int
	minLength   int
	lastUpdate  time.Time
}

func NewProgressBar(start float64, stop float64, symbol string, hideProgressBar bool) *ProgressBar {
	return &ProgressBar{
		start:       start,
		stop:        stop,
		last:        -1,
		out:         color.New(color.FgGreen),
		disabled:    hideProgressBar,
		symbol:      symbol,
		symbolWidth: 2, // Width of the symbol in characters
		minLength:   1, // Minimum bar length in symbols
		lastUpdate:  time.Now(),
	}
}

func (bar *ProgressBar) getTermWidth() int {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || width <= 0 {
		width = 80 // Default width if unable to get terminal size
	}
	return width
}

func (bar *ProgressBar) calculateDimensions(percentage int) (int, int) {
	width := bar.getTermWidth()
	fixedLength := 1 + 1 + 1 + 4 // '[', ']', space, '100%'
	availableWidth := width - fixedLength
	availableWidth -= availableWidth % bar.symbolWidth
	barLength := availableWidth / bar.symbolWidth
	if barLength < bar.minLength {
		barLength = bar.minLength
	}
	filledLength := int(float64(barLength) * float64(percentage) / 100.0)
	if filledLength > barLength {
		filledLength = barLength
	}
	emptyLength := barLength - filledLength
	return emptyLength, filledLength
}

func (bar *ProgressBar) Update(currentTime float64) {
	if bar.disabled {
		return
	}
	now := time.Now()
	if now.Sub(bar.lastUpdate) < time.Millisecond*100 {
		return
	}
	bar.lastUpdate = now

	percentage := int((currentTime - bar.start) / (bar.stop - bar.start) * 100)
	if percentage > 100 {
		percentage = 100
	}
	if percentage < 0 {
		percentage = 0
	}

	if percentage > bar.last {
		width := bar.getTermWidth()
		if width < 10 {
			// Don't display the progress bar if the terminal is too small
			return
		}
		bar.last = percentage

		emptyLength, filledLength := bar.calculateDimensions(percentage)

		filledSymbols := strings.Repeat(bar.symbol, filledLength)
		emptySymbols := strings.Repeat("  ", emptyLength)

		progressBar := fmt.Sprintf("\r[%s%s] %3d%%",
			filledSymbols,
			emptySymbols,
			percentage)

		// Clear the line
		bar.out.Print("\r\033[K")

		// Write the progress bar
		bar.out.Print(progressBar)
	}
}

func (bar *ProgressBar) Finish() {
	if !bar.disabled {
		_, filledLength := bar.calculateDimensions(100)

		filledSymbols := strings.Repeat(bar.symbol, filledLength)

		progressBar := fmt.Sprintf("\r[%s] 100%%\n",
			filledSymbols)

		// Clear the line
		bar.out.Print("\r\033[K")

		// Write the progress bar
		bar.out.Print(progressBar)
	}
}
