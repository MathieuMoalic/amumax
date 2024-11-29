package cuda

// for compatibility between mesh_old and mesh
type MeshLike interface {
	Size() [3]int
	PBC() [3]int
	CellSize() [3]float64
	PBC_code() byte
	WorldSize() [3]float64
}
