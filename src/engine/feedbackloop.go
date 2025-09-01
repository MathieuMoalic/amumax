package engine

import (
	"github.com/MathieuMoalic/amumax/src/data"
)

func init() {
	DeclFunc("FeedbackLoop", FeedbackLoop,
		"Add m(src)*multiplier [T] into B_eff at dst each iteration.")
}

// FeedbackLoop adds m(src)*multiplier [T] into B_eff at dst each iteration.
func FeedbackLoop(input_mask, output_mask *data.Slice, multiplier float64) {
}
