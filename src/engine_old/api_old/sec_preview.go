package api_old

import (
	"math"
	"net/http"
	"time"

	"github.com/MathieuMoalic/amumax/src/cuda_old"
	"github.com/MathieuMoalic/amumax/src/engine_old"
	"github.com/MathieuMoalic/amumax/src/engine_old/data_old"
	"github.com/MathieuMoalic/amumax/src/engine_old/log_old"
	"github.com/labstack/echo/v4"
)

type PreviewState struct {
	ws                   *WebSocketManager
	globalQuantities     []string
	layerMask            [][]float32
	Quantity             string               `msgpack:"quantity"`
	Unit                 string               `msgpack:"unit"`
	Component            string               `msgpack:"component"`
	Layer                int                  `msgpack:"layer"`
	Type                 string               `msgpack:"type"`
	VectorFieldValues    []map[string]float32 `msgpack:"vectorFieldValues"`
	VectorFieldPositions []map[string]int     `msgpack:"vectorFieldPositions"`
	ScalarField          [][3]float32         `msgpack:"scalarField"`
	Min                  float32              `msgpack:"min"`
	Max                  float32              `msgpack:"max"`
	Refresh              bool                 `msgpack:"refresh"`
	NComp                int                  `msgpack:"nComp"`

	MaxPoints       int   `msgpack:"maxPoints"`
	DataPointsCount int   `msgpack:"dataPointsCount"`
	XPossibleSizes  []int `msgpack:"xPossibleSizes"`
	YPossibleSizes  []int `msgpack:"yPossibleSizes"`
	XChosenSize     int   `msgpack:"xChosenSize"`
	YChosenSize     int   `msgpack:"yChosenSize"`
}

func initPreviewAPI(e *echo.Group, ws *WebSocketManager) *PreviewState {
	previewState := &PreviewState{
		Quantity:             "m",
		Component:            "3D",
		Layer:                0,
		MaxPoints:            8192,
		Type:                 "3D",
		VectorFieldValues:    nil,
		VectorFieldPositions: nil,
		ScalarField:          nil,
		Min:                  0,
		Max:                  0,
		Refresh:              true,
		NComp:                3,
		DataPointsCount:      0,
		XPossibleSizes:       nil,
		YPossibleSizes:       nil,
		XChosenSize:          engine_old.Mesh.Nx,
		YChosenSize:          engine_old.Mesh.Ny,
		ws:                   ws,
		globalQuantities:     []string{"B_demag", "B_ext", "B_eff", "Edens_demag", "Edens_ext", "Edens_eff", "geom"},
	}
	previewState.addPossibleDownscaleSizes()
	e.POST("/api/preview/component", previewState.postPreviewComponent)
	e.POST("/api/preview/quantity", previewState.postPreviewQuantity)
	e.POST("/api/preview/layer", previewState.postPreviewLayer)
	e.POST("/api/preview/maxpoints", previewState.postPreviewMaxPoints)
	e.POST("/api/preview/refresh", previewState.postPreviewRefresh)
	e.POST("/api/preview/XChosenSize", previewState.postXChosenSize)
	e.POST("/api/preview/YChosenSize", previewState.postYChosenSize)

	return previewState
}

func (s *PreviewState) getQuantity() engine_old.Quantity {
	quantity, exists := engine_old.Quantities[s.Quantity]
	if !exists {
		log_old.Log.Err("Quantity not found: %v", s.Quantity)
	}
	return quantity
}

func (s *PreviewState) getComponent() int {
	return compStringToIndex(s.Component)
}

func (s *PreviewState) Update() {
	engine_old.InjectAndWait(s.UpdateQuantityBuffer)
}

func (s *PreviewState) UpdateQuantityBuffer() {
	defer func() {
		if r := recover(); r != nil {
			log_old.Log.Warn("Recovered from panic in UpdateQuantityBuffer: %v", r)
			s.ScalarField = nil
			s.VectorFieldPositions = nil
			s.VectorFieldValues = nil
		}
	}()
	if s.layerMask == nil {
		s.updateMask()
	}
	if s.XChosenSize == 0 || s.YChosenSize == 0 {
		log_old.Log.Debug("XChosenSize or YChosenSize is 0")
		return
	}
	componentCount := 1
	if s.Type == "3D" {
		componentCount = 3
	}
	GPU_in := engine_old.ValueOf(s.getQuantity())
	defer cuda_old.Recycle(GPU_in)

	CPU_out := data_old.NewSlice(componentCount, [3]int{s.XChosenSize, s.YChosenSize, 1})
	GPU_out := cuda_old.NewSlice(1, [3]int{s.XChosenSize, s.YChosenSize, 1})
	defer GPU_out.Free()

	if s.Type == "3D" {
		for c := 0; c < componentCount; c++ {
			cuda_old.Resize(GPU_out, GPU_in.Comp(c), s.Layer)
			data_old.Copy(CPU_out.Comp(c), GPU_out)
		}
		s.normalizeVectors(CPU_out)
		s.UpdateVectorField(CPU_out.Vectors())
	} else {
		if s.getQuantity().NComp() > 1 {
			cuda_old.Resize(GPU_out, GPU_in.Comp(s.getComponent()), s.Layer)
			data_old.Copy(CPU_out.Comp(0), GPU_out)
		} else {
			cuda_old.Resize(GPU_out, GPU_in.Comp(0), s.Layer)
			data_old.Copy(CPU_out.Comp(0), GPU_out)
		}
		s.UpdateScalarField(CPU_out.Scalars())
	}
}

func (s *PreviewState) normalizeVectors(f *data_old.Slice) {
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
				if s.layerMask != nil {
					if s.layerMask[posy][posx] == 0 {
						continue
					}
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
		log_old.Log.Warn("No data in scalar field")
	}

	s.Min = min
	s.Max = max
	s.ScalarField = valArray
	s.VectorFieldValues = nil
	s.VectorFieldPositions = nil
	s.DataPointsCount = len(valArray)
}

func (s *PreviewState) updateMask() {
	defer func() {
		if r := recover(); r != nil {
			// Handle the panic, possibly log the error
			log_old.Log.Warn("Recovered from panic in updateMask: %v", r)
			// Optionally, reset s.layerMask or handle it appropriately
			s.layerMask = nil
		}
	}()
	if s.XChosenSize == 0 || s.YChosenSize == 0 {
		log_old.Log.Debug("XChosenSize or YChosenSize is 0")
		return
	}
	// cuda full size geom
	geom := engine_old.Geometry
	GPU_fullsize := cuda_old.Buffer(geom.NComp(), geom.Buffer.Size())
	geom.EvalTo(GPU_fullsize)
	defer cuda_old.Recycle(GPU_fullsize)

	// resize geom in GPU
	GPU_resized := cuda_old.NewSlice(1, [3]int{s.XChosenSize, s.YChosenSize, 1})
	defer GPU_resized.Free()
	cuda_old.Resize(GPU_resized, GPU_fullsize.Comp(0), s.Layer)

	// copy resized geom from GPU to CPU
	CPU_out := data_old.NewSlice(1, [3]int{s.XChosenSize, s.YChosenSize, 1})
	defer CPU_out.Free()
	data_old.Copy(CPU_out.Comp(0), GPU_resized)

	// extract mask from CPU slice
	s.layerMask = CPU_out.Scalars()[0]
}

func contains(arr []string, val string) bool {
	for _, item := range arr {
		if item == val {
			return true
		}
	}
	return false
}

func closestInArray(arr []int, target int) int {
	closest := arr[0]
	minDiff := math.Abs(float64(target - closest))

	for _, num := range arr {
		diff := math.Abs(float64(target - num))
		if diff < minDiff {
			minDiff = diff
			closest = num
		}
	}

	return closest
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
	log_old.Log.ErrAndExit("Invalid component string")
	return -2
}

// A valid destination size is a positive integer less than or equal to srcsize that evenly divides srcsize.
func (s *PreviewState) addPossibleDownscaleSizes() {
	// retry until engine.Mesh.Nx and engine.Mesh.Ny are not 0
	for engine_old.Mesh.Nx == 0 || engine_old.Mesh.Ny == 0 {
		time.Sleep(1 * time.Second)
	}
	if engine_old.Mesh.Nx == 0 || engine_old.Mesh.Ny == 0 {
		log_old.Log.Err("Nx or Ny is 0")
	}
	// iterate over engine.Mesh.Nx and engine.Mesh.Ny
	for dstsize := 1; dstsize <= engine_old.Mesh.Nx; dstsize++ {
		if engine_old.Mesh.Nx%dstsize == 0 {
			s.XPossibleSizes = append(s.XPossibleSizes, dstsize)
		}
	}
	for dstsize := 1; dstsize <= engine_old.Mesh.Ny; dstsize++ {
		if engine_old.Mesh.Ny%dstsize == 0 {
			s.YPossibleSizes = append(s.YPossibleSizes, dstsize)
		}
	}
	if len(s.YPossibleSizes) == 0 || len(s.XPossibleSizes) == 0 {
		log_old.Log.Err("No possible sizes found")
	}
	if engine_old.PreviewXDataPoints != 0 {
		s.XChosenSize = closestInArray(s.XPossibleSizes, engine_old.PreviewXDataPoints)
	} else {
		s.XChosenSize = closestInArray(s.XPossibleSizes, 100)
	}
	if engine_old.PreviewYDataPoints != 0 {
		s.YChosenSize = closestInArray(s.YPossibleSizes, engine_old.PreviewYDataPoints)
	} else {
		s.YChosenSize = closestInArray(s.YPossibleSizes, 100)
	}
}

func (s *PreviewState) updatePreviewType() {
	var fieldType string
	isVectorField := s.NComp == 3 && s.getComponent() == -1
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
	s.NComp = s.getQuantity().NComp()
	switch s.NComp {
	case 1:
		s.Component = "None"
	case 3:
		if s.Component == "None" {
			s.Component = "3D"
		}
	default:
		log_old.Log.Err("Invalid number of components")
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
		log_old.Log.Err("%v", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request payload"})
	}
	s.Component = req.Component
	s.validateComponent()
	s.updatePreviewType()
	s.Refresh = true
	s.ws.broadcastEngineState()
	return c.JSON(http.StatusOK, nil)
}

func (s *PreviewState) postPreviewQuantity(c echo.Context) error {
	type Request struct {
		Quantity string `msgpack:"quantity"`
	}
	req := new(Request)
	if err := c.Bind(req); err != nil {
		log_old.Log.Err("%v", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request payload"})
	}
	_, exists := engine_old.Quantities[req.Quantity]
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
		log_old.Log.Err("%v", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request payload"})
	}

	s.Layer = req.Layer
	s.Refresh = true
	engine_old.InjectAndWait(s.updateMask)
	s.ws.broadcastEngineState()
	return c.JSON(http.StatusOK, nil)
}

func (s *PreviewState) postPreviewMaxPoints(c echo.Context) error {
	type Request struct {
		MaxPoints int `msgpack:"maxPoints"`
	}
	req := new(Request)
	if err := c.Bind(req); err != nil {
		log_old.Log.Err("%v", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request payload"})
	}
	if req.MaxPoints < 8 {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "MaxPoints must be at least 8"})
	}
	s.MaxPoints = req.MaxPoints
	s.Refresh = true
	engine_old.InjectAndWait(s.updateMask)
	s.ws.broadcastEngineState()
	return c.JSON(http.StatusOK, nil)
}

func (s *PreviewState) postPreviewRefresh(c echo.Context) error {
	s.ws.broadcastEngineState()
	return c.JSON(http.StatusOK, nil)
}

func containsInt(arr []int, target int) bool {
	for _, v := range arr {
		if v == target {
			return true
		}
	}
	return false
}

func (s *PreviewState) postXChosenSize(c echo.Context) error {
	type Request struct {
		XChosenSize int `msgpack:"xChosenSize"`
	}
	req := new(Request)
	if err := c.Bind(req); err != nil {
		log_old.Log.Err("%v", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request payload"})
	}
	if !containsInt(s.XPossibleSizes, req.XChosenSize) {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid xChosenSize"})
	}
	s.XChosenSize = req.XChosenSize
	s.Refresh = true
	engine_old.InjectAndWait(s.updateMask)
	s.ws.broadcastEngineState()
	return c.JSON(http.StatusOK, nil)
}

func (s *PreviewState) postYChosenSize(c echo.Context) error {
	type Request struct {
		YChosenSize int `msgpack:"yChosenSize"`
	}
	req := new(Request)
	if err := c.Bind(req); err != nil {
		log_old.Log.Err("%v", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request payload"})
	}
	if !containsInt(s.YPossibleSizes, req.YChosenSize) {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid yChosenSize"})
	}
	s.YChosenSize = req.YChosenSize
	s.Refresh = true
	engine_old.InjectAndWait(s.updateMask)
	s.ws.broadcastEngineState()
	return c.JSON(http.StatusOK, nil)
}
