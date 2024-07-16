package api

import (
	"net/http"
	"strconv"

	"github.com/MathieuMoalic/amumax/engine"
	"github.com/labstack/echo/v4"
)

type Mesh struct {
	Dx   float64 `msgpack:"dx"`
	Dy   float64 `msgpack:"dy"`
	Dz   float64 `msgpack:"dz"`
	Nx   int     `msgpack:"Nx"`
	Ny   int     `msgpack:"Ny"`
	Nz   int     `msgpack:"Nz"`
	Tx   float64 `msgpack:"Tx"`
	Ty   float64 `msgpack:"Ty"`
	Tz   float64 `msgpack:"Tz"`
	PBCx int     `msgpack:"PBCx"`
	PBCy int     `msgpack:"PBCy"`
	PBCz int     `msgpack:"PBCz"`
}

func newMesh() *Mesh {
	return &Mesh{
		Dx:   engine.Dx,
		Dy:   engine.Dy,
		Dz:   engine.Dz,
		Nx:   engine.Nx,
		Ny:   engine.Ny,
		Nz:   engine.Nz,
		Tx:   engine.Tx,
		Ty:   engine.Ty,
		Tz:   engine.Tz,
		PBCx: engine.PBCx,
		PBCy: engine.PBCy,
		PBCz: engine.PBCz,
	}
}

func postMesh(c echo.Context) error {
	type Request struct {
		Runtime float64 `msgpack:"runtime"`
	}

	req := new(Request)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request payload"})
	}
	engine.Break()
	engine.Inject <- func() { EvalGUI("Run(" + strconv.FormatFloat(req.Runtime, 'f', -1, 64) + ")") }
	return c.JSON(http.StatusOK, "")
}
