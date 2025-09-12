package engine

import (
	"sync"

	"github.com/MathieuMoalic/amumax/src/data"
	"github.com/MathieuMoalic/amumax/src/log"
)

// register in the scripting API (optional, but handy)
func init() {
	DeclFunc("FeedbackLoop", FeedbackLoop, "FeedbackLoop(input_mask, output_mask *data.Slice, multiplier float64)")
}

// FeedbackLoop reads the “current” from the output antenna (as the dot of m with
// output_mask), multiplies it by `multiplier`, and reinjects it through the input
// antenna as a dynamic contribution to B_ext.
// - input_mask: field mask of the drive antenna. Accepted sizes: (Nx,Ny,Nz) or (Nx,1,Nz)
// - output_mask: field mask used as pickup (usually (Nx,1,Nz))
// - multiplier: scalar gain applied to the measured signal before reinjection
func FeedbackLoop(input_mask, output_mask *data.Slice, multiplier float64) {
	Nx, Ny, Nz := Mesh.GetNi()

	if input_mask.NComp() != 3 || output_mask.NComp() != 3 {
		log.Log.Err("%s", "FeedbackLoop expects vector masks (3 components)")
	}

	// Ensure the drive slice matches the mesh. If input_mask has Ny=1 (typical mask),
	// broadcast it along y so B_ext.AddGo can use it directly.
	inT := input_mask.Tensors() // [c][z][y][x]
	xDim := len(inT[0][0][0])
	yDim := len(inT[0][0])
	zDim := len(inT[0])

	var drive *data.Slice
	switch {
	case xDim == Nx && yDim == Ny && zDim == Nz:
		drive = input_mask
	case xDim == Nx && yDim == 1 && zDim == Nz:
		drive = data.NewSlice(3, [3]int{Nx, Ny, Nz})
		for iz := 0; iz < Nz; iz++ {
			for iy := 0; iy < Ny; iy++ {
				for ix := 0; ix < Nx; ix++ {
					drive.Set(0, ix, iy, iz, float64(inT[0][iz][0][ix]))
					drive.Set(1, ix, iy, iz, float64(inT[1][iz][0][ix]))
					drive.Set(2, ix, iy, iz, float64(inT[2][iz][0][ix]))
				}
			}
		}
	default:
		log.Log.Err("%s", "FeedbackLoop: input_mask must be sized (Nx,Ny,Nz) or (Nx,1,Nz)")
	}

	// Auto-zero: capture the initial pickup as baseline and subtract it forever after.
	var once sync.Once
	var baseline float64

	B_ext.AddGo(drive, func() float64 {
		s := magModulatedByMask(output_mask)
		once.Do(func() { baseline = s }) // set on first call
		return multiplier * (s - baseline)
	})
}
