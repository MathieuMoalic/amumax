package api

import (
	"net/http"

	"encoding/binary"
	"math"

	"github.com/MathieuMoalic/amumax/cuda"
	"github.com/MathieuMoalic/amumax/data"
	"github.com/MathieuMoalic/amumax/engine"
	"github.com/labstack/echo/v4"
)

type Preview struct {
	Quantity  string `json:"displayQuantity"`
	Component int    `json:"displayComponent"`
	Layer     int    `json:"displayLayer"`
}

func postPreviewState(c echo.Context) error {
	req := new(Preview)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request payload"})
	}
	return c.JSON(http.StatusOK, "")
}

var DisplayVectorField []byte

// scaleOutputBuffer scales down the image size until the number of points are < 1000
func scaleOutputBuffer(originalSize [3]int, maxPoints int) [3]int {
	width, height := originalSize[0], originalSize[1]
	points := width * height

	// Calculate the scaling factor
	for points >= maxPoints {
		width = (width + 1) / 2
		height = (height + 1) / 2
		points = width * height
	}

	return [3]int{width, height, 1}
}

// GetVectorField gets the vector field from the mesh and converts it to a binary byte slice.
func GetVectorField() {
	const maxPoints = 1000
	quant := &engine.M
	renderLayer := 0

	originalSize := engine.MeshOf(quant).Size()
	outputSize := scaleOutputBuffer(originalSize, maxPoints)
	scaledSlice := data.NewSlice(3, outputSize)
	tempBuffer := cuda.NewSlice(1, outputSize)

	GPUBuffer := engine.ValueOf(quant)
	defer cuda.Recycle(GPUBuffer)
	defer tempBuffer.Free()

	for c := 0; c < quant.NComp(); c++ {
		cuda.Resize(tempBuffer, GPUBuffer.Comp(c), renderLayer)
		data.Copy(scaledSlice.Comp(c), tempBuffer)
	}
	engine.Normalize(scaledSlice)
	vectorField := scaledSlice.Vectors()
	DisplayVectorField = ConvertVectorFieldToBinary(vectorField)
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
