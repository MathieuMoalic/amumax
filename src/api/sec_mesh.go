package api

import (
	"net/http"

	"github.com/MathieuMoalic/amumax/src/engine"
	"github.com/labstack/echo/v4"
)

type MeshState struct {
	ws   *WebSocketManager
	Dx   *float64 `msgpack:"dx"`
	Dy   *float64 `msgpack:"dy"`
	Dz   *float64 `msgpack:"dz"`
	Nx   *int     `msgpack:"Nx"`
	Ny   *int     `msgpack:"Ny"`
	Nz   *int     `msgpack:"Nz"`
	Tx   *float64 `msgpack:"Tx"`
	Ty   *float64 `msgpack:"Ty"`
	Tz   *float64 `msgpack:"Tz"`
	PBCx *int     `msgpack:"PBCx"`
	PBCy *int     `msgpack:"PBCy"`
	PBCz *int     `msgpack:"PBCz"`
}

func initMeshAPI(e *echo.Echo, ws *WebSocketManager) *MeshState {
	meshState := MeshState{
		ws:   ws,
		Dx:   &engine.Dx,
		Dy:   &engine.Dy,
		Dz:   &engine.Dz,
		Nx:   &engine.Nx,
		Ny:   &engine.Ny,
		Nz:   &engine.Nz,
		Tx:   &engine.Tx,
		Ty:   &engine.Ty,
		Tz:   &engine.Tz,
		PBCx: &engine.PBCx,
		PBCy: &engine.PBCy,
		PBCz: &engine.PBCz,
	}
	e.POST("/api/mesh", meshState.postMesh)
	return &meshState
}

func (m *MeshState) Update() {}

func (m *MeshState) postMesh(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, "")
}
