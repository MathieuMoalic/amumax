package api

import (
	"net/http"

	"github.com/MathieuMoalic/amumax/src/engine"
	"github.com/MathieuMoalic/amumax/src/engine/log"
	"github.com/labstack/echo/v4"
)

type TablePlotState struct {
	ws               *WebSocketManager
	AutoSaveInterval *float64    `msgpack:"autoSaveInterval"`
	Columns          []string    `msgpack:"columns"`
	XColumn          string      `msgpack:"xColumn"`
	YColumn          string      `msgpack:"yColumn"`
	XColumnUnit      string      `msgpack:"xColumnUnit"`
	YColumnUnit      string      `msgpack:"yColumnUnit"`
	Data             [][]float64 `msgpack:"data"`
	XMin             float64     `msgpack:"xmin"`
	XMax             float64     `msgpack:"xmax"`
	YMin             float64     `msgpack:"ymin"`
	YMax             float64     `msgpack:"ymax"`
	MaxPoints        int         `msgpack:"maxPoints"`
	Step             int         `msgpack:"step"`
}

func initTablePlotAPI(e *echo.Group, ws *WebSocketManager) *TablePlotState {
	t := TablePlotState{
		ws:               ws,
		AutoSaveInterval: &engine.Table.AutoSavePeriod,
		XColumn:          "t",
		YColumn:          "mx",
		MaxPoints:        10000,
		Step:             1,
	}
	e.POST("/api/tableplot/autosave-interval", t.postTablePlotAutoSaveInterval)
	e.POST("/api/tableplot/xcolumn", t.postTablePlotXColumn)
	e.POST("/api/tableplot/ycolumn", t.postTablePlotYColumn)
	e.POST("/api/tableplot/maxpoints", t.postTablePlotMaxPoints)
	e.POST("/api/tableplot/step", t.postTablePlotStep)
	t.Update()
	return &t
}

func (t *TablePlotState) Update() {
	t.Columns = t.GetTableNames()
	t.XColumnUnit = t.GetUnit(t.XColumn)
	t.YColumnUnit = t.GetUnit(t.YColumn)
	t.GetTablePlotData()
	t.GetMinMaxXY()
}

func (t *TablePlotState) GetTablePlotData() {
	xData := t.GetColumnData(t.XColumn)
	yData := t.GetColumnData(t.YColumn)
	data := make([][]float64, len(xData)) // [ [x1, y1], [x2, y2], ... ]
	if len(xData) == 0 {
		t.Data = data
		return
	}
	if len(xData) > t.MaxPoints {
		xData = xData[len(xData)-t.MaxPoints:]
		yData = yData[len(yData)-t.MaxPoints:]
	}
	for i := 0; i < len(xData); i++ {
		data[i] = []float64{xData[i], yData[i]}
	}
	t.Data = data
}

func (t *TablePlotState) GetMinMaxXY() {
	if len(t.Data) == 0 {
		return
	}
	// Initialize with the first valid data point
	xmin, xmax := t.Data[0][0], t.Data[0][0]
	ymin, ymax := t.Data[0][1], t.Data[0][1]
	for _, i := range t.Data {
		if len(i) < 2 {
			continue // Skip this entry if it doesn't have enough elements
		}
		if i[0] < xmin {
			xmin = i[0]
		}
		if i[0] > xmax {
			xmax = i[0]
		}
		if i[1] < ymin {
			ymin = i[1]
		}
		if i[1] > ymax {
			ymax = i[1]
		}
	}
	t.YMin = ymin
	t.YMax = ymax
	t.XMin = xmin
	t.XMax = xmax
}

func (t *TablePlotState) GetColumnData(column string) []float64 {
	engine.Table.Mu.Lock() // Lock the mutex before reading the map
	defer engine.Table.Mu.Unlock()
	originalData := engine.Table.Data[column]
	originalLen := len(originalData)
	newLen := (originalLen + 1) / t.Step
	result := make([]float64, 0, newLen)
	for i := 0; i < originalLen; i += t.Step {
		result = append(result, originalData[i])
	}
	return result
}

func (t *TablePlotState) GetUnit(name string) string {
	for _, i := range engine.Table.Columns {
		if i.Name == name {
			return i.Unit
		}
	}
	return ""
}

func (t *TablePlotState) ColumnExists(name string) bool {
	for _, i := range engine.Table.Columns {
		if i.Name == name {
			return true
		}
	}
	return false
}

func (t *TablePlotState) GetTableNames() []string {
	names := []string{}
	for _, column := range engine.Table.Columns {
		names = append(names, column.Name)
	}
	return names
}

func (t *TablePlotState) postTablePlotAutoSaveInterval(c echo.Context) error {
	type Request struct {
		AutoSaveInterval string `msgpack:"autoSaveInterval"`
	}
	req := new(Request)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request payload"})
	}
	engine.InjectAndWait(func() { engine.EvalTryRecover("TableAutoSave(" + req.AutoSaveInterval + ")") })
	t.ws.broadcastEngineState()
	return c.JSON(http.StatusOK, nil)
}

func (t *TablePlotState) postTablePlotXColumn(c echo.Context) error {
	type Request struct {
		XColumn string `msgpack:"XColumn"`
	}
	req := new(Request)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request payload"})
	}
	if !t.ColumnExists(req.XColumn) {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Table column not found"})
	}
	t.XColumn = req.XColumn
	t.ws.broadcastEngineState()
	return c.JSON(http.StatusOK, nil)
}

func (t *TablePlotState) postTablePlotYColumn(c echo.Context) error {
	type Request struct {
		YColumn string `msgpack:"YColumn"`
	}
	req := new(Request)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request payload"})
	}
	if !t.ColumnExists(req.YColumn) {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Table column not found"})
	}
	t.YColumn = req.YColumn
	t.ws.broadcastEngineState()
	return c.JSON(http.StatusOK, nil)
}

func (t *TablePlotState) postTablePlotMaxPoints(c echo.Context) error {
	type Request struct {
		MaxPoints int `msgpack:"maxPoints"`
	}
	req := new(Request)
	if err := c.Bind(req); err != nil {
		log.Log.Err("%v", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request payload"})
	}
	t.MaxPoints = req.MaxPoints
	t.ws.engineState.Preview.Refresh = true
	t.ws.broadcastEngineState()
	return c.JSON(http.StatusOK, nil)
}

func (t *TablePlotState) postTablePlotStep(c echo.Context) error {
	type Request struct {
		Step int `msgpack:"step"`
	}
	req := new(Request)
	if err := c.Bind(req); err != nil {
		log.Log.Err("%v", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request payload"})
	}
	if req.Step < 1 {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Step must be at least 1"})
	}
	t.Step = req.Step
	t.ws.broadcastEngineState()
	return c.JSON(http.StatusOK, nil)
}
