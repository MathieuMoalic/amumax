package api

import (
	"net/http"

	"github.com/MathieuMoalic/amumax/engine"
	"github.com/MathieuMoalic/amumax/util"
	"github.com/labstack/echo/v4"
)

type Preview struct {
	Quantity   string  `msgpack:"quantity"`
	Unit       string  `msgpack:"unit"`
	Component  string  `msgpack:"component"`
	Layer      int     `msgpack:"layer"`
	MaxPoints  int     `msgpack:"maxPoints"`
	Dimensions [3]int  `msgpack:"dimensions"`
	Type       string  `msgpack:"type"`
	Buffer     *[]byte `msgpack:"buffer"`
	Min        float32 `msgpack:"min"`
	Max        float32 `msgpack:"max"`
}

func newPreview() *Preview {
	return &Preview{
		Quantity:   engine.NameOf(bps.Quantity),
		Unit:       engine.UnitOf(bps.Quantity),
		Component:  compIndexToString(bps.Component),
		Layer:      bps.Layer,
		MaxPoints:  bps.MaxPoints,
		Dimensions: bps.Dimensions,
		Type:       bps.Type,
		Buffer:     &bps.Buffer,
		Min:        bps.Min,
		Max:        bps.Max,
	}
}

func updateType() {
	isVectorField := bps.Quantity.NComp() == 3 && bps.Component == -1
	if isVectorField {
		bps.Type = "vector"
	} else {
		bps.Type = "scalar"
	}
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
	bps.Component = compStringToIndex(req.Component)
	updateType()
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
	bps.Quantity = quantity
	updateType()
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
	bps.Layer = req.Layer
	broadcastEngineState()
	return c.JSON(http.StatusOK, nil)
}

func postPreviewMaxPoints(c echo.Context) error {
	type Request struct {
		MaxPoints int `msgpack:"maxPoints"`
	}
	req := new(Request)
	if err := c.Bind(req); err != nil {
		util.LogErr(err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request payload"})
	}
	bps.MaxPoints = req.MaxPoints
	broadcastEngineState()
	return c.JSON(http.StatusOK, nil)
}

func postPreviewRefresh(c echo.Context) error {
	broadcastEngineState()
	return c.JSON(http.StatusOK, nil)
}
