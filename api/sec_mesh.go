package api

import "github.com/MathieuMoalic/amumax/engine"

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
