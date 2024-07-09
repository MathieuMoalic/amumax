package api

import (
	"encoding/binary"

	"math"

	"github.com/MathieuMoalic/amumax/cuda"
	"github.com/MathieuMoalic/amumax/data"
	"github.com/MathieuMoalic/amumax/engine"
	"github.com/MathieuMoalic/amumax/util"
)

var bps BackendPreviewState

func init() {
	bps = BackendPreviewState{
		Quantity:   &engine.M,
		Component:  -1,
		Layer:      0,
		MaxPoints:  10000,
		Dimensions: [3]int{0, 0, 0},
		Type:       "vector",
		Buffer:     []byte{},
	}
}

type BackendPreviewState struct {
	Quantity   engine.Quantity `msgpack:"quantity"`
	Component  int             `msgpack:"component"`
	Layer      int             `msgpack:"layer"`
	MaxPoints  int             `msgpack:"maxPoints"`
	Dimensions [3]int          `msgpack:"dimensions"`
	Type       string          `msgpack:"type"`
	Buffer     []byte          `msgpack:"buffer"`
}

func compStringToIndex(comp string) int {
	switch comp {
	case "x":
		return 0
	case "y":
		return 1
	case "z":
		return 2
	case "All":
		return -1
	}
	util.Fatal("Invalid component string")
	return -2
}

func compIndexToString(comp int) string {
	switch comp {
	case 0:
		return "x"
	case 1:
		return "y"
	case 2:
		return "z"
	case -1:
		return "All"
	}
	util.Fatal("Invalid component index")
	return ""
}

var PreviewBuffer []byte

// scaleOutputBuffer scales down the image size until the number of points are < MaxVectors
func scaleOutputBuffer(originalSize [3]int) [3]int {
	width, height := originalSize[0], originalSize[1]
	points := width * height

	// Calculate the scaling factor
	for points >= bps.MaxPoints {
		width = (width + 1) / 2
		height = (height + 1) / 2
		points = width * height
	}
	return [3]int{width, height, 1}
}

// ConvertVectorFieldToBinary converts a 4D vector field to a binary byte slice.
func ConvertVectorFieldToBinary(vectorField [3][][][]float32) []byte {
	// Calculate the total number of elements
	xLen := len(vectorField[0])
	yLen := len(vectorField[0][0])
	zLen := len(vectorField[0][0][0])

	// Allocate buffer for binary data
	// Each vector field element has 3 float32 values, each float32 is 4 bytes.
	buf := make([]byte, xLen*yLen*zLen*3*4)

	index := 0
	for i := 0; i < xLen; i++ {
		for j := 0; j < yLen; j++ {
			for k := 0; k < zLen; k++ {
				// Put X, Y, Z values into the buffer
				binary.LittleEndian.PutUint32(buf[index:], math.Float32bits(vectorField[0][i][j][k]))
				index += 4
				binary.LittleEndian.PutUint32(buf[index:], math.Float32bits(vectorField[1][i][j][k]))
				index += 4
				binary.LittleEndian.PutUint32(buf[index:], math.Float32bits(vectorField[2][i][j][k]))
				index += 4
			}
		}
	}

	return buf
}

func ConvertScalarFieldToBinary(scalarField [][][]float32) []byte {
	xLen := len(scalarField)
	yLen := len(scalarField[0])
	zLen := len(scalarField[0][0])

	buf := make([]byte, xLen*yLen*zLen*4)

	index := 0
	for i := 0; i < xLen; i++ {
		for j := 0; j < yLen; j++ {
			for k := 0; k < zLen; k++ {
				binary.LittleEndian.PutUint32(buf[index:], math.Float32bits(scalarField[i][j][k]))
				index += 4
			}
		}
	}

	return buf
}

func PreparePreviewBuffer() {
	nComp := 1
	if bps.Type == "vector" {
		nComp = 3
	}
	originalSize := engine.MeshOf(bps.Quantity).Size()
	outputSize := scaleOutputBuffer(originalSize)
	bps.Dimensions = outputSize
	scaledSlice := data.NewSlice(nComp, outputSize)
	tempBuffer := cuda.NewSlice(1, outputSize)

	GPUBuffer := engine.ValueOf(bps.Quantity)
	defer cuda.Recycle(GPUBuffer)
	defer tempBuffer.Free()

	if bps.Type == "vector" {
		for c := 0; c < nComp; c++ {
			cuda.Resize(tempBuffer, GPUBuffer.Comp(c), bps.Layer)
			data.Copy(scaledSlice.Comp(c), tempBuffer)
		}
	} else {
		if bps.Quantity.NComp() == 3 {
			cuda.Resize(tempBuffer, GPUBuffer.Comp(bps.Component), bps.Layer)
			data.Copy(scaledSlice.Comp(0), tempBuffer)
		} else {
			cuda.Resize(tempBuffer, GPUBuffer.Comp(0), bps.Layer)
			data.Copy(scaledSlice.Comp(0), tempBuffer)
		}
	}

	if bps.Type == "vector" {
		normalizeVectors(scaledSlice)
		bps.Buffer = ConvertVectorFieldToBinary(scaledSlice.Vectors())
	} else {
		normalizeScalars(scaledSlice)
		bps.Buffer = ConvertScalarFieldToBinary(scaledSlice.Scalars())
	}
}

func normalizeVectors(f *data.Slice) {
	a := f.Vectors()
	maxnorm := 0.
	for i := range a[0] {
		for j := range a[0][i] {
			for k := range a[0][i][j] {

				x, y, z := a[0][i][j][k], a[1][i][j][k], a[2][i][j][k]
				norm := math.Sqrt(float64(x*x + y*y + z*z))
				if norm > maxnorm {
					maxnorm = norm
				}

			}
		}
	}
	factor := float32(1 / maxnorm)

	for i := range a[0] {
		for j := range a[0][i] {
			for k := range a[0][i][j] {
				a[0][i][j][k] *= factor
				a[1][i][j][k] *= factor
				a[2][i][j][k] *= factor

			}
		}
	}
}

func normalizeScalars(f *data.Slice) {
	a := f.Scalars()
	maxnorm := 0.
	for i := range a {
		for j := range a[i] {
			for k := range a[i][j] {

				norm := math.Abs(float64(a[i][j][k]))
				if norm > maxnorm {
					maxnorm = norm
				}

			}
		}
	}
	factor := float32(1 / maxnorm)

	for i := range a {
		for j := range a[i] {
			for k := range a[i][j] {
				a[i][j][k] *= factor

			}
		}
	}
}
