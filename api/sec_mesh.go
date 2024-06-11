package api

import "github.com/MathieuMoalic/amumax/engine"

type Mesh struct {
	Dx   float64 `json:"dx"`
	Dy   float64 `json:"dy"`
	Dz   float64 `json:"dz"`
	Nx   int     `json:"Nx"`
	Ny   int     `json:"Ny"`
	Nz   int     `json:"Nz"`
	Tx   float64 `json:"Tx"`
	Ty   float64 `json:"Ty"`
	Tz   float64 `json:"Tz"`
	PBCx int     `json:"PBCx"`
	PBCy int     `json:"PBCy"`
	PBCz int     `json:"PBCz"`
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
