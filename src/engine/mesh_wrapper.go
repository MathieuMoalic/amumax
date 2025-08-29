package engine

import (
	"github.com/MathieuMoalic/amumax/src/engine/cuda"
	"github.com/MathieuMoalic/amumax/src/engine/mag"
)

// This package is a wrapper around data.Mesh, to allow for engine initialization during Mesh creation.

func CreateMesh() {
	setBusy(true)
	defer setBusy(false)
	Mesh.Create()
	NormMag.Alloc()
	Regions.Alloc()
}

func SmoothMesh(smoothx, smoothy, smoothz bool) {
	setBusy(true)
	defer setBusy(false)
	Mesh.SmoothMesh(smoothx, smoothy, smoothz)
	NormMag.Alloc()
	Regions.Alloc()
	EngineState.Metadata.AddMesh(&Mesh)
}

// buggy and unused for now
func ReCreateMesh(Nx, Ny, Nz int, dx, dy, dz float64, PBCx, PBCy, PBCz int) {
	setBusy(true)
	defer setBusy(false)
	Mesh.ReCreate(Nx, Ny, Nz, dx, dy, dz, PBCx, PBCy, PBCz)
	NormMag.Alloc()
	Regions.Alloc()
	kernel := mag.DemagKernel(Mesh.Size(), Mesh.PBC(), Mesh.CellSize(), DemagAccuracy, CacheDir, HideProgresBar)
	conv_ = cuda.NewDemag(Mesh.Size(), Mesh.PBC(), kernel, SelfTest)
}
