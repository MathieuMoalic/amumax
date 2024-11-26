package data

import (
	"math"

	"github.com/MathieuMoalic/amumax/src/log"
)

// MeshType stores info of a finite-difference mesh.
type MeshType struct {
	Nx, Ny, Nz                      int
	Dx, Dy, Dz                      float64
	Tx, Ty, Tz                      float64
	PBCx, PBCy, PBCz                int
	AutoMeshx, AutoMeshy, AutoMeshz bool
	created                         bool
}

// NewMesh creates a new mesh.
func NewMesh(nx, ny, nz int, dx, dy, dz float64, pbcx, pbcy, pbcz int) *MeshType {
	return &MeshType{nx, ny, nz, dx, dy, dz, 0, 0, 0, pbcx, pbcy, pbcz, false, false, false, false}
}

func (m *MeshType) Size() [3]int {
	return [3]int{m.Nx, m.Ny, m.Nz}
}

func (m *MeshType) CellSize() [3]float64 {
	return [3]float64{m.Dx, m.Dy, m.Dz}
}

// Returns pbc (periodic boundary conditions), as passed to constructor.
func (m *MeshType) PBC() [3]int {
	return [3]int{m.PBCx, m.PBCy, m.PBCz}
}

// Total number of cells, not taking into account PBCs.
func (m *MeshType) NCell() int {
	return m.Nx * m.Ny * m.Nz
}

// WorldSize equals (grid)Size x CellSize.
func (m *MeshType) WorldSize() [3]float64 {
	return [3]float64{float64(m.Nx) * m.Dx, float64(m.Ny) * m.Dy, float64(m.Nz) * m.Dz}
}

// 3 bools, packed in one byte, indicating whether there are periodic boundary conditions in
// X (LSB), Y(LSB<<1), Z(LSB<<2)
func (m *MeshType) PBC_code() byte {
	var code byte
	if m.PBCx != 0 {
		code = 1
	}
	if m.PBCy != 0 {
		code |= 2
	}
	if m.PBCz != 0 {
		code |= 4
	}
	return code
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

func (m *MeshType) smoothMesh() {
	if m.Nx*m.Ny*m.Nz < 10000 {
		log.Log.Info("No optimization to be made for small meshes")
		return
	}
	log.Log.Info("Original mesh: ")
	log.Log.Info("Cell size: %e, %e, %e", m.Dx, m.Dy, m.Dz)
	log.Log.Info("Grid Size: %d, %d, %d", m.Nx, m.Ny, m.Nz)
	if m.AutoMeshx {
		NewNx := closestSevenSmooth(m.Nx)
		m.Dx = m.Dx * float64(m.Nx) / float64(NewNx)
		m.Nx = NewNx
	}
	if m.AutoMeshy {
		NewNy := closestSevenSmooth(m.Ny)
		m.Dy = m.Dy * float64(m.Ny) / float64(NewNy)
		m.Ny = NewNy
	}
	if m.AutoMeshz {
		NewNz := closestSevenSmooth(m.Nz)
		m.Dz = m.Dz * float64(m.Nz) / float64(NewNz)
		m.Nz = NewNz
	}
	log.Log.Info("Smoothed mesh: ")
	log.Log.Info("Cell size: %e, %e, %e", m.Dx, m.Dy, m.Dz)
	log.Log.Info("Grid Size: %d, %d, %d", m.Nx, m.Ny, m.Nz)
}

func (m *MeshType) validateGridSize() {
	max_threshold := 1000000
	Ni_list := []string{"m.Nx", "m.Ny", "m.Nz"}
	for i, N := range []int{m.Nx, m.Ny, m.Nz} {
		if N == 0.0 {
			log.Log.ErrAndExit("Error: You have to specify  %v", Ni_list[i])
		} else if N > max_threshold {
			log.Log.ErrAndExit("Error: %s shouldn't be more than %d", Ni_list[i], max_threshold)
		} else if N < 0 {
			Ti := []float64{m.Tx, m.Ty, m.Tz}[i]
			di := []float64{m.Dx, m.Dy, m.Dz}[i]
			log.Log.ErrAndExit("Error: %s=%d shouldn't be negative, Ti: %e m, di: %e m", Ni_list[i], N, Ti, di)
		}
	}
	log.Log.Debug("Grid size: %d, %d, %d", m.Nx, m.Ny, m.Nz)
}

func (m *MeshType) checkLargestPrimeFactor(N int, AutoMesh bool, axisName string) {
	if largestPrimeFactor(N) > 127 && !AutoMesh {
		log.Log.ErrAndExit("Error: %s (%d) has a prime factor larger than 127 so the mesh cannot"+
			" be calculated. Use `AutoMesh%s = True` or change the value of %s manually", axisName, N, axisName, axisName)
	}
}

func (m *MeshType) validateCellSize() {
	min_threshold := 0.25e-9
	max_threshold := 500e-9
	names := []string{"dx", "dy", "dz"}
	for i, d := range []float64{m.Dx, m.Dy, m.Dz} {
		if d == 0.0 {
			log.Log.ErrAndExit("Error: You have to specify  %v", names[i])
		} else if d < min_threshold {
			log.Log.Warn("Warning: %s shouldn't be less than %f", names[i], min_threshold)
		} else if d > max_threshold {
			log.Log.Warn("Warning: %s shouldn't be more than %f", names[i], max_threshold)
		}
	}
	m.checkLargestPrimeFactor(m.Nx, m.AutoMeshx, "m.Nx")
	m.checkLargestPrimeFactor(m.Ny, m.AutoMeshy, "m.Ny")
	m.checkLargestPrimeFactor(m.Nz, m.AutoMeshz, "m.Nz")
	log.Log.Debug("Cell size: %e, %e, %e", m.Dx, m.Dy, m.Dz)
}

func (m *MeshType) isMeshCreated() bool {
	return m.created
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

func (m *MeshType) ReadyToCreate() bool {
	if m.created {
		return false
	} else if m.Dx != 0.0 && m.Dy != 0.0 && m.Dz != 0.0 && m.Nx != 0 && m.Ny != 0 && m.Nz != 0 {
		return true
	}
	return false
}

// check if mesh is set, otherwise, it creates it
func (m *MeshType) CreateMesh() {
	if !m.isMeshCreated() {
		setTiDiNi(&m.Tx, &m.Dx, &m.Nx, "x")
		setTiDiNi(&m.Ty, &m.Dy, &m.Ny, "y")
		setTiDiNi(&m.Tz, &m.Dz, &m.Nz, "z")
		m.validateGridSize()
		m.validateCellSize()
		if m.AutoMeshx || m.AutoMeshy || m.AutoMeshz {
			m.smoothMesh()
		}
	}
}

func ReCreateMesh(Nx_new, Ny_new, Nz_new int, dx_new, dy_new, dz_new float64, PBCx_new, PBCy_new, PBCz_new int) {
	// setBusy(true)
	// defer setBusy(false)
	// m.Nx = m.Nx_new
	// m.Ny = m.Ny_new
	// m.Nz = m.Nz_new
	// m.Dx = dx_new
	// m.Dy = dy_new
	// m.Dz = dz_new
	// PBCx = PBCx_new
	// PBCy = PBCy_new
	// PBCz = PBCz_new
	// m.Tx = 0.0
	// m.Ty = 0.0
	// m.Tz = 0.0
	// setTiDiNi(&m.Tx, &m.Dx, &m.Nx, "x")
	// setTiDiNi(&m.Ty, &m.Dy, &m.Ny, "y")
	// setTiDiNi(&m.Tz, &m.Dz, &m.Nz, "z")
	// validateGridSize()
	// validateCellSize()
	// if AutoMeshx || AutoMeshy || AutoMeshz {
	// 	smoothMesh()
	// }
	// globalmesh_ = *NewMesh(m.Nx, m.Ny, m.Nz, m.Dx, m.Dy, m.Dz, PBCx, PBCy, PBCz)
	// normMag.alloc()
	// Regions.alloc()
	// // these 2 lines make sure the progress bar doesn't break when calculating the kernel
	// fmt.Print("\033[2K\r") // clearline ANSI escape code
	// // kernel := mag.DemagKernel(GetMesh().Size(), GetMesh().PBC(), GetMesh().CellSize(), DemagAccuracy, CacheDir, ShowProgresBar)
	// // conv_ = cuda.NewDemag(GetMesh().Size(), GetMesh().PBC(), kernel, SelfTest)

}
