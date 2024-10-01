package data

// Mesh stores info of a finite-difference mesh.
type Mesh struct {
	Nx   int
	Ny   int
	Nz   int
	Dx   float64
	Dy   float64
	Dz   float64
	PBCx int
	PBCy int
	PBCz int
}

// NewMesh creates a new mesh.
func NewMesh(nx, ny, nz int, dx, dy, dz float64, pbcx, pbcy, pbcz int) *Mesh {
	return &Mesh{nx, ny, nz, dx, dy, dz, pbcx, pbcy, pbcz}
}

func (m *Mesh) Size() [3]int {
	return [3]int{m.Nx, m.Ny, m.Nz}
}

func (m *Mesh) CellSize() [3]float64 {
	return [3]float64{m.Dx, m.Dy, m.Dz}
}

// Returns pbc (periodic boundary conditions), as passed to constructor.
func (m *Mesh) PBC() [3]int {
	return [3]int{m.PBCx, m.PBCy, m.PBCz}
}

// Total number of cells, not taking into account PBCs.
func (m *Mesh) NCell() int {
	return m.Nx * m.Ny * m.Nz
}

// WorldSize equals (grid)Size x CellSize.
func (m *Mesh) WorldSize() [3]float64 {
	return [3]float64{float64(m.Nx) * m.Dx, float64(m.Ny) * m.Dy, float64(m.Nz) * m.Dz}
}

// 3 bools, packed in one byte, indicating whether there are periodic boundary conditions in
// X (LSB), Y(LSB<<1), Z(LSB<<2)
func (m *Mesh) PBC_code() byte {
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
