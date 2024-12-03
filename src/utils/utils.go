package utils

import (
	"encoding/binary"
	"fmt"
	"math"

	"github.com/MathieuMoalic/amumax/src/quantity"
)

func SizeOf(block [][][]float32) [3]int {
	return [3]int{len(block[0][0]), len(block[0]), len(block)}
}

func Reshape(array []float32, size [3]int) [][][]float32 {
	Nx, Ny, Nz := size[0], size[1], size[2]
	if Nx*Ny*Nz != len(array) {
		panic(fmt.Errorf("reshape: size mismatch: %v*%v*%v != %v", Nx, Ny, Nz, len(array)))
	}
	sliced := make([][][]float32, Nz)
	for i := range sliced {
		sliced[i] = make([][]float32, Ny)
	}
	for i := range sliced {
		for j := range sliced[i] {
			sliced[i][j] = array[(i*Ny+j)*Nx+0 : (i*Ny+j)*Nx+Nx]
		}
	}
	return sliced
}

// mumax specific formatting (Slice -> average, etc).
func CustomFmt(msg []interface{}) (fmtMsg string) {
	for _, m := range msg {
		if e, ok := m.(quantity.Quantity); ok {
			str := fmt.Sprint(e.Average())
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

func Prod(s [3]int) int {
	return s[0] * s[1] * s[2]
}

// average of slice over universe
// func AverageSlice(s *slice.Slice) []float64 {
// 	nCell := float64(Prod(s.Size()))
// 	avg := make([]float64, s.NComp())
// 	for i := range avg {
// 		avg[i] = float64(cuda.Sum(s.Comp(i))) / nCell
// 		if math.IsNaN(avg[i]) {
// 			panic("NaN")
// 		}
// 	}
// 	return avg
// }

func Sign(x float64) float64 {
	switch {
	case x > 0:
		return 1
	case x < 0:
		return -1
	default:
		return 0
	}
}

func Fmod(a, b float64) float64 {
	if b == 0 || math.IsInf(b, 1) {
		return a
	}
	if math.Abs(a) > b/2 {
		return Sign(a) * (math.Mod(math.Abs(a+b/2), b) - b/2)
	} else {
		return a
	}
}

func Sqr64(x float64) float64 { return x * x }

func Float64ToBytes(f float64) []byte {
	var buf [8]byte
	binary.LittleEndian.PutUint64(buf[:], math.Float64bits(f))
	return buf[:]
}

func Float32ToBytes(f float32) []byte {
	var buf [4]byte
	binary.LittleEndian.PutUint32(buf[:], math.Float32bits(f))
	return buf[:]
}

func BytesToFloat32(b []byte) float32 {
	return math.Float32frombits(binary.LittleEndian.Uint32(b))
}

func Unslice(v []float64) [3]float64 {
	if len(v) != 3 {
		panic(fmt.Errorf("length mismatch: input slice must have exactly 3 elements in unslice"))
	}
	return [3]float64{v[0], v[1], v[2]}
}

func Slice(v [3]float64) []float64 {
	return v[:]
}
