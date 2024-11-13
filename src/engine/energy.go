package engine

// Total energy calculation

import (
	"github.com/MathieuMoalic/amumax/src/cuda"
	"github.com/MathieuMoalic/amumax/src/data"
)

// TODO: Integrate(Edens)
// TODO: consistent naming SetEdensTotal, ...

var (
	energyTerms []func() float64        // all contributions to total energy
	edensTerms  []func(dst *data.Slice) // all contributions to total energy density (add to dst)
	Edens_total = newScalarField("Edens_total", "J/m3", "Total energy density", setTotalEdens)
	E_total     = newScalarValue("E_total", "J", "total energy", getTotalEnergy)
)

// add energy term to global energy
func registerEnergy(term func() float64, dens func(*data.Slice)) {
	energyTerms = append(energyTerms, term)
	edensTerms = append(edensTerms, dens)
}

// Returns the total energy in J.
func getTotalEnergy() float64 {
	E := 0.
	for _, f := range energyTerms {
		E += f()
	}
	checkNaN1(E)
	return E
}

// Set dst to total energy density in J/m3
func setTotalEdens(dst *data.Slice) {
	cuda.Zero(dst)
	for _, addTerm := range edensTerms {
		addTerm(dst)
	}
}

// volume of one cell in m3
func cellVolume() float64 {
	c := GetMesh().CellSize()
	return c[0] * c[1] * c[2]
}

// returns a function that adds to dst the energy density:
//
//	prefactor * dot (M_full, field)
func makeEdensAdder(field Quantity, prefactor float64) func(*data.Slice) {
	return func(dst *data.Slice) {
		B := ValueOf(field)
		defer cuda.Recycle(B)
		m := ValueOf(M_full)
		defer cuda.Recycle(m)
		factor := float32(prefactor)
		cuda.AddDotProduct(dst, factor, B, m)
	}
}

// vector dot product
func dot(a, b Quantity) float64 {
	A := ValueOf(a)
	defer cuda.Recycle(A)
	B := ValueOf(b)
	defer cuda.Recycle(B)
	return float64(cuda.Dot(A, B))
}
