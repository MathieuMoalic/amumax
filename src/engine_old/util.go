package engine_old

import (
	"fmt"
	"math"
	"os"
	"path"
	"strings"

	"github.com/MathieuMoalic/amumax/src/cuda"
	"github.com/MathieuMoalic/amumax/src/data"
	"github.com/MathieuMoalic/amumax/src/fsutil_old"
	"github.com/MathieuMoalic/amumax/src/log_old"
	"github.com/MathieuMoalic/amumax/src/oommf"
	"github.com/MathieuMoalic/amumax/src/zarr"
)

func regionFromCoordinate(x, y, z int) int {
	return Regions.GetCell(x, y, z)
}

// Returns a new new slice (3D array) with given number of components and size.
func newSlice(ncomp, Nx, Ny, Nz int) *data.Slice {
	return data.NewSlice(ncomp, [3]int{Nx, Ny, Nz})
}

func newVectorMask(Nx, Ny, Nz int) *data.Slice {
	return data.NewSlice(3, [3]int{Nx, Ny, Nz})
}

func newScalarMask(Nx, Ny, Nz int) *data.Slice {
	return data.NewSlice(1, [3]int{Nx, Ny, Nz})
}

// Constructs a vector
func vector(x, y, z float64) data.Vector {
	return data.Vector{x, y, z}
}

// Test if have lies within want +/- maxError,
// and print suited message.
func expect(msg string, have, want, maxError float64) {
	if math.IsNaN(have) || math.IsNaN(want) || math.Abs(have-want) > maxError {
		log_old.Log.Info(msg, ":", " have: ", have, " want: ", want, "±", maxError)
		CleanExit()
		os.Exit(1)
	} else {
		log_old.Log.Info(msg, ":", have, "OK")
	}
	// note: we also check "want" for NaN in case "have" and "want" are switched.
}

func expectV(msg string, have, want data.Vector, maxErr float64) {
	for c := 0; c < 3; c++ {
		expect(fmt.Sprintf("%v[%v]", msg, c), have[c], want[c], maxErr)
	}
}

// Append msg to file. Used to write aggregated output of many simulations in one file.
func fprintln(filename string, msg ...interface{}) {
	if !path.IsAbs(filename) {
		filename = OD() + filename
	}
	err := fsutil_old.Touch(filename)
	log_old.Log.PanicIfError(err)
	err = fsutil_old.Append(filename, []byte(fmt.Sprintln(customFmt(msg))))
	log_old.Log.PanicIfError(err)
}

func loadFile(fname string) *data.Slice {
	var s *data.Slice
	s, err := zarr.Read(fname, OD())
	log_old.Log.PanicIfError(err)
	return s
}

func loadOvfFile(fname string) *data.Slice {
	in, err := fsutil_old.Open(fname)
	log_old.Log.PanicIfError(err)
	var s *data.Slice
	s, _, err = oommf.Read(in)
	log_old.Log.PanicIfError(err)
	return s
}

// download a quantity to host,
// or just return its data when already on host.
func download(q Quantity) *data.Slice {
	// TODO: optimize for Buffer()
	buf := ValueOf(q)
	defer cuda.Recycle(buf)
	if buf.CPUAccess() {
		return buf
	} else {
		return buf.HostCopy()
	}
}

// print with special formatting for some known types
func myprint(msg ...interface{}) {
	log_old.Log.Info("%v", customFmt(msg))
}

// mumax specific formatting (Slice -> average, etc).
func customFmt(msg []interface{}) (fmtMsg string) {
	for _, m := range msg {
		if e, ok := m.(Quantity); ok {
			str := fmt.Sprint(AverageOf(e))
			str = str[1 : len(str)-1] // remove [ ]
			fmtMsg += fmt.Sprintf("%v, ", str)
		} else {
			fmtMsg += fmt.Sprintf("%v, ", m)
		}
	}
	// remove trailing ", "
	if len(fmtMsg) > 2 {
		fmtMsg = fmtMsg[:len(fmtMsg)-2]
	}
	return
}

// converts cell index to coordinate, internal coordinates
func index2Coord(ix, iy, iz int) data.Vector {
	m := GetMesh()
	n := m.Size()
	c := m.CellSize()
	x := c[X]*(float64(ix)-0.5*float64(n[X]-1)) - totalShift
	y := c[Y]*(float64(iy)-0.5*float64(n[Y]-1)) - totalYShift
	z := c[Z] * (float64(iz) - 0.5*float64(n[Z]-1))
	return data.Vector{x, y, z}
}

func sign(x float64) float64 {
	switch {
	case x > 0:
		return 1
	case x < 0:
		return -1
	default:
		return 0
	}
}

// shortcut for slicing unaddressable_vector()[:]
func slice(v [3]float64) []float64 {
	return v[:]
}

func unslice(v []float64) [3]float64 {
	log_old.AssertMsg(len(v) == 3, "Length mismatch: input slice must have exactly 3 elements in unslice")
	return [3]float64{v[0], v[1], v[2]}
}

func assureGPU(s *data.Slice) *data.Slice {
	if s.GPUAccess() {
		return s
	} else {
		return cuda.GPUCopy(s)
	}
}

func checkNaN1(x float64) {
	if math.IsNaN(x) {
		panic("NaN")
	}
}

// trim trailing newlines
func rmln(a string) string {
	for strings.HasSuffix(a, "\n") {
		a = a[:len(a)-1]
	}
	return a
}

const (
	X = 0
	Y = 1
	Z = 2
)

const (
	SCALAR = 1
	VECTOR = 3
)
