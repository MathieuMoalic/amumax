package api

import (
	"math"
	"net/http"

	"github.com/MathieuMoalic/amumax/src/cuda"
	"github.com/MathieuMoalic/amumax/src/data"
	"github.com/MathieuMoalic/amumax/src/engine"
	"github.com/MathieuMoalic/amumax/src/log"
	"github.com/labstack/echo/v4"
)

type PreviewState struct {
	ws                   *WebSocketManager
	globalQuantities     []string
	mask                 [][][]float32
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
	NComp                int                  `msgpack:"nComp"`
	DataPointsCount      int                  `msgpack:"dataPointsCount"`
}

func initPreviewAPI(e *echo.Echo, ws *WebSocketManager) *PreviewState {
	previewState := &PreviewState{
		Quantity:             "m",
		Component:            "3D",
		Layer:                0,
		MaxPoints:            8192,
		Dimensions:           [3]int{0, 0, 0},
		Type:                 "3D",
		VectorFieldValues:    nil,
		VectorFieldPositions: nil,
		ScalarField:          nil,
		Min:                  0,
		Max:                  0,
		Refresh:              true,
		NComp:                3,
		DataPointsCount:      0,
		ws:                   ws,
		globalQuantities:     []string{"B_demag", "B_ext", "B_eff", "Edens_demag", "Edens_ext", "Edens_eff", "geom"},
	}
	e.POST("/api/preview/component", previewState.postPreviewComponent)
	e.POST("/api/preview/quantity", previewState.postPreviewQuantity)
	e.POST("/api/preview/layer", previewState.postPreviewLayer)
	e.POST("/api/preview/maxpoints", previewState.postPreviewMaxPoints)
	e.POST("/api/preview/refresh", previewState.postPreviewRefresh)

	return previewState
}

func (s *PreviewState) GetQuantity() engine.Quantity {
	quantity, exists := engine.Quantities[s.Quantity]
	if !exists {
		log.Log.Err("Quantity not found: %v", s.Quantity)
	}
	return quantity
}

func (s *PreviewState) GetComponent() int {
	return compStringToIndex(s.Component)
}

func (s *PreviewState) Update() {
	engine.InjectAndWait(s.UpdateQuantityBuffer)
}

func (s *PreviewState) UpdateQuantityBuffer() {
	// s.ScaleDimensions()
	if s.mask == nil {
		s.updateMask()
	}
	componentCount := 1
	if s.Type == "3D" {
		componentCount = 3
	}
	GPU_in := engine.ValueOf(s.GetQuantity())
	defer cuda.Recycle(GPU_in)

	CPU_out := data.NewSlice(componentCount, s.Dimensions)
	GPU_out := cuda.NewSlice(1, s.Dimensions)
	defer GPU_out.Free()

	if s.Type == "3D" {
		for c := 0; c < componentCount; c++ {
			cuda.Resize(GPU_out, GPU_in.Comp(c), s.Layer)
			data.Copy(CPU_out.Comp(c), GPU_out)
		}
		s.normalizeVectors(CPU_out)
		s.UpdateVectorField(CPU_out.Vectors())
	} else {
		if s.GetQuantity().NComp() > 1 {
			cuda.Resize(GPU_out, GPU_in.Comp(s.GetComponent()), s.Layer)
			data.Copy(CPU_out.Comp(0), GPU_out)
		} else {
			cuda.Resize(GPU_out, GPU_in.Comp(0), s.Layer)
			data.Copy(CPU_out.Comp(0), GPU_out)
		}
		s.UpdateScalarField(CPU_out.Scalars())
	}
}

func (s *PreviewState) ScaleDimensions() {
	originalSize := engine.MeshOf(s.GetQuantity()).Size()
	width, height := float32(originalSize[0]), float32(originalSize[1])
	points := width * height
	if points <= float32(s.MaxPoints) {
		s.Dimensions = [3]int{originalSize[0], originalSize[1], 1}
		return
	}

	// Calculate the scaling factor
	for points >= float32(s.MaxPoints) {
		width = width / 2
		height = height / 2
		points = width * height
	}
	s.Dimensions = [3]int{int(width), int(height), 1}
}

func (s *PreviewState) normalizeVectors(f *data.Slice) {
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

func (s *PreviewState) UpdateVectorField(vectorField [3][][][]float32) {
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
			if (valx == 0 && valy == 0 && valz == 0) || (math.IsNaN(float64(valx))) {
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
	s.VectorFieldPositions = posArray
	s.VectorFieldValues = valArray
	s.ScalarField = nil
	s.DataPointsCount = len(valArray)
}

func (s *PreviewState) UpdateScalarField(scalarField [][][]float32) {
	xLen := len(scalarField[0][0])
	yLen := len(scalarField[0])
	min, max := scalarField[0][0][0], scalarField[0][0][0]

	var valArray [][3]float32
	for posx := 0; posx < xLen; posx++ {
		for posy := 0; posy < yLen; posy++ {
			// Some quantities exist where the magnetic materials are not present
			// and we don't want to filter them out
			if !contains(s.globalQuantities, s.Quantity) {
				if s.mask[0][posy][posx] == 0 {
					continue
				}
			}
			val := scalarField[0][posy][posx]
			if val < min {
				min = val
			}
			if val > max {
				max = val
			}
			valArray = append(valArray, [3]float32{float32(posx), float32(posy), val})
		}
	}
	if len(valArray) == 0 {
		log.Log.Warn("No data in scalar field")
	}

	s.Min = min
	s.Max = max
	s.ScalarField = valArray
	s.VectorFieldValues = nil
	s.VectorFieldPositions = nil
	s.DataPointsCount = len(valArray)
}

func (s *PreviewState) updateMask() {
	s.ScaleDimensions()

	// cuda full size geom
	geom := engine.Geometry
	GPU_fullsize := cuda.Buffer(geom.NComp(), geom.Buffer.Size())
	geom.EvalTo(GPU_fullsize)
	defer cuda.Recycle(GPU_fullsize)

	// resize geom in GPU
	GPU_resized := cuda.NewSlice(1, s.Dimensions)
	defer GPU_resized.Free()
	cuda.Resize(GPU_resized, GPU_fullsize.Comp(0), 0)

	// copy resized geom from GPU to CPU
	CPU_out := data.NewSlice(1, s.Dimensions)
	defer CPU_out.Free()
	data.Copy(CPU_out.Comp(0), GPU_resized)

	// extract mask from CPU slice
	s.mask = CPU_out.Scalars()
}

func contains(arr []string, val string) bool {
	for _, item := range arr {
		if item == val {
			return true
		}
	}
	return false
}

func compStringToIndex(comp string) int {
	switch comp {
	case "x":
		return 0
	case "y":
		return 1
	case "z":
		return 2
	case "3D":
		return -1
	case "None":
		return 0
	}
	log.Log.ErrAndExit("Invalid component string")
	return -2
}

func (s *PreviewState) updatePreviewType() {
	var fieldType string
	isVectorField := s.NComp == 3 && s.GetComponent() == -1
	if isVectorField {
		fieldType = "3D"
	} else {
		fieldType = "2D"
	}
	if fieldType != s.Type {
		s.Type = fieldType
		s.Refresh = true
	}
}

func (s *PreviewState) validateComponent() {
	s.NComp = s.GetQuantity().NComp()
	switch s.NComp {
	case 1:
		s.Component = "None"
	case 3:
		if s.Component == "None" {
			s.Component = "3D"
		}
	default:
		log.Log.Err("Invalid number of components")
		// reset to default
		s.Quantity = "m"
		s.Component = "3D"
	}
}

func (s *PreviewState) postPreviewComponent(c echo.Context) error {
	type Request struct {
		Component string `msgpack:"component"`
	}
	req := new(Request)
	if err := c.Bind(req); err != nil {
		log.Log.Err("%v", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request payload"})
	}
	s.Component = req.Component
	s.validateComponent()
	s.Refresh = true
	s.updatePreviewType()
	s.ws.broadcastEngineState()
	return c.JSON(http.StatusOK, nil)
}

func (s *PreviewState) postPreviewQuantity(c echo.Context) error {
	type Request struct {
		Quantity string `msgpack:"quantity"`
	}
	req := new(Request)
	if err := c.Bind(req); err != nil {
		log.Log.Err("%v", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request payload"})
	}
	_, exists := engine.Quantities[req.Quantity]
	if !exists {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Quantity not found"})
	}
	s.Quantity = req.Quantity
	s.validateComponent()
	s.Refresh = true
	s.updatePreviewType()
	s.ws.broadcastEngineState()
	return c.JSON(http.StatusOK, nil)
}

func (s *PreviewState) postPreviewLayer(c echo.Context) error {
	type Request struct {
		Layer int `msgpack:"layer"`
	}
	req := new(Request)
	if err := c.Bind(req); err != nil {
		log.Log.Err("%v", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request payload"})
	}

	s.Layer = req.Layer
	s.Refresh = true
	s.ws.broadcastEngineState()
	return c.JSON(http.StatusOK, nil)
}

func (s *PreviewState) postPreviewMaxPoints(c echo.Context) error {
	type Request struct {
		MaxPoints int `msgpack:"maxPoints"`
	}
	req := new(Request)
	if err := c.Bind(req); err != nil {
		log.Log.Err("%v", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request payload"})
	}
	if req.MaxPoints < 8 {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "MaxPoints must be at least 8"})
	}
	s.MaxPoints = req.MaxPoints
	s.Refresh = true
	engine.InjectAndWait(s.updateMask)
	s.ws.broadcastEngineState()
	return c.JSON(http.StatusOK, nil)
}

func (s *PreviewState) postPreviewRefresh(c echo.Context) error {
	s.ws.broadcastEngineState()
	return c.JSON(http.StatusOK, nil)
}
