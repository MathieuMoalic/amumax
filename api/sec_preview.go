package api

import (
	"encoding/binary"
	"net/http"
	"strconv"

	"math"

	"github.com/MathieuMoalic/amumax/cuda"
	"github.com/MathieuMoalic/amumax/data"
	"github.com/MathieuMoalic/amumax/engine"
	"github.com/MathieuMoalic/amumax/util"
	"github.com/labstack/echo/v4"
)

type Preview struct {
	Quantity   string  `msgpack:"quantity"`
	Component  string  `msgpack:"component"`
	Layer      int     `msgpack:"layer"`
	MaxPoints  int     `msgpack:"maxPoints"`
	Dimensions [3]int  `msgpack:"dimensions"`
	Type       string  `msgpack:"type"`
	Buffer     *[]byte `msgpack:"buffer"`
}

func newPreview() *Preview {
	return &Preview{
		Quantity:   engine.NameOf(backendPreviewState.Quantity),
		Component:  compIndexToString(backendPreviewState.Component),
		Layer:      backendPreviewState.Layer,
		MaxPoints:  backendPreviewState.MaxPoints,
		Dimensions: backendPreviewState.Dimensions,
		Type:       backendPreviewState.Type,
		Buffer:     &backendPreviewState.Buffer,
	}
}

var backendPreviewState BackendPreviewState

func init() {
	backendPreviewState = BackendPreviewState{
		Quantity:   &engine.M,
		Component:  0,
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

func postPreviewComponent(c echo.Context) error {
	type Request struct {
		Component string `msgpack:"component"`
	}
	req := new(Request)
	if err := c.Bind(req); err != nil {
		util.LogErr(err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request payload"})
	}
	backendPreviewState.Component = compStringToIndex(req.Component)
	broadcastEngineState()
	return c.JSON(http.StatusOK, nil)
}

func postPreviewQuantity(c echo.Context) error {
	type Request struct {
		Quantity string `msgpack:"quantity"`
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
	backendPreviewState.Quantity = quantity
	broadcastEngineState()
	return c.JSON(http.StatusOK, nil)
}

func postPreviewLayer(c echo.Context) error {
	type Request struct {
		Layer int `msgpack:"layer"`
	}
	req := new(Request)
	if err := c.Bind(req); err != nil {
		util.LogErr(err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request payload"})
	}
	backendPreviewState.Layer = req.Layer
	broadcastEngineState()
	return c.JSON(http.StatusOK, nil)
}

func postPreviewMaxPoints(c echo.Context) error {
	type Request struct {
		MaxPoints string `msgpack:"maxPoints"`
	}
	req := new(Request)
	if err := c.Bind(req); err != nil {
		util.LogErr(err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request payload"})
	}
	MaxPoints, err := strconv.Atoi(req.MaxPoints)
	if err != nil {
		util.LogErr(err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Could not parse maxPoints as integer"})
	}
	backendPreviewState.MaxPoints = MaxPoints
	broadcastEngineState()
	return c.JSON(http.StatusOK, nil)
}

func postPreviewRefresh(c echo.Context) error {
	broadcastEngineState()
	return c.JSON(http.StatusOK, nil)
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
	for points >= backendPreviewState.MaxPoints {
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
	isVectorField := backendPreviewState.Quantity.NComp() == 3 && backendPreviewState.Component == -1

	nComp := 1
	if isVectorField {
		nComp = 3
	}
	util.Log("PreparePreviewBuffer: ", engine.NameOf(backendPreviewState.Quantity), ", component:", backendPreviewState.Component, ",layer: ", backendPreviewState.Layer, ", Max points:", backendPreviewState.MaxPoints, ", isVectorField:", isVectorField)

	originalSize := engine.MeshOf(backendPreviewState.Quantity).Size()
	outputSize := scaleOutputBuffer(originalSize)
	backendPreviewState.Dimensions = outputSize
	scaledSlice := data.NewSlice(nComp, outputSize)
	tempBuffer := cuda.NewSlice(1, outputSize)

	GPUBuffer := engine.ValueOf(backendPreviewState.Quantity)
	defer cuda.Recycle(GPUBuffer)
	defer tempBuffer.Free()

	if isVectorField {
		for c := 0; c < nComp; c++ {
			cuda.Resize(tempBuffer, GPUBuffer.Comp(c), backendPreviewState.Layer)
			data.Copy(scaledSlice.Comp(c), tempBuffer)
		}
	} else {
		if backendPreviewState.Quantity.NComp() == 3 {
			cuda.Resize(tempBuffer, GPUBuffer.Comp(nComp), backendPreviewState.Layer)
			data.Copy(scaledSlice.Comp(0), tempBuffer)
		} else {
			cuda.Resize(tempBuffer, GPUBuffer.Comp(0), backendPreviewState.Layer)
			data.Copy(scaledSlice.Comp(0), tempBuffer)
		}
	}

	if isVectorField {
		normalizeVectors(scaledSlice)
		backendPreviewState.Buffer = ConvertVectorFieldToBinary(scaledSlice.Vectors())
	} else {
		normalizeScalars(scaledSlice)
		backendPreviewState.Buffer = ConvertScalarFieldToBinary(scaledSlice.Scalars())
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
