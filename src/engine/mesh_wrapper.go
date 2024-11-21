package engine

import (
	"github.com/MathieuMoalic/amumax/src/cuda"
	"github.com/MathieuMoalic/amumax/src/mag"
)

// This package is a wrapper around data.Mesh, to allow for engine initialization during Mesh creation.

func CreateMesh() {
	setBusy(true)
	defer setBusy(false)
	Mesh.Create()
	normMag.alloc()
	Regions.alloc()
}

func ReCreateMesh(Nx, Ny, Nz int, dx, dy, dz float64, PBCx, PBCy, PBCz int) {
	setBusy(true)
	defer setBusy(false)
	Mesh.ReCreate(Nx, Ny, Nz, dx, dy, dz, PBCx, PBCy, PBCz)
	normMag.alloc()
	Regions.alloc()
	kernel := mag.DemagKernel(Mesh.Size(), Mesh.PBC(), Mesh.CellSize(), DemagAccuracy, CacheDir, ShowProgresBar)
	conv_ = cuda.NewDemag(Mesh.Size(), Mesh.PBC(), kernel, SelfTest)
}

func SmoothMesh(smoothx, smoothy, smoothz bool) {
	setBusy(true)
	defer setBusy(false)
	Mesh.SmoothMesh(smoothx, smoothy, smoothz)
	normMag.alloc()
	Regions.alloc()
}
