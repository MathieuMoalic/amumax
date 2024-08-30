package api

import (
	"math"
	"net/http"

	"github.com/MathieuMoalic/amumax/cuda"
	"github.com/MathieuMoalic/amumax/data"
	"github.com/MathieuMoalic/amumax/engine"
	"github.com/MathieuMoalic/amumax/util"
	"github.com/labstack/echo/v4"
)

var preview Preview

func init() {
	preview = Preview{
		Quantity:             "m",
		Component:            "All",
		Layer:                0,
		MaxPoints:            10000,
		Dimensions:           [3]int{0, 0, 0},
		Type:                 "3D",
		VectorFieldValues:    nil,
		VectorFieldPositions: nil,
		ScalarField:          nil,
		Min:                  0,
		Max:                  0,
		Refresh:              true,
	}
}

type Preview struct {
	Quantity             string               `msgpack:"quantity"`
	Unit                 string               `msgpack:"unit"`
	Component            string               `msgpack:"component"`
	Layer                int                  `msgpack:"layer"`
	MaxPoints            int                  `msgpack:"maxPoints"`
	Dimensions           [3]int               `msgpack:"dimensions"`
	Type                 string               `msgpack:"type"`
	VectorFieldValues    []map[string]float32 `msgpack:"vectorFieldValues"`
	VectorFieldPositions []map[string]int     `msgpack:"vectorFieldPositions"`
	ScalarField          [][3]float32         `msgpack:"scalarField"`
	Min                  float32              `msgpack:"min"`
	Max                  float32              `msgpack:"max"`
	Refresh              bool                 `msgpack:"refresh"`
}

func (p *Preview) GetQuantity() engine.Quantity {
	quantity, exists := engine.Quantities[p.Quantity]
	if !exists {
		util.Log.Err("Quantity not found: %v", p.Quantity)
	}
	return quantity
}

func (p *Preview) GetComponent() int {
	return compStringToIndex(p.Component)
}

func (p *Preview) Update() {
	engine.InjectAndWait(p.UpdateQuantityBuffer)
}

func (p *Preview) UpdateQuantityBuffer() {
	nComp := 1
	if p.Type == "3D" {
		nComp = 3
	}
	p.ScaleDimensions()
	GPU_in := engine.ValueOf(p.GetQuantity())

	defer cuda.Recycle(GPU_in)
	CPU_out := data.NewSlice(nComp, p.Dimensions)
	GPU_out := cuda.NewSlice(1, p.Dimensions)
	defer GPU_out.Free()

	if p.Type == "3D" {
		for c := 0; c < nComp; c++ {
			cuda.Resize(GPU_out, GPU_in.Comp(c), p.Layer)
			data.Copy(CPU_out.Comp(c), GPU_out)
			// data.Copy(CPU_out.Comp(c), GPU_in.Comp(c))
		}
		p.normalizeVectors(CPU_out)
		p.UpdateVectorField(CPU_out.Vectors())
	} else {
		if p.GetQuantity().NComp() >= 1 {
			cuda.Resize(GPU_out, GPU_in.Comp(p.GetComponent()), p.Layer)
			data.Copy(CPU_out.Comp(0), GPU_out)
			// data.Copy(CPU_out.Comp(0), GPU_in.Comp(p.GetComponent()))
		} else {
			cuda.Resize(GPU_out, GPU_in.Comp(0), p.Layer)
			data.Copy(CPU_out.Comp(0), GPU_out)
			// data.Copy(CPU_out.Comp(0), GPU_in.Comp(0))
		}
		p.UpdateScalarField(CPU_out.Scalars())
	}
}

func (p *Preview) ScaleDimensions() {
	originalSize := engine.MeshOf(p.GetQuantity()).Size()
	width, height := float32(originalSize[0]), float32(originalSize[1])
	points := width * height
	if points <= float32(p.MaxPoints) {
		p.Dimensions = [3]int{originalSize[0], originalSize[1], 1}
		return
	}

	// Calculate the scaling factor
	for points >= float32(p.MaxPoints) {
		width = width / 2
		height = height / 2
		points = width * height
	}
	p.Dimensions = [3]int{int(width), int(height), 1}
}

func (p *Preview) normalizeVectors(f *data.Slice) {
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

func (p *Preview) UpdateVectorField(vectorField [3][][][]float32) {
	// Calculate the total number of elements
	yLen := len(vectorField[0][0])
	xLen := len(vectorField[0][0][0])

	// Create a slice to hold the array of {x, y, z} objects
	var valArray []map[string]float32
	var posArray []map[string]int

	for posx := 0; posx < xLen; posx++ {
		for posy := 0; posy < yLen; posy++ {
			valx := vectorField[0][0][posy][posx]
			valy := vectorField[1][0][posy][posx]
			valz := vectorField[2][0][posy][posx]
			if valx == 0 && valy == 0 && valz == 0 {
				continue
			}
			posArray = append(posArray,
				map[string]int{
					"x": posx,
					"y": posy,
					"z": 0,
				})
			valArray = append(valArray,
				map[string]float32{
					"x": valx,
					"y": valy,
					"z": valz,
				})
		}
	}
	p.VectorFieldPositions = posArray
	p.VectorFieldValues = valArray
	p.ScalarField = nil
}

func (p *Preview) UpdateScalarField(scalarField [][][]float32) {
	xLen := len(scalarField[0][0])
	yLen := len(scalarField[0])
	min, max := scalarField[0][0][0], scalarField[0][0][0]
	var valArray [][3]float32
	for posx := 0; posx < xLen; posx++ {
		for posy := 0; posy < yLen; posy++ {
			val := scalarField[0][posy][posx]
			if val == 0 {
				continue
			}
			if val < min {
				min = val
			}
			if val > max {
				max = val
			}
			valArray = append(valArray, [3]float32{float32(posx), float32(posy), val})
		}
	}
	p.Min = min
	p.Max = max
	p.ScalarField = valArray
	p.VectorFieldValues = nil
	p.VectorFieldPositions = nil
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
	util.Log.ErrAndExit("Invalid component string")
	return -2
}

func newPreview() *Preview {
	preview.Update()
	return &preview
}

func updatePreviewType() {
	var fieldType string
	isVectorField := preview.GetQuantity().NComp() == 3 && preview.GetComponent() == -1
	if isVectorField {
		fieldType = "3D"
	} else {
		fieldType = "2D"
	}
	if fieldType != preview.Type {
		preview.Type = fieldType
		preview.Refresh = true
	}
}

func postPreviewComponent(c echo.Context) error {
	type Request struct {
		Component string `msgpack:"component"`
	}
	req := new(Request)
	if err := c.Bind(req); err != nil {
		util.Log.Err("%v", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request payload"})
	}
	if preview.GetQuantity().NComp() == 1 && req.Component == "All" {
		preview.Component = "x"
	} else {
		preview.Component = req.Component
	}
	preview.Refresh = true
	updatePreviewType()
	broadcastEngineState()
	return c.JSON(http.StatusOK, nil)
}

func postPreviewQuantity(c echo.Context) error {
	type Request struct {
		Quantity string `msgpack:"quantity"`
	}
	req := new(Request)
	if err := c.Bind(req); err != nil {
		util.Log.Err("%v", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request payload"})
	}
	_, exists := engine.Quantities[req.Quantity]
	if !exists {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Quantity not found"})
	}
	preview.Quantity = req.Quantity
	if preview.GetQuantity().NComp() == 1 && preview.Component == "All" {
		preview.Component = "x"
	}
	updatePreviewType()
	broadcastEngineState()
	return c.JSON(http.StatusOK, nil)
}

func postPreviewLayer(c echo.Context) error {
	type Request struct {
		Layer int `msgpack:"layer"`
	}
	req := new(Request)
	if err := c.Bind(req); err != nil {
		util.Log.Err("%v", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request payload"})
	}
	preview.Layer = req.Layer
	preview.Refresh = true
	broadcastEngineState()
	return c.JSON(http.StatusOK, nil)
}

func postPreviewMaxPoints(c echo.Context) error {
	type Request struct {
		MaxPoints int `msgpack:"maxPoints"`
	}
	req := new(Request)
	if err := c.Bind(req); err != nil {
		util.Log.Err("%v", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request payload"})
	}
	preview.MaxPoints = req.MaxPoints
	preview.Refresh = true
	broadcastEngineState()
	return c.JSON(http.StatusOK, nil)
}

func postPreviewRefresh(c echo.Context) error {
	broadcastEngineState()
	return c.JSON(http.StatusOK, nil)
}
