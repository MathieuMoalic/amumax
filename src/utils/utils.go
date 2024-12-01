package utils

import (
	"encoding/binary"
	"fmt"
	"math"

	"github.com/MathieuMoalic/amumax/src/cuda"
	"github.com/MathieuMoalic/amumax/src/data"
	"github.com/MathieuMoalic/amumax/src/quantity"
)

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
func AverageSlice(s *data.Slice) []float64 {
	nCell := float64(Prod(s.Size()))
	avg := make([]float64, s.NComp())
	for i := range avg {
		avg[i] = float64(cuda.Sum(s.Comp(i))) / nCell
		if math.IsNaN(avg[i]) {
			panic("NaN")
		}
	}
	return avg
}
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

// func bytesToFloat64(bytes []byte) float64 {
// 	bits := binary.LittleEndian.Uint64(bytes)
// 	float := math.Float64frombits(bits)
// 	return float
// }

// func bytesToFloat32(bytes []byte) float32 {
// 	bits := binary.LittleEndian.Uint32(bytes)
// 	float := math.Float32frombits(bits)
// 	return float
// }

// func float32ToBytes(f float32) []byte {
// 	var buf [4]byte
// 	binary.LittleEndian.PutUint32(buf[:], math.Float32bits(f))
// 	return buf[:]
// }
