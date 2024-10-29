package engine

import (
	"fmt"
	"math"

	"github.com/MathieuMoalic/amumax/src/cuda"
	"github.com/MathieuMoalic/amumax/src/data"
	"github.com/MathieuMoalic/amumax/src/log"
	"github.com/MathieuMoalic/amumax/src/mag"
	"github.com/MathieuMoalic/amumax/src/script"
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

func getMesh() *data.Mesh {
	createMesh()
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

func smoothMesh() {
	if Nx*Ny*Nz < 10000 {
		log.Log.Info("No optimization to be made for small meshes")
		return
	}
	log.Log.Info("Original mesh: ")
	log.Log.Info("Cell size: %e, %e, %e", Dx, Dy, Dz)
	log.Log.Info("Grid Size: %d, %d, %d", Nx, Ny, Nz)
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
	log.Log.Info("Smoothed mesh: ")
	log.Log.Info("Cell size: %e, %e, %e", Dx, Dy, Dz)
	log.Log.Info("Grid Size: %d, %d, %d", Nx, Ny, Nz)
}

func validateGridSize() {
	max_threshold := 1000000
	Ni_list := []string{"Nx", "Ny", "Nz"}
	for i, N := range []int{Nx, Ny, Nz} {
		if N == 0.0 {
			log.Log.ErrAndExit("Error: You have to specify  %v", Ni_list[i])
		} else if N > max_threshold {
			log.Log.ErrAndExit("Error: %s shouldn't be more than %d", Ni_list[i], max_threshold)
		} else if N < 0 {
			Ti := []float64{Tx, Ty, Tz}[i]
			di := []float64{Dx, Dy, Dz}[i]
			log.Log.ErrAndExit("Error: %s=%d shouldn't be negative, Ti: %e m, di: %e m", Ni_list[i], N, Ti, di)
		}
	}
	log.Log.Debug("Grid size: %d, %d, %d", Nx, Ny, Nz)
}

func checkLargestPrimeFactor(N int, AutoMesh bool, axisName string) {
	if largestPrimeFactor(N) > 127 && !AutoMesh {
		log.Log.ErrAndExit("Error: %s (%d) has a prime factor larger than 127 so the mesh cannot"+
			" be calculated. Use `AutoMesh%s = True` or change the value of %s manually", axisName, N, axisName, axisName)
	}
}

func validateCellSize() {
	min_threshold := 0.25e-9
	max_threshold := 500e-9
	names := []string{"dx", "dy", "dz"}
	for i, d := range []float64{Dx, Dy, Dz} {
		if d == 0.0 {
			log.Log.ErrAndExit("Error: You have to specify  %v", names[i])
		} else if d < min_threshold {
			log.Log.Warn("Warning: %s shouldn't be less than %f", names[i], min_threshold)
		} else if d > max_threshold {
			log.Log.Warn("Warning: %s shouldn't be more than %f", names[i], max_threshold)
		}
	}
	checkLargestPrimeFactor(Nx, AutoMeshx, "Nx")
	checkLargestPrimeFactor(Ny, AutoMeshy, "Ny")
	checkLargestPrimeFactor(Nz, AutoMeshz, "Nz")
	log.Log.Debug("Cell size: %e, %e, %e", Dx, Dy, Dz)
}

func isMeshCreated() bool {
	return globalmesh_.Size() != [3]int{0, 0, 0}
}

func setTiDiNi(Ti, di *float64, Ni *int, comp string) {
	if (*Ti != 0.0) && (*di != 0.0) && (*Ni != 0) {
		log.Log.ErrAndExit("Error: Only 2 of [N%s,d%s,T%s] are needed to define the mesh, you can't define all 3 of them.", comp, comp, comp)
	} else if (*Ti != 0.0) && (*di != 0.0) {
		*Ni = int(math.Round(*Ti / *di))
	} else if (*Ni != 0) && (*di != 0.0) {
		*Ti = *di * float64(*Ni)
	} else if (*Ni != 0) && (*Ti != 0.0) {
		*di = *Ti / float64(*Ni)
	}
}

// check if mesh is set, otherwise, it creates it
func createMesh() {
	if !isMeshCreated() {
		log.Log.Info("Creating mesh")
		setBusy(true)
		defer setBusy(false)
		setTiDiNi(&Tx, &Dx, &Nx, "x")
		setTiDiNi(&Ty, &Dy, &Ny, "y")
		setTiDiNi(&Tz, &Dz, &Nz, "z")
		validateGridSize()
		validateCellSize()
		if AutoMeshx || AutoMeshy || AutoMeshz {
			smoothMesh()
		}
		globalmesh_ = *data.NewMesh(Nx, Ny, Nz, Dx, Dy, Dz, PBCx, PBCy, PBCz)
		normMag.alloc()
		Regions.alloc()
		script.MMetadata.Init(OD(), StartTime, Dx, Dy, Dz, Nx, Ny, Nz, Tx, Ty, Tz, PBCx, PBCy, PBCz, cuda.GPUInfo)
	}
}

func reCreateMesh(Nx_new, Ny_new, Nz_new int, dx_new, dy_new, dz_new float64, PBCx_new, PBCy_new, PBCz_new int) {
	setBusy(true)
	defer setBusy(false)
	Nx = Nx_new
	Ny = Ny_new
	Nz = Nz_new
	Dx = dx_new
	Dy = dy_new
	Dz = dz_new
	PBCx = PBCx_new
	PBCy = PBCy_new
	PBCz = PBCz_new
	Tx = 0.0
	Ty = 0.0
	Tz = 0.0
	setTiDiNi(&Tx, &Dx, &Nx, "x")
	setTiDiNi(&Ty, &Dy, &Ny, "y")
	setTiDiNi(&Tz, &Dz, &Nz, "z")
	validateGridSize()
	validateCellSize()
	if AutoMeshx || AutoMeshy || AutoMeshz {
		smoothMesh()
	}
	globalmesh_ = *data.NewMesh(Nx, Ny, Nz, Dx, Dy, Dz, PBCx, PBCy, PBCz)
	normMag.alloc()
	Regions.alloc()
	// these 2 lines make sure the progress bar doesn't break when calculating the kernel
	fmt.Print("\033[2K\r") // clearline ANSI escape code
	kernel := mag.DemagKernel(getMesh().Size(), getMesh().PBC(), getMesh().CellSize(), DemagAccuracy, CacheDir, ShowProgresBar)
	conv_ = cuda.NewDemag(getMesh().Size(), getMesh().PBC(), kernel, SelfTest)

}
