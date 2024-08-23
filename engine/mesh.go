package engine

import (
	"fmt"
	"math"

	"github.com/MathieuMoalic/amumax/cuda"
	"github.com/MathieuMoalic/amumax/data"
	"github.com/MathieuMoalic/amumax/mag"
	"github.com/MathieuMoalic/amumax/script"
	"github.com/MathieuMoalic/amumax/util"
)

var (
	Nx          int
	Ny          int
	Nz          int
	Dx          float64
	Dy          float64
	Dz          float64
	Tx          float64
	Ty          float64
	Tz          float64
	PBCx        int
	PBCy        int
	PBCz        int
	AutoMeshx   bool
	AutoMeshy   bool
	AutoMeshz   bool
	globalmesh_ data.Mesh
)

func init() {
	DeclVar("AutoMeshx", &AutoMeshx, "")
	DeclVar("AutoMeshy", &AutoMeshy, "")
	DeclVar("AutoMeshz", &AutoMeshz, "")
	DeclVar("Tx", &Tx, "")
	DeclVar("Ty", &Ty, "")
	DeclVar("Tz", &Tz, "")
	DeclVar("Nx", &Nx, "")
	DeclVar("Ny", &Ny, "")
	DeclVar("Nz", &Nz, "")
	DeclVar("dx", &Dx, "")
	DeclVar("dy", &Dy, "")
	DeclVar("dz", &Dz, "")
	DeclVar("PBCx", &PBCx, "")
	DeclVar("PBCy", &PBCy, "")
	DeclVar("PBCz", &PBCz, "")
	DeclFunc("ReCreateMesh", ReCreateMesh, "")
}

func GetMesh() *data.Mesh {
	CreateMesh()
	return &globalmesh_
}

func largestPrimeFactor(n int) int {
	maxPrime := -1
	for n%2 == 0 {
		maxPrime = 2
		n /= 2
	}
	for i := 3; i < int(math.Sqrt(float64(n))); i = i + 2 {
		for n%i == 0 {
			n /= i
		}
	}
	if n > 2 {
		maxPrime = n
	}
	return int(maxPrime)
}
func closestSevenSmooth(n int) int {
	for largestPrimeFactor(n) > 7 {
		n -= 1
	}
	return n
}

func SmoothMesh() {
	if Nx*Ny*Nz < 10000 {
		util.Log.Comment("No optimization to be made for small meshes")
		return
	}
	util.Log.Comment("Original mesh: ")
	util.Log.Comment("Cell size: %f, %f, %f", Dx, Dy, Dz)
	util.Log.Comment("Grid Size: %d, %d, %d", Nx, Ny, Nz)
	if AutoMeshx {
		NewNx := closestSevenSmooth(Nx)
		Dx = Dx * float64(Nx) / float64(NewNx)
		Nx = NewNx
	}
	if AutoMeshy {
		NewNy := closestSevenSmooth(Ny)
		Dy = Dy * float64(Ny) / float64(NewNy)
		Ny = NewNy
	}
	if AutoMeshz {
		NewNz := closestSevenSmooth(Nz)
		Dz = Dz * float64(Nz) / float64(NewNz)
		Nz = NewNz
	}
	util.Log.Comment("Smoothed mesh: ")
	util.Log.Comment("Cell size: %f, %f, %f", Dx, Dy, Dz)
	util.Log.Comment("Grid Size: %d, %d, %d", Nx, Ny, Nz)
}

func IsValidCellSize(cellSizeX, cellSizeY, cellSizeZ float64) bool {
	threshold := 2.5e-10
	if (cellSizeX < threshold) || (cellSizeY < threshold) || (cellSizeZ < threshold) {
		return false
	} else {
		return true
	}
}

func ValidateGridSize() {
	max_threshold := 1000000
	names := []string{"Nx", "Ny", "Nz"}
	for i, N := range []int{Nx, Ny, Nz} {
		if N == 0.0 {
			util.Log.ErrAndExit("Error: You have to specify  %v", names[i])
		} else if N > max_threshold {
			util.Log.ErrAndExit("Error: %s shouldn't be more than %d", names[i], max_threshold)
		}
	}
}

func ValidateCellSize() {
	min_threshold := 0.25e-9
	max_threshold := 500e-9
	names := []string{"dx", "dy", "dz"}
	for i, d := range []float64{Dx, Dy, Dz} {
		if d == 0.0 {
			util.Log.ErrAndExit("Error: You have to specify  %v", names[i])
		} else if d < min_threshold {
			util.Log.ErrAndExit("Error: %s shouldn't be less than %f", names[i], min_threshold)
		} else if d > max_threshold {
			util.Log.ErrAndExit("Error: %s shouldn't be more than %f", names[i], max_threshold)
		}
	}
}

func IsMeshCreated() bool {
	return globalmesh_.Size() != [3]int{0, 0, 0}
}

func SetTiDiNi(Ti, di *float64, Ni *int, comp string) {
	if (*Ti != 0.0) && (*di != 0.0) && (*Ni != 0) {
		util.Log.ErrAndExit(fmt.Sprintf("Error: Only 2 of [N%s,d%s,T%s] are needed to define the mesh, you can't define all 3 of them.", comp, comp, comp))
	} else if (*Ti != 0.0) && (*di != 0.0) {
		*Ni = int(math.Round(*Ti / *di))
	} else if (*Ni != 0) && (*di != 0.0) {
		*Ti = *di * float64(*Ni)
	} else if (*Ni != 0) && (*Ti != 0.0) {
		*di = *Ti / float64(*Ni)
	}
}

// check if mesh is set, otherwise, it creates it
func CreateMesh() {
	if !IsMeshCreated() {
		util.Log.Comment("Creating mesh")
		SetBusy(true)
		defer SetBusy(false)
		SetTiDiNi(&Tx, &Dx, &Nx, "x")
		SetTiDiNi(&Ty, &Dy, &Ny, "y")
		SetTiDiNi(&Tz, &Dz, &Nz, "z")
		ValidateGridSize()
		ValidateCellSize()
		if AutoMeshx || AutoMeshy || AutoMeshz {
			SmoothMesh()
		}
		globalmesh_ = *data.NewMesh(Nx, Ny, Nz, Dx, Dy, Dz, PBCx, PBCy, PBCz)
		M.alloc()
		Regions.alloc()
		script.MMetadata.Init(OD(), StartTime, Dx, Dy, Dz, Nx, Ny, Nz, Tx, Ty, Tz, PBCx, PBCy, PBCz, cuda.GPUInfo)
	}
}

func ReCreateMesh(Nx, Ny, Nz int, dx, dy, dz float64, PBCx, PBCy, PBCz int) {
	SetBusy(true)
	defer SetBusy(false)
	globalmesh_ = *data.NewMesh(Nx, Ny, Nz, dx, dy, dz, PBCx, PBCy, PBCz)
	M.alloc()
	Regions.alloc()
	script.MMetadata.Init(OD(), StartTime, dx, dy, dz, Nx, Ny, Nz, Tx, Ty, Tz, PBCx, PBCy, PBCz, cuda.GPUInfo)

	SetBusy(true)
	defer SetBusy(false)
	// these 2 lines make sure the progress bar doesn't break when calculating the kernel
	fmt.Print("\033[2K\r") // clearline ANSI escape code
	kernel := mag.DemagKernel(GetMesh().Size(), GetMesh().PBC(), GetMesh().CellSize(), DemagAccuracy, *Flag_cachedir, *Flag_magnets)
	conv_ = cuda.NewDemag(GetMesh().Size(), GetMesh().PBC(), kernel, *Flag_selftest)

}
