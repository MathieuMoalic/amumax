package engine_old

import (
	"fmt"
	"math"

	"github.com/MathieuMoalic/amumax/src/engine_old/data_old"
	"github.com/MathieuMoalic/amumax/src/engine_old/log_old"
	"github.com/MathieuMoalic/amumax/src/engine_old/mag_old"
	"github.com/MathieuMoalic/amumax/src/engine_old/mesh_old"
)

// For a nanowire magnetized in-plane, with mx = mxLeft on the left end and
// mx = mxRight on the right end (both -1 or +1), add a B field needed to compensate
// for the surface charges on the left and right edges.
// This will mimic an infinitely long wire.
func removeLRSurfaceCharge(region int, mxLeft, mxRight float64) {
	setBusy(true)
	defer setBusy(false)
	log_old.AssertMsg(mxLeft == 1 || mxLeft == -1,
		"Invalid value for mxLeft: must be either +1 or -1 in removeLRSurfaceCharge")
	log_old.AssertMsg(mxRight == 1 || mxRight == -1,
		"Invalid value for mxRight: must be either +1 or -1 in removeLRSurfaceCharge")

	bsat := Msat.GetRegion(region) * mag_old.Mu0
	log_old.AssertMsg(bsat != 0,
		"RemoveSurfaceCharges: Msat is zero in region "+fmt.Sprint(region))

	B_ext.Add(compensateLRSurfaceCharges(GetMesh(), mxLeft, mxRight, bsat), nil)
}

func compensateLRSurfaceCharges(m *mesh_old.Mesh, mxLeft, mxRight float64, bsat float64) *data_old.Slice {
	h := data_old.NewSlice(3, m.Size())
	H := h.Vectors()
	world := m.WorldSize()
	cell := m.CellSize()
	size := m.Size()
	q := cell[Z] * cell[Y] * bsat
	q1 := q * mxLeft
	q2 := q * (-mxRight)

	// surface loop (source)
	for I := 0; I < size[Z]; I++ {
		for J := 0; J < size[Y]; J++ {

			y := (float64(J) + 0.5) * cell[Y]
			z := (float64(I) + 0.5) * cell[Z]
			source1 := [3]float64{0, y, z}        // left surface source
			source2 := [3]float64{world[X], y, z} // right surface source

			// volume loop (destination)
			for iz := range H[0] {
				for iy := range H[0][iz] {
					for ix := range H[0][iz][iy] {

						dst := [3]float64{ // destination coordinate
							(float64(ix) + 0.5) * cell[X],
							(float64(iy) + 0.5) * cell[Y],
							(float64(iz) + 0.5) * cell[Z]}

						h1 := hfield(q1, source1, dst)
						h2 := hfield(q2, source2, dst)

						// add this surface charges' field to grand total
						for c := 0; c < 3; c++ {
							H[c][iz][iy][ix] += float32(h1[c] + h2[c])
						}
					}
				}
			}
		}
	}
	return h
}

// H field of charge at location source, evaluated in location dest.
func hfield(charge float64, source, dest [3]float64) [3]float64 {
	var R [3]float64
	R[0] = dest[0] - source[0]
	R[1] = dest[1] - source[1]
	R[2] = dest[2] - source[2]
	r := math.Sqrt(R[0]*R[0] + R[1]*R[1] + R[2]*R[2])
	qr3pi4 := charge / ((4 * math.Pi) * r * r * r)
	var h [3]float64
	h[0] = R[0] * qr3pi4
	h[1] = R[1] * qr3pi4
	h[2] = R[2] * qr3pi4
	return h
}
