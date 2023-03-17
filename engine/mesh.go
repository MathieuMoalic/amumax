package engine

import (
	"math"

	"github.com/MathieuMoalic/amumax/cuda"
	"github.com/MathieuMoalic/amumax/data"
	"github.com/MathieuMoalic/amumax/util"
	"github.com/MathieuMoalic/amumax/zarr"
)

var globalmesh_ data.Mesh // mesh for m and everything that has the same size
var ZarrMeta zarr.MetaStruct

var chunks Chunks
var Nx int
var Ny int
var Nz int
var dx float64
var dy float64
var dz float64
var PBCx int
var PBCy int
var PBCz int

type Chunk struct {
	len int
	nb  int
}

type Chunks struct {
	x Chunk
	y Chunk
	z Chunk
	c Chunk
}

func init() {
	DeclFunc("InitMesh", InitMesh, "")
	DeclFunc("InitAutoMesh", InitAutoMesh, "")
	DeclFunc("Chunkx", Chunkx, "")
	DeclFunc("Chunky", Chunky, "")
	DeclFunc("Chunkz", Chunkz, "")
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
	checkMesh()
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
		LogOut("No optmization to be made for small meshes")
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
			util.Fatal("Error: ", names[i], "shouldn't be more than ", max_threshold)
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

func InitAutoMesh() {
	InitMeshInner(true)
}

func InitMesh() {
	InitMeshInner(false)
}
func InitMeshInner(smooth bool) {
	SetBusy(true)
	defer SetBusy(false)
	ValidateGridSize()
	ValidateCellSize()
	if smooth {
		SmoothMesh()
	}
	if globalmesh_.Size() == [3]int{0, 0, 0} {
		globalmesh_ = *data.NewMesh(Nx, Ny, Nz, dx, dy, dz, PBCx, PBCy, PBCz)
		M.alloc()
		regions.alloc()
	} else {
		panic("Can't resize the mesh in amumax sorry :/")
	}
	chunks.Init()
	ZarrMeta.Init(globalmesh_, OD(), cuda.GPUInfo)
}

func Chunkx(nb_of_chunks int) {
	chunks.x.Assign(0, nb_of_chunks)
}
func Chunky(nb_of_chunks int) {
	chunks.y.Assign(1, nb_of_chunks)
}
func Chunkz(nb_of_chunks int) {
	chunks.z.Assign(2, nb_of_chunks)
}

func (chunks *Chunks) Init() {
	*chunks = Chunks{
		Chunk{Nx, 1},
		Chunk{Ny, 1},
		Chunk{Nz, 1},
		Chunk{3, 1},
	}
}

func (chunk *Chunk) Assign(N_index, nb_of_chunks int) {
	name := []string{"Nx", "Ny", "Nz"}[N_index]
	N := []int{Nx, Ny, Nz}[N_index]
	checkMesh()
	if chunk.nb != 1 {
		util.Fatal("Error: You cannot change the chunking once it's set.")
	}
	if nb_of_chunks < 1 || (nb_of_chunks > N) {
		util.Fatal("Error: The number of chunks must be between 1 and ", name)
	}
	new_nb_of_chunks := closestDivisor(N, nb_of_chunks)
	if new_nb_of_chunks != nb_of_chunks {
		LogOut("Warning: The number of chunks for", name, "has been automatically resized from", nb_of_chunks, "to", new_nb_of_chunks)
	}
	nb_of_chunks = new_nb_of_chunks
	*chunk = Chunk{N / nb_of_chunks, nb_of_chunks}
}

func printf(f float64) float32 {
	return float32(f)
}

func closestDivisor(N int, D int) int {
	closest := 0
	minDist := math.MaxInt32
	for i := 1; i <= N; i++ {
		if N%i == 0 {
			dist := i - D
			if dist < 0 {
				dist = -dist
			}
			if dist < minDist {
				minDist = dist
				closest = i
			}
		}
	}
	return closest
}

// check if mesh is set
func checkMesh() {
	if globalmesh_.Size() == [3]int{0, 0, 0} {
		util.Fatal("Error: You need to set the mesh first")
	}
}
