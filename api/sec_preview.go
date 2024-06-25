package api

import (
	"net/http"

	"encoding/binary"
	"math"

	"github.com/MathieuMoalic/amumax/cuda"
	"github.com/MathieuMoalic/amumax/data"
	"github.com/MathieuMoalic/amumax/engine"
	"github.com/MathieuMoalic/amumax/util"
	"github.com/labstack/echo/v4"
)

var preview Preview

func init() {
	preview = Preview{
		Quantity:  &engine.M,
		Component: 0,
		Layer:     0,
		MaxPoints: 10000,
	}
}

type Preview struct {
	Quantity  engine.Quantity `msgpack:"quantity"`
	Component int             `msgpack:"component"`
	Layer     int             `msgpack:"layer"`
	MaxPoints int             `msgpack:"maxPoints"`
}

func postPreviewState(c echo.Context) error {
	type Request struct {
		Quantity  string `msgpack:"quantity"`
		Component string `msgpack:"component"`
		Layer     int    `msgpack:"layer"`
		MaxPoints int    `msgpack:"maxPoints"`
	}
	req := new(Request)
	if err := c.Bind(req); err != nil {
		util.LogErr(err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request payload"})
	}

	quantity, exists := engine.Quantities[req.Quantity]
	if !exists {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Quantity not found"})
	}
	if quantity.NComp() == 3 {
		// vector field
		if req.Component == "All" {
			// do nothing
		}
	} else if quantity.NComp() == 1 {
		// scalar field
	} else {
		util.Fatal("Invalid quantity component count")
	}
	preview = Preview{
		Quantity:  quantity,
		Component: compStringToIndex(req.Component),
		Layer:     req.Layer,
		MaxPoints: req.MaxPoints,
	}
	util.Log("Preview state updated:", preview)

	return c.JSON(http.StatusOK, preview)
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

var PreviewBuffer []byte

// scaleOutputBuffer scales down the image size until the number of points are < MaxVectors
func scaleOutputBuffer(originalSize [3]int) [3]int {
	width, height := originalSize[0], originalSize[1]
	points := width * height

	// Calculate the scaling factor
	for points >= preview.MaxPoints {
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
	isVectorField := preview.Quantity.NComp() == 3 && preview.Component == -1

	nComp := 1
	if isVectorField {
		nComp = 3
	}
	util.Log("PreparePreviewBuffer: ", engine.NameOf(preview.Quantity), ", component:", preview.Component, ",layer: ", preview.Layer, ", Max points:", preview.MaxPoints, ", isVectorField:", isVectorField)

	originalSize := engine.MeshOf(preview.Quantity).Size()
	outputSize := scaleOutputBuffer(originalSize)
	scaledSlice := data.NewSlice(nComp, outputSize)
	tempBuffer := cuda.NewSlice(1, outputSize)

	GPUBuffer := engine.ValueOf(preview.Quantity)
	defer cuda.Recycle(GPUBuffer)
	defer tempBuffer.Free()

	if isVectorField {
		for c := 0; c < nComp; c++ {
			cuda.Resize(tempBuffer, GPUBuffer.Comp(c), preview.Layer)
			data.Copy(scaledSlice.Comp(c), tempBuffer)
		}
	} else {
		if preview.Quantity.NComp() == 3 {
			cuda.Resize(tempBuffer, GPUBuffer.Comp(nComp), preview.Layer)
			data.Copy(scaledSlice.Comp(0), tempBuffer)
		} else {
			cuda.Resize(tempBuffer, GPUBuffer.Comp(0), preview.Layer)
			data.Copy(scaledSlice.Comp(0), tempBuffer)
		}
	}

	if isVectorField {
		normalizeVectors(scaledSlice)
		vectorField := scaledSlice.Vectors()
		PreviewBuffer = ConvertVectorFieldToBinary(vectorField)
	} else {
		normalizeScalars(scaledSlice)
		scalarField := scaledSlice.Scalars()
		PreviewBuffer = ConvertScalarFieldToBinary(scalarField)
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
