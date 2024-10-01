package api

import (
	"net/http"

	"github.com/MathieuMoalic/amumax/src/engine"
	"github.com/MathieuMoalic/amumax/src/util"
	"github.com/labstack/echo/v4"
)

func newTablePlot() *TablePlot {
	tablePlot.Update()
	return &tablePlot
}

type TablePlot struct {
	AutoSaveInterval float64     `msgpack:"autoSaveInterval"`
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

var tablePlot TablePlot

func init() {
	tablePlot = TablePlot{
		AutoSaveInterval: engine.Table.AutoSavePeriod,
		XColumn:          "t",
		YColumn:          "mx",
		MaxPoints:        10000,
		Step:             1,
	}
	tablePlot.Update()
}

func (t *TablePlot) Update() {
	t.AutoSaveInterval = engine.Table.AutoSavePeriod
	t.Columns = t.GetTableNames()
	t.XColumnUnit = t.GetUnit(t.XColumn)
	t.YColumnUnit = t.GetUnit(t.YColumn)
	t.GetTablePlotData()
	t.GetMinMaxXY()
}

func (t *TablePlot) GetTablePlotData() {
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

func (t *TablePlot) GetMinMaxXY() {
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

func (t *TablePlot) GetColumnData(column string) []float64 {
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

func (t *TablePlot) GetUnit(name string) string {
	for _, i := range engine.Table.Columns {
		if i.Name == name {
			return i.Unit
		}
	}
	return ""
}

func (t *TablePlot) ColumnExists(name string) bool {
	for _, i := range engine.Table.Columns {
		if i.Name == name {
			return true
		}
	}
	return false
}

func (t *TablePlot) GetTableNames() []string {
	names := []string{}
	for _, column := range engine.Table.Columns {
		names = append(names, column.Name)
	}
	return names
}

func postTablePlotAutoSaveInterval(c echo.Context) error {
	type Request struct {
		AutoSaveInterval float64 `msgpack:"autoSaveInterval"`
	}
	req := new(Request)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request payload"})
	}
	engine.Table.AutoSavePeriod = req.AutoSaveInterval
	broadcastEngineState()
	return c.JSON(http.StatusOK, nil)
}

func postTablePlotXColumn(c echo.Context) error {
	type Request struct {
		XColumn string `msgpack:"XColumn"`
	}
	req := new(Request)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request payload"})
	}
	if !tablePlot.ColumnExists(req.XColumn) {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Table column not found"})
	}
	tablePlot.XColumn = req.XColumn
	broadcastEngineState()
	return c.JSON(http.StatusOK, nil)
}

func postTablePlotYColumn(c echo.Context) error {
	type Request struct {
		YColumn string `msgpack:"YColumn"`
	}
	req := new(Request)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request payload"})
	}
	if !tablePlot.ColumnExists(req.YColumn) {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Table column not found"})
	}
	tablePlot.YColumn = req.YColumn
	broadcastEngineState()
	return c.JSON(http.StatusOK, nil)
}

func postTablePlotMaxPoints(c echo.Context) error {
	type Request struct {
		MaxPoints int `msgpack:"maxPoints"`
	}
	req := new(Request)
	if err := c.Bind(req); err != nil {
		util.Log.Err("%v", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request payload"})
	}
	tablePlot.MaxPoints = req.MaxPoints
	preview.Refresh = true
	broadcastEngineState()
	return c.JSON(http.StatusOK, nil)
}

func postTablePlotStep(c echo.Context) error {
	type Request struct {
		Step int `msgpack:"step"`
	}
	req := new(Request)
	if err := c.Bind(req); err != nil {
		util.Log.Err("%v", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request payload"})
	}
	if req.Step < 1 {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Step must be at least 1"})
	}
	tablePlot.Step = req.Step
	broadcastEngineState()
	return c.JSON(http.StatusOK, nil)
}
