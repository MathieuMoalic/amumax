package api

import (
	"net/http"

	"github.com/MathieuMoalic/amumax/engine"
	"github.com/labstack/echo/v4"
)

type TablePlot struct {
	AutoSaveInterval float64     `msgpack:"autoSaveInterval"`
	Columns          []string    `msgpack:"columns"`
	XColumn          string      `msgpack:"xColumn"`
	YColumn          string      `msgpack:"yColumn"`
	XColumnUnit      string      `msgpack:"xColumnUnit"`
	YColumnUnit      string      `msgpack:"yColumnUnit"`
	Data             [][]float64 `msgpack:"data"`
	Min              float64     `msgpack:"min"`
	Max              float64     `msgpack:"max"`
}

func GetTablePlotData() ([][]float64, float64, float64) {
	xData := engine.Table.GetXData()
	yData := engine.Table.GetYData()
	// xUnit := engine.Table.
	data := make([][]float64, len(xData)) // [ [x1, y1], [x2, y2], ... ]
	if len(xData) == 0 {
		return data, 0, 0
	}
	min := yData[0]
	max := yData[0]
	for i := 0; i < len(xData); i++ {
		data[i] = []float64{xData[i], yData[i]}
		if min > yData[i] {
			min = yData[i]
		}
		if max < yData[i] {
			max = yData[i]
		}
	}
	return data, min, max
}
func newTablePlot() *TablePlot {
	data, min, max := GetTablePlotData()
	return &TablePlot{
		AutoSaveInterval: engine.Table.AutoSavePeriod,
		Columns:          engine.Table.GetTableNames(),
		XColumn:          engine.Table.XColumn,
		YColumn:          engine.Table.YColumn,
		XColumnUnit:      engine.Table.GetUnit(engine.Table.XColumn),
		YColumnUnit:      engine.Table.GetUnit(engine.Table.YColumn),
		Data:             data,
		Min:              min,
		Max:              max,
	}
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
	if !engine.Table.ColumnExists(req.XColumn) {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Table column not found"})
	}
	engine.Table.XColumn = req.XColumn
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
	if !engine.Table.ColumnExists(req.YColumn) {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Table column not found"})
	}
	engine.Table.YColumn = req.YColumn
	broadcastEngineState()
	return c.JSON(http.StatusOK, nil)
}
