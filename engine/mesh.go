package engine

import (
	"fmt"
	"math"

	"github.com/MathieuMoalic/amumax/cuda"
	"github.com/MathieuMoalic/amumax/data"
	"github.com/MathieuMoalic/amumax/util"
	"github.com/MathieuMoalic/amumax/zarr"
)

var globalmesh_ data.Mesh // mesh for m and everything that has the same size
var ZarrMeta zarr.MetaStruct
var AutoKernel bool

func init() {
	DeclFunc("SetGridSize", SetGridSize, `Sets the number of cells for X,Y,Z`)
	DeclFunc("SetCellSize", SetCellSize, `Sets the X,Y,Z cell size in meters`)
	DeclFunc("SetMesh", SetMesh, `Sets GridSize, CellSize and PBC at the same time`)
	DeclFunc("SetPBC", SetPBC, "Sets the number of repetitions in X,Y,Z to create periodic boundary "+
		"conditions. The number of repetitions determines the cutoff range for the demagnetization.")
	DeclVar("AutoKernel", &AutoKernel, "Smoothes the number of cells to make the FFT kernels faster, may change the number of cells")
}

func Mesh() *data.Mesh {
	checkMesh()
	return &globalmesh_
}

func arg(msg string, test bool) {
	if !test {
		panic(UserErr(msg + ": illegal arugment"))
	}
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

func SmoothKernel(Nx, Ny int, cellSizeX, cellSizeY float64) (int, int, float64, float64) {
	NewNx := closestSevenSmooth(Nx)
	NewNy := closestSevenSmooth(Ny)
	NewCellSizeX := cellSizeX * float64(Nx) / float64(NewNx)
	NewCellSizeY := cellSizeY * float64(Ny) / float64(NewNy)
	return NewNx, NewNy, NewCellSizeX, NewCellSizeY
}

func IsValidCellSize(cellSizeX, cellSizeY, cellSizeZ float64) bool {
	threshold := 2.5e-10
	if (cellSizeX < threshold) || (cellSizeY < threshold) || (cellSizeZ < threshold) {
		return false
	} else {
		return true
	}
}

// Set the simulation mesh to Nx x Ny x Nz cells of given size.
// Can be set only once at the beginning of the simulation.
// TODO: dedup arguments from globals
func SetMesh(Nx, Ny, Nz int, cellSizeX, cellSizeY, cellSizeZ float64, pbcx, pbcy, pbcz int) {
	SetBusy(true)
	defer SetBusy(false)

	arg("GridSize", Nx > 0 && Ny > 0 && Nz > 0)
	arg("CellSize", cellSizeX > 0 && cellSizeY > 0 && cellSizeZ > 0)
	arg("PBC", pbcx >= 0 && pbcy >= 0 && pbcz >= 0)

	if !IsValidCellSize(cellSizeX, cellSizeY, cellSizeZ) {
		util.Println("/!\\/!\\/!\\/!\\ Warning: Cell sizes might be too small. /!\\/!\\/!\\/!\\")
	}
	if Nx*Ny > 10000 {
		NewNx, NewNy, NewcellSizeX, NewcellSizeY := SmoothKernel(Nx, Ny, cellSizeX, cellSizeY)
		if AutoKernel {
			Nx, Ny, cellSizeX, cellSizeY = NewNx, NewNy, NewcellSizeX, NewcellSizeY
			fmt.Println("Kernel smoothed to: ", Nx, Ny, cellSizeX, cellSizeY)
		} else if (Nx != NewNx) || (Ny != NewNy) {
			util.Println("/!\\/!\\/!\\/!\\ Warning: Your gridsize doesn't appear to be optimized, expect longer simulations /!\\/!\\/!\\/!\\")

		}
	}

	prevSize := globalmesh_.Size()
	pbc := []int{pbcx, pbcy, pbcz}

	if globalmesh_.Size() == [3]int{0, 0, 0} {
		// first time mesh is set
		globalmesh_ = *data.NewMesh(Nx, Ny, Nz, cellSizeX, cellSizeY, cellSizeZ, pbc...)
		M.alloc()
		regions.alloc()
	} else {
		// here be dragons
		LogOut("resizing...")

		// free everything
		conv_.Free()
		conv_ = nil
		mfmconv_.Free()
		mfmconv_ = nil
		cuda.FreeBuffers()

		// resize everything
		globalmesh_ = *data.NewMesh(Nx, Ny, Nz, cellSizeX, cellSizeY, cellSizeZ, pbc...)
		M.resize()
		regions.resize()
		geometry.buffer.Free()
		geometry.buffer = data.NilSlice(1, Mesh().Size())
		geometry.setGeom(geometry.shape)

		// remove excitation extra terms if they don't fit anymore
		// up to the user to add them again
		if Mesh().Size() != prevSize {
			B_ext.RemoveExtraTerms()
			J.RemoveExtraTerms()
		}

		if Mesh().Size() != prevSize {
			B_therm.noise.Free()
			B_therm.noise = nil
		}
	}
	lazy_gridsize = []int{Nx, Ny, Nz}
	lazy_cellsize = []float64{cellSizeX, cellSizeY, cellSizeZ}
	lazy_pbc = []int{pbcx, pbcy, pbcz}
	ZarrMeta.Init(globalmesh_, OD(), cuda.GPUInfo)
	if chunks.x.nb == 0 {
		chunks = Chunks{
			Chunk{globalmesh_.Size()[X], 1},
			Chunk{globalmesh_.Size()[Y], 1},
			Chunk{globalmesh_.Size()[Z], 1},
			Chunk{3, 1},
		}
	}
}

func printf(f float64) float32 {
	return float32(f)
}

// for lazy setmesh: set gridsize and cellsize in separate calls
var (
	lazy_gridsize []int
	lazy_cellsize []float64
	lazy_pbc      = []int{0, 0, 0}
)

func SetGridSize(Nx, Ny, Nz int) {
	lazy_gridsize = []int{Nx, Ny, Nz}
	if lazy_cellsize != nil {
		SetMesh(Nx, Ny, Nz, lazy_cellsize[X], lazy_cellsize[Y], lazy_cellsize[Z], lazy_pbc[X], lazy_pbc[Y], lazy_pbc[Z])
	}
}

func SetCellSize(cx, cy, cz float64) {
	lazy_cellsize = []float64{cx, cy, cz}
	if lazy_gridsize != nil {
		SetMesh(lazy_gridsize[X], lazy_gridsize[Y], lazy_gridsize[Z], cx, cy, cz, lazy_pbc[X], lazy_pbc[Y], lazy_pbc[Z])
	}
}

func SetPBC(nx, ny, nz int) {
	lazy_pbc = []int{nx, ny, nz}
	if lazy_gridsize != nil && lazy_cellsize != nil {
		SetMesh(lazy_gridsize[X], lazy_gridsize[Y], lazy_gridsize[Z],
			lazy_cellsize[X], lazy_cellsize[Y], lazy_cellsize[Z],
			lazy_pbc[X], lazy_pbc[Y], lazy_pbc[Z])
	}
}

// check if mesh is set
func checkMesh() {
	if globalmesh_.Size() == [3]int{0, 0, 0} {
		panic("need to set mesh first")
	}
}
