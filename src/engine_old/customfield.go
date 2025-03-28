package engine_old

// Add arbitrary terms to B_eff, Edens_total.

import (
	"fmt"

	"github.com/MathieuMoalic/amumax/src/engine_old/cuda_old"
	"github.com/MathieuMoalic/amumax/src/engine_old/data_old"
	"github.com/MathieuMoalic/amumax/src/engine_old/log_old"
	"github.com/MathieuMoalic/amumax/src/engine_old/mesh_old"
)

var (
	B_custom       = newVectorField("B_custom", "T", "User-defined field", addCustomField)
	Edens_custom   = newScalarField("Edens_custom", "J/m3", "Energy density of user-defined field.", addCustomEnergyDensity)
	E_custom       = newScalarValue("E_custom", "J", "total energy of user-defined field", getCustomEnergy)
	customTerms    []Quantity // vector
	customEnergies []Quantity // scalar
)

func init() {
	registerEnergy(getCustomEnergy, addCustomEnergyDensity)
}

// Removes all customfields
func removeCustomFields() {
	customTerms = nil
}

// addFieldTerm adds an effective field function (returning Teslas) to B_eff.
// Be sure to also add the corresponding energy term using AddEnergyTerm.
func addFieldTerm(b Quantity) {
	customTerms = append(customTerms, b)
}

// AddEnergyTerm adds an energy density function (returning Joules/m³) to Edens_total.
// Needed when AddFieldTerm was used and a correct energy is needed
// (e.g. for Relax, Minimize, ...).
func addEdensTerm(e Quantity) {
	customEnergies = append(customEnergies, e)
}

// addCustomField evaluates the user-defined custom field terms
// and adds the result to dst.
func addCustomField(dst *data_old.Slice) {
	for _, term := range customTerms {
		buf := ValueOf(term)
		cuda_old.Add(dst, dst, buf)
		cuda_old.Recycle(buf)
	}
}

// Adds the custom energy densities (defined with AddCustomE
func addCustomEnergyDensity(dst *data_old.Slice) {
	for _, term := range customEnergies {
		buf := ValueOf(term)
		cuda_old.Add(dst, dst, buf)
		cuda_old.Recycle(buf)
	}
}

func getCustomEnergy() float64 {
	buf := cuda_old.Buffer(1, GetMesh().Size())
	defer cuda_old.Recycle(buf)
	cuda_old.Zero(buf)
	addCustomEnergyDensity(buf)
	return cellVolume() * float64(cuda_old.Sum(buf))
}

type constValue struct {
	value []float64
}

func (c *constValue) NComp() int { return len(c.value) }

func (d *constValue) EvalTo(dst *data_old.Slice) {
	for c, v := range d.value {
		cuda_old.Memset(dst.Comp(c), float32(v))
	}
}

// constScalar returns a constant (uniform) scalar quantity,
// that can be used to construct custom field terms.
func constScalar(v float64) Quantity {
	return &constValue{[]float64{v}}
}

// constVector returns a constant (uniform) vector quantity,
// that can be used to construct custom field terms.
func constVector(x, y, z float64) Quantity {
	return &constValue{[]float64{x, y, z}}
}

// fieldOp holds the abstract functionality for operations
// (like add, multiply, ...) on space-dependend quantites
// (like M, B_sat, ...)
type fieldOp struct {
	a, b  Quantity
	nComp int
}

func (o fieldOp) NComp() int {
	return o.nComp
}

type dotProduct struct {
	fieldOp
}

type crossProduct struct {
	fieldOp
}

type addition struct {
	fieldOp
}

type mAddition struct {
	fieldOp
	fac1, fac2 float64
}

type mulmv struct {
	ax, ay, az, b Quantity
}

// mulMV returns a new Quantity that evaluates to the
// matrix-vector product (Ax·b, Ay·b, Az·b).
func mulMV(Ax, Ay, Az, b Quantity) Quantity {
	log_old.AssertMsg(Ax.NComp() == 3 &&
		Ay.NComp() == 3 &&
		Az.NComp() == 3 &&
		b.NComp() == 3,
		"Component mismatch: Ax, Ay, Az, and b must all have 3 components in mulMV")
	return &mulmv{Ax, Ay, Az, b}
}

func (q *mulmv) EvalTo(dst *data_old.Slice) {
	log_old.AssertMsg(dst.NComp() == 3, "Component mismatch: dst must have 3 components in EvalTo")
	cuda_old.Zero(dst)
	b := ValueOf(q.b)
	defer cuda_old.Recycle(b)

	{
		Ax := ValueOf(q.ax)
		cuda_old.AddDotProduct(dst.Comp(X), 1, Ax, b)
		cuda_old.Recycle(Ax)
	}
	{
		Ay := ValueOf(q.ay)
		cuda_old.AddDotProduct(dst.Comp(Y), 1, Ay, b)
		cuda_old.Recycle(Ay)
	}
	{
		Az := ValueOf(q.az)
		cuda_old.AddDotProduct(dst.Comp(Z), 1, Az, b)
		cuda_old.Recycle(Az)
	}
}

func (q *mulmv) NComp() int {
	return 3
}

// DotProduct creates a new quantity that is the dot product of
// quantities a and b. E.g.:
//
//	DotProct(&M, &B_ext)
func dotProductFunc(a, b Quantity) Quantity {
	return &dotProduct{fieldOp{a, b, 1}}
}

func (d *dotProduct) EvalTo(dst *data_old.Slice) {
	A := ValueOf(d.a)
	defer cuda_old.Recycle(A)
	B := ValueOf(d.b)
	defer cuda_old.Recycle(B)
	cuda_old.Zero(dst)
	cuda_old.AddDotProduct(dst, 1, A, B)
}

// CrossProduct creates a new quantity that is the cross product of
// quantities a and b. E.g.:
//
//	CrossProct(&M, &B_ext)
func cross(a, b Quantity) Quantity {
	return &crossProduct{fieldOp{a, b, 3}}
}

func (d *crossProduct) EvalTo(dst *data_old.Slice) {
	A := ValueOf(d.a)
	defer cuda_old.Recycle(A)
	B := ValueOf(d.b)
	defer cuda_old.Recycle(B)
	cuda_old.Zero(dst)
	cuda_old.CrossProduct(dst, A, B)
}

func add(a, b Quantity) Quantity {
	if a.NComp() != b.NComp() {
		panic(fmt.Sprintf("Cannot point-wise Add %v components by %v components", a.NComp(), b.NComp()))
	}
	return &addition{fieldOp{a, b, a.NComp()}}
}

func (d *addition) EvalTo(dst *data_old.Slice) {
	A := ValueOf(d.a)
	defer cuda_old.Recycle(A)
	B := ValueOf(d.b)
	defer cuda_old.Recycle(B)
	cuda_old.Zero(dst)
	cuda_old.Add(dst, A, B)
}

type pointwiseMul struct {
	fieldOp
}

func madd(a, b Quantity, fac1, fac2 float64) *mAddition {
	if a.NComp() != b.NComp() {
		panic(fmt.Sprintf("Cannot point-wise add %v components by %v components", a.NComp(), b.NComp()))
	}
	return &mAddition{fieldOp{a, b, a.NComp()}, fac1, fac2}
}

func (o *mAddition) EvalTo(dst *data_old.Slice) {
	A := ValueOf(o.a)
	defer cuda_old.Recycle(A)
	B := ValueOf(o.b)
	defer cuda_old.Recycle(B)
	cuda_old.Zero(dst)
	cuda_old.Madd2(dst, A, B, float32(o.fac1), float32(o.fac2))
}

// mul returns a new quantity that evaluates to the pointwise product a and b.
func mul(a, b Quantity) Quantity {
	nComp := -1
	switch {
	case a.NComp() == b.NComp():
		nComp = a.NComp() // vector*vector, scalar*scalar
	case a.NComp() == 1:
		nComp = b.NComp() // scalar*something
	case b.NComp() == 1:
		nComp = a.NComp() // something*scalar
	default:
		panic(fmt.Sprintf("Cannot point-wise multiply %v components by %v components", a.NComp(), b.NComp()))
	}

	return &pointwiseMul{fieldOp{a, b, nComp}}
}

func (d *pointwiseMul) EvalTo(dst *data_old.Slice) {
	cuda_old.Zero(dst)
	a := ValueOf(d.a)
	defer cuda_old.Recycle(a)
	b := ValueOf(d.b)
	defer cuda_old.Recycle(b)

	switch {
	case a.NComp() == b.NComp():
		mulNN(dst, a, b) // vector*vector, scalar*scalar
	case a.NComp() == 1:
		mul1N(dst, a, b)
	case b.NComp() == 1:
		mul1N(dst, b, a)
	default:
		panic(fmt.Sprintf("Cannot point-wise multiply %v components by %v components", a.NComp(), b.NComp()))
	}
}

// mulNN pointwise multiplies two N-component vectors,
// yielding an N-component vector stored in dst.
func mulNN(dst, a, b *data_old.Slice) {
	cuda_old.Mul(dst, a, b)
}

// mul1N pointwise multiplies a scalar (1-component) with an N-component vector,
// yielding an N-component vector stored in dst.
func mul1N(dst, a, b *data_old.Slice) {
	log_old.AssertMsg(a.NComp() == 1, "Component mismatch: a must have 1 component in mul1N")
	log_old.AssertMsg(dst.NComp() == b.NComp(), "Component mismatch: dst and b must have the same number of components in mul1N")
	for c := 0; c < dst.NComp(); c++ {
		cuda_old.Mul(dst.Comp(c), a, b.Comp(c))
	}
}

type pointwiseDiv struct {
	fieldOp
}

// div returns a new quantity that evaluates to the pointwise product a and b.
func div(a, b Quantity) Quantity {
	nComp := -1
	switch {
	case a.NComp() == b.NComp():
		nComp = a.NComp() // vector/vector, scalar/scalar
	case b.NComp() == 1:
		nComp = a.NComp() // something/scalar
	default:
		panic(fmt.Sprintf("Cannot point-wise divide %v components by %v components", a.NComp(), b.NComp()))
	}
	return &pointwiseDiv{fieldOp{a, b, nComp}}
}

func (d *pointwiseDiv) EvalTo(dst *data_old.Slice) {
	a := ValueOf(d.a)
	defer cuda_old.Recycle(a)
	b := ValueOf(d.b)
	defer cuda_old.Recycle(b)

	switch {
	case a.NComp() == b.NComp():
		divNN(dst, a, b) // vector*vector, scalar*scalar
	case b.NComp() == 1:
		divN1(dst, a, b)
	default:
		panic(fmt.Sprintf("Cannot point-wise divide %v components by %v components", a.NComp(), b.NComp()))
	}

}

func divNN(dst, a, b *data_old.Slice) {
	cuda_old.Div(dst, a, b)
}

func divN1(dst, a, b *data_old.Slice) {
	log_old.AssertMsg(dst.NComp() == a.NComp(), "Component mismatch: dst and a must have the same number of components in divN1")
	log_old.AssertMsg(b.NComp() == 1, "Component mismatch: b must have 1 component in divN1")
	for c := 0; c < dst.NComp(); c++ {
		cuda_old.Div(dst.Comp(c), a.Comp(c), b)
	}
}

type shifted struct {
	orig       Quantity
	dx, dy, dz int
}

// shiftedQuant returns a new Quantity that evaluates to
// the original, shifted over dx, dy, dz cells.
func shiftedQuant(q Quantity, dx, dy, dz int) Quantity {
	log_old.AssertMsg(dx != 0 || dy != 0 || dz != 0, "Invalid shift: at least one of dx, dy, or dz must be non-zero in shiftedQuant")
	return &shifted{q, dx, dy, dz}
}

func (q *shifted) EvalTo(dst *data_old.Slice) {
	orig := ValueOf(q.orig)
	defer cuda_old.Recycle(orig)
	for i := 0; i < q.NComp(); i++ {
		dsti := dst.Comp(i)
		origi := orig.Comp(i)
		if q.dx != 0 {
			cuda_old.ShiftX(dsti, origi, q.dx, 0, 0)
		}
		if q.dy != 0 {
			cuda_old.ShiftY(dsti, origi, q.dy, 0, 0)
		}
		if q.dz != 0 {
			cuda_old.ShiftZ(dsti, origi, q.dz, 0, 0)
		}
	}
}

func (q *shifted) NComp() int {
	return q.orig.NComp()
}

// Masks a quantity with a shape
// The shape will be only evaluated once on the mesh,
// and will be re-evaluated after mesh change,
// because otherwise too slow
func maskedQuant(q Quantity, shape shape) Quantity {
	return &masked{q, shape, nil, mesh_old.Mesh{}}
}

type masked struct {
	orig  Quantity
	shape shape
	mask  *data_old.Slice
	mesh  mesh_old.Mesh
}

func (q *masked) EvalTo(dst *data_old.Slice) {
	if q.mesh != *GetMesh() {
		// When mesh is changed, mask needs an update
		q.createMask()
	}
	orig := ValueOf(q.orig)
	defer cuda_old.Recycle(orig)
	mul1N(dst, q.mask, orig)
}

func (q *masked) NComp() int {
	return q.orig.NComp()
}

func (q *masked) createMask() {
	size := GetMesh().Size()
	// Prepare mask on host
	maskhost := data_old.NewSlice(SCALAR, size)
	defer maskhost.Free()
	maskScalars := maskhost.Scalars()
	for iz := 0; iz < size[Z]; iz++ {
		for iy := 0; iy < size[Y]; iy++ {
			for ix := 0; ix < size[X]; ix++ {
				r := index2Coord(ix, iy, iz)
				if q.shape(r[X], r[Y], r[Z]) {
					maskScalars[iz][iy][ix] = 1
				}
			}
		}
	}
	// Update mask
	q.mask.Free()
	q.mask = cuda_old.NewSlice(SCALAR, size)
	data_old.Copy(q.mask, maskhost)
	q.mesh = *GetMesh()
	// Remove mask from host
}

// normalizedQuant returns a quantity that evaluates to the unit vector of q
func normalizedQuant(q Quantity) Quantity {
	return &normalized{q}
}

type normalized struct {
	orig Quantity
}

func (q *normalized) NComp() int {
	return 3
}

func (q *normalized) EvalTo(dst *data_old.Slice) {
	log_old.AssertMsg(dst.NComp() == q.NComp(), "Component mismatch: dst must have the same number of components as the normalized quantity in EvalTo")
	q.orig.EvalTo(dst)
	cuda_old.Normalize(dst, nil)
}
