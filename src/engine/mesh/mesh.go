package mesh

import (
	"math"

	"github.com/MathieuMoalic/amumax/src/engine/log"
)

// Mesh stores info of a finite-difference mesh.
type Mesh struct {
	Nx, Ny, Nz       int
	Dx, Dy, Dz       float64
	Tx, Ty, Tz       float64
	PBCx, PBCy, PBCz int
	created          bool
}

// NewMesh creates an empty new mesh.
func NewMesh(nx, ny, nz int, dx, dy, dz float64, pbcx, pbcy, pbcz int) *Mesh {
	return &Mesh{nx, ny, nz, dx, dy, dz, 0, 0, 0, pbcx, pbcy, pbcz, false}
}

func (m Mesh) prettyPrint() {
	log.Log.Info("+----------------+------------+------------+------------+")
	log.Log.Info("| Axis           |     X      |     Y      |     Z      |")
	log.Log.Info("| Gridsize       | %10d | %10d | %10d |", m.Nx, m.Ny, m.Nz)
	log.Log.Info("| CellSize       | %10.3e | %10.3e | %10.3e |", m.Dx, m.Dy, m.Dz)
	log.Log.Info("| TotalSize      | %10.3e | %10.3e | %10.3e |", m.Tx, m.Ty, m.Tz)
	log.Log.Info("| PBC            | %10d | %10d | %10d |", m.PBCx, m.PBCy, m.PBCz)
	log.Log.Info("+----------------+------------+------------+------------+")
}

func (m *Mesh) Size() [3]int {
	return [3]int{m.Nx, m.Ny, m.Nz}
}
func (m *Mesh) GetNi() (int, int, int) {
	if !m.created {
		panic("Mesh not created yet")
	}
	return m.Nx, m.Ny, m.Nz
}

func (m *Mesh) CellSize() [3]float64 {
	if !m.created {
		panic("Mesh not created yet")
	}
	return [3]float64{m.Dx, m.Dy, m.Dz}
}
func (m *Mesh) GetDi() (float64, float64, float64) {
	if !m.created {
		panic("Mesh not created yet")
	}
	return m.Dx, m.Dy, m.Dz
}

// Returns pbc (periodic boundary conditions), as passed to constructor.
func (m *Mesh) PBC() [3]int {
	if !m.created {
		panic("Mesh not created yet")
	}
	return [3]int{m.PBCx, m.PBCy, m.PBCz}
}

// Total number of cells, not taking into account PBCs.
func (m *Mesh) NCell() int {
	if !m.created {
		panic("Mesh not created yet")
	}
	return m.Nx * m.Ny * m.Nz
}

// WorldSize equals (grid)Size x CellSize.
func (m *Mesh) WorldSize() [3]float64 {
	if !m.created {
		panic("Mesh not created yet")
	}
	return [3]float64{m.Tx, m.Ty, m.Tz}
}

// 3 bools, packed in one byte, indicating whether there are periodic boundary conditions in
// X (LSB), Y(LSB<<1), Z(LSB<<2)
func (m *Mesh) PBC_code() byte {
	if !m.created {
		panic("Mesh not created yet")
	}
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

func (m *Mesh) largestPrimeFactor(n int) int {
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

func (m *Mesh) closestSevenSmooth(n int) int {
	for m.largestPrimeFactor(n) > 7 {
		n -= 1
	}
	return n
}

func (m *Mesh) SmoothMesh(smoothx, smoothy, smoothz bool) {
	if !m.created {
		panic("Mesh not created yet")
	}
	if m.Nx*m.Ny*m.Nz < 10000 && m.Nx < 128 && m.Ny < 128 && m.Nz < 128 {
		log.Log.Info("No optimization to be made for small meshes")
		return
	}
	if !m.created {
		log.Log.ErrAndExit("Mesh not created yet")
	}
	if smoothx {
		NewNx := m.closestSevenSmooth(m.Nx)
		m.Dx = m.Dx * float64(m.Nx) / float64(NewNx)
		m.Nx = NewNx
		m.Tx = m.Dx * float64(m.Nx)
	}
	if smoothy {
		NewNy := m.closestSevenSmooth(m.Ny)
		m.Dy = m.Dy * float64(m.Ny) / float64(NewNy)
		m.Ny = NewNy
		m.Ty = m.Dy * float64(m.Ny)
	}
	if smoothz {
		NewNz := m.closestSevenSmooth(m.Nz)
		m.Dz = m.Dz * float64(m.Nz) / float64(NewNz)
		m.Nz = NewNz
		m.Tz = m.Dz * float64(m.Nz)
	}
	log.Log.Info("Smoothed mesh: ")
	m.prettyPrint()
}

func (m *Mesh) SetGridSize(Nx, Ny, Nz int) {
	m.Nx = Nx
	m.Ny = Ny
	m.Nz = Nz
}

func (m *Mesh) SetCellSize(Dx, Dy, Dz float64) {
	m.Dx = Dx
	m.Dy = Dy
	m.Dz = Dz
}

func (m *Mesh) SetTotalSize(Tx, Ty, Tz float64) {
	m.Tx = Tx
	m.Ty = Ty
	m.Tz = Tz
}

func (m *Mesh) SetPBC(PBCx, PBCy, PBCz int) {
	m.PBCx = PBCx
	m.PBCy = PBCy
	m.PBCz = PBCz
}

func (m *Mesh) SetMesh(Nx, Ny, Nz int, Dx, Dy, Dz float64, PBCx, PBCy, PBCz int) {
	m.SetGridSize(Nx, Ny, Nz)
	m.SetCellSize(Dx, Dy, Dz)
	m.SetPBC(PBCx, PBCy, PBCz)
}

func (m *Mesh) validateGridSize() {
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
}

func (m *Mesh) checkLargestPrimeFactor(N int, axisName string) {
	if m.largestPrimeFactor(N) > 127 {
		log.Log.Warn("%s (%d) has a prime factor larger than 127 so the mesh cannot", axisName, N)
		log.Log.Warn("be calculated. Use `AutoMesh(bool,bool,bool)` or change the value")
		log.Log.Warn("of %s manually or you might have CUDA errors.", axisName)
	}
}

func (m *Mesh) validateCellSize() {
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
	m.checkLargestPrimeFactor(m.Nx, "m.Nx")
	m.checkLargestPrimeFactor(m.Ny, "m.Ny")
	m.checkLargestPrimeFactor(m.Nz, "m.Nz")
}

func (m *Mesh) isAxisReadyToCreate(Ti, di float64, Ni int) bool {
	// if 2 of the 3 values are set, we return true
	if (Ti != 0.0 && di != 0.0) || (Ti != 0.0 && Ni != 0) || (di != 0.0 && Ni != 0) {
		return true
	}
	return false
}

func (m *Mesh) ReadyToCreate() bool {
	if m.created {
		return false
	} else if m.isAxisReadyToCreate(m.Tx, m.Dx, m.Nx) && m.isAxisReadyToCreate(m.Ty, m.Dy, m.Ny) && m.isAxisReadyToCreate(m.Tz, m.Dz, m.Nz) {
		return true
	}
	return false
}

func (m *Mesh) setTiDiNi(Ti, di *float64, Ni *int, comp string) {
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
func (m *Mesh) Create() {
	if !m.created {
		m.setTiDiNi(&m.Tx, &m.Dx, &m.Nx, "x")
		m.setTiDiNi(&m.Ty, &m.Dy, &m.Ny, "y")
		m.setTiDiNi(&m.Tz, &m.Dz, &m.Nz, "z")
		m.validateGridSize()
		m.validateCellSize()
		m.created = true
		m.prettyPrint()
	}
}

func (m *Mesh) ReCreate(Nx, Ny, Nz int, dx, dy, dz float64, PBCx, PBCy, PBCz int) {
	m.SetGridSize(Nx, Ny, Nz)
	m.SetCellSize(dx, dy, dz)
	m.SetPBC(PBCx, PBCy, PBCz)
	m.Tx = 0.0
	m.Ty = 0.0
	m.Tz = 0.0
	m.setTiDiNi(&m.Tx, &m.Dx, &m.Nx, "x")
	m.setTiDiNi(&m.Ty, &m.Dy, &m.Ny, "y")
	m.setTiDiNi(&m.Tz, &m.Dz, &m.Nz, "z")
	m.validateGridSize()
	m.validateCellSize()

}
