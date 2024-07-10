package api

import (
	"net/http"

	"github.com/MathieuMoalic/amumax/engine"
	"github.com/labstack/echo/v4"
)

type TablePlotData struct {
	X float64 `msgpack:"x"`
	Y float64 `msgpack:"y"`
}
type TablePlot struct {
	AutoSaveInterval float64         `msgpack:"autoSaveInterval"`
	Columns          []string        `msgpack:"columns"`
	XColumn          string          `msgpack:"xColumn"`
	YColumn          string          `msgpack:"yColumn"`
	Data             []TablePlotData `msgpack:"data"`
}

func newTablePlot() *TablePlot {
	tablePlotData := make([]TablePlotData, len(engine.Table.Data[engine.Tableplot.X]))
	for i := range tablePlotData {
		tablePlotData[i] = TablePlotData{
			X: engine.Table.Data[engine.Tableplot.X][i],
			Y: engine.Table.Data[engine.Tableplot.Y][i],
		}
	}
	data := TablePlot{
		AutoSaveInterval: engine.Table.AutoSavePeriod,
		Columns:          engine.Table.GetTableNames(),
		XColumn:          engine.Tableplot.X,
		YColumn:          engine.Tableplot.Y,
		Data:             tablePlotData,
	}
	return &data
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
	engine.Tableplot.X = req.XColumn
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
	engine.Tableplot.Y = req.YColumn
	broadcastEngineState()
	return c.JSON(http.StatusOK, nil)
}
