package api

import (
	"net/http"

	"github.com/MathieuMoalic/amumax/src/engine_old"
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

func initMeshAPI(e *echo.Group, ws *WebSocketManager) *MeshState {
	meshState := MeshState{
		ws:   ws,
		Dx:   &engine_old.Mesh.Dx,
		Dy:   &engine_old.Mesh.Dy,
		Dz:   &engine_old.Mesh.Dz,
		Nx:   &engine_old.Mesh.Nx,
		Ny:   &engine_old.Mesh.Ny,
		Nz:   &engine_old.Mesh.Nz,
		Tx:   &engine_old.Mesh.Tx,
		Ty:   &engine_old.Mesh.Ty,
		Tz:   &engine_old.Mesh.Tz,
		PBCx: &engine_old.Mesh.PBCx,
		PBCy: &engine_old.Mesh.PBCy,
		PBCz: &engine_old.Mesh.PBCz,
	}
	e.POST("/api/mesh", meshState.postMesh)
	return &meshState
}

func (m *MeshState) Update() {}

func (m *MeshState) postMesh(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, "")
}
