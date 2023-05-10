package engine

import (
	"fmt"
	"math"

	"github.com/MathieuMoalic/amumax/data"
	"github.com/MathieuMoalic/amumax/util"
)

var (
	Nx          int
	Ny          int
	Nz          int
	dx          float64
	dy          float64
	dz          float64
	Tx          float64
	Ty          float64
	Tz          float64
	PBCx        int
	PBCy        int
	PBCz        int
	AutoMesh    bool
	globalmesh_ data.Mesh
)

func init() {
	DeclVar("AutoMesh", &AutoMesh, "")
	DeclVar("Tx", &Tx, "")
	DeclVar("Ty", &Ty, "")
	DeclVar("Tz", &Tz, "")
	DeclVar("Nx", &Nx, "")
	DeclVar("Ny", &Ny, "")
	DeclVar("Nz", &Nz, "")
	DeclVar("dx", &dx, "")
	DeclVar("dy", &dy, "")
	DeclVar("dz", &dz, "")
	DeclVar("PBCx", &PBCx, "")
	DeclVar("PBCy", &PBCy, "")
	DeclVar("PBCz", &PBCz, "")
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
		LogOut("No optimization to be made for small meshes")
		return
	}
	LogOut("Original mesh: ")
	LogOut("Cell size: ", dx, dy, dz)
	LogOut("Grid Size: ", Nx, Ny, Nz)
	NewNx := closestSevenSmooth(Nx)
	NewNy := closestSevenSmooth(Ny)
	NewNz := closestSevenSmooth(Nz)
	dx := dx * float64(Nx) / float64(NewNx)
	dy := dy * float64(Ny) / float64(NewNy)
	dz := dz * float64(Nz) / float64(NewNz)
	Nx = NewNx
	Ny = NewNy
	Nz = NewNz
	LogOut("Smoothed mesh: ")
	LogOut("Cell size: ", dx, dy, dz)
	LogOut("Grid Size: ", Nx, Ny, Nz)
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
			util.Fatal("Error: You have to specify ", names[i])
		} else if N > max_threshold {
			util.Fatal("Error: ", names[i], " shouldn't be more than ", max_threshold)
		}
	}
}

func ValidateCellSize() {
	min_threshold := 0.25e-9
	max_threshold := 500e-9
	names := []string{"dx", "dy", "dz"}
	for i, d := range []float64{dx, dy, dz} {
		if d == 0.0 {
			util.Fatal("Error: You have to specify ", names[i])
		} else if d < min_threshold {
			util.Fatal("Error: ", names[i], "shouldn't be less than ", min_threshold)
		} else if d > max_threshold {
			util.Fatal("Error: ", names[i], "shouldn't be more than ", max_threshold)
		}
	}
}

func IsMeshCreated() bool {
	return globalmesh_.Size() == [3]int{0, 0, 0}
}

func SetTiDiNi(Ti, di *float64, Ni *int, comp string) {
	if (*Ti != 0.0) && (*di != 0.0) && (*Ni != 0) {
		util.Fatal(fmt.Sprintf("Error: Only 2 of [N%s,d%s,T%s] are needed to define the mesh, you can't define all 3 of them.", comp, comp, comp))
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
	if IsMeshCreated() {
		SetBusy(true)
		defer SetBusy(false)
		SetTiDiNi(&Tx, &dx, &Nx, "x")
		SetTiDiNi(&Ty, &dy, &Ny, "y")
		SetTiDiNi(&Tz, &dz, &Nz, "z")
		ValidateGridSize()
		ValidateCellSize()
		if AutoMesh {
			SmoothMesh()
		}
		globalmesh_ = *data.NewMesh(Nx, Ny, Nz, dx, dy, dz, PBCx, PBCy, PBCz)
		M.alloc()
		regions.alloc()
		InitMetadata()
	}
}
