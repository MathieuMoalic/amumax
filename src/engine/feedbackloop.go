package engine

import (
	"sync"

	"github.com/MathieuMoalic/amumax/src/data"
	"github.com/MathieuMoalic/amumax/src/log"
)

func init() {
	DeclFunc("FeedbackLoop", FeedbackLoop, "FeedbackLoop(input_mask, output_mask *data.Slice, multiplier float64)")
}

// FeedbackLoop reads the “current” from the output antenna (as the dot of m with
// output_mask), multiplies it by `multiplier`, and reinjects it through the input
// antenna as a dynamic contribution to B_ext.
// - input_mask: field mask of the drive antenna
// - output_mask: field mask used as pickup
// - multiplier: scalar gain applied to the measured signal before reinjection
func FeedbackLoop(inputMask, outputMask *data.Slice, multiplier float64) {
	if inputMask.NComp() != 3 || outputMask.NComp() != 3 {
		log.Log.Err("%s", "FeedbackLoop expects vector masks (3 components)")
	}
	// Auto-zero: capture the initial pickup as baseline and subtract it forever after.
	var once sync.Once
	var baseline float64

	BExt.AddGo(inputMask, func() float64 {
		// Get the current value of m projected on outputMask
		// (this is the "measured" signal from the output antenna)
		s := magModulatedByMask(outputMask)

		// On the first call, set the baseline to the current value.
		// Thereafter, always subtract the baseline.
		once.Do(func() { baseline = s }) // set on first call

		// Reinject the (baseline-subtracted) signal through the input antenna
		// scaled by the multiplier.
		return multiplier * (s - baseline)
	})
}
