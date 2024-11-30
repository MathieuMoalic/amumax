package engine_old

import (
	"github.com/MathieuMoalic/amumax/src/cuda"
	"github.com/MathieuMoalic/amumax/src/cuda/curand"
	"github.com/MathieuMoalic/amumax/src/data"
	"github.com/MathieuMoalic/amumax/src/engine_old/log_old"
	"github.com/MathieuMoalic/amumax/src/engine_old/mesh_old"
	"github.com/MathieuMoalic/amumax/src/mag"

	"math"
)

var (
	Temp                      = newScalarParam("Temp", "K", "Temperature")
	E_therm                   = newScalarValue("E_therm", "J", "Thermal energy", getThermalEnergy)
	Edens_therm               = newScalarField("Edens_therm", "J/m3", "Thermal energy density", AddThermalEnergyDensity)
	B_therm                   thermField // Thermal effective field (T)
	printedWarningTempOddGrid = false
)

var AddThermalEnergyDensity = makeEdensAdder(&B_therm, -1)

// thermField calculates and caches thermal noise.
type thermField struct {
	seed      int64            // seed for generator
	generator curand.Generator //
	noise     *data.Slice      // noise buffer
	step      int              // solver step corresponding to noise
	dt        float64          // solver timestep corresponding to noise
}

func init() {
	registerEnergy(getThermalEnergy, AddThermalEnergyDensity)
	B_therm.step = -1 // invalidate noise cache
	declROnly("B_therm", &B_therm, "Thermal field (T)")
}

func (b *thermField) AddTo(dst *data.Slice) {
	if !Temp.isZero() {
		b.update()
		cuda.Add(dst, dst, b.noise)
	}
}

func (b *thermField) update() {
	// we need to fix the time step here because solver will not yet have done it before the first step.
	// FixDt as an lvalue that sets Dt_si on change might be cleaner.
	if FixDt != 0 {
		Dt_si = FixDt
	}

	if b.generator == 0 {
		b.generator = curand.CreateGenerator(curand.PSEUDO_DEFAULT)
		b.generator.SetSeed(b.seed)
	}
	if b.noise == nil {
		b.noise = cuda.NewSlice(b.NComp(), b.Mesh().Size())
		// when noise was (re-)allocated it's invalid for sure.
		B_therm.step = -1
		B_therm.dt = -1
	}

	if Temp.isZero() {
		cuda.Memset(b.noise, 0, 0, 0)
		b.step = NSteps
		b.dt = Dt_si
		return
	}

	// keep constant during time step
	if NSteps == b.step && Dt_si == b.dt {
		return
	}

	// after a bad step the timestep is rescaled and the noise should be rescaled accordingly, instead of redrawing the random numbers
	if NSteps == b.step && Dt_si != b.dt {
		for c := 0; c < 3; c++ {
			cuda.Madd2(b.noise.Comp(c), b.noise.Comp(c), b.noise.Comp(c), float32(math.Sqrt(b.dt/Dt_si)), 0.)
		}
		b.dt = Dt_si
		return
	}

	N := GetMesh().NCell()

	if !printedWarningTempOddGrid && N%2 > 0 { // T is nonzero if we have gotten this far. As noted in issue #314, this means the grid size must be even.
		printedWarningTempOddGrid = true
		log_old.Log.Warn("nonzero temperature requires an even amount of grid cells, but all axes have "+
			"an odd number of cells: %v. This may cause a CURAND_STATUS_LENGTH_NOT_MULTIPLE error.", GetMesh().Size())
	}
	k2_VgammaDt := 2 * mag.Kb / (gammaLL * cellVolume() * Dt_si)
	noise := cuda.Buffer(1, GetMesh().Size())
	defer cuda.Recycle(noise)

	const mean = 0
	const stddev = 1
	dst := b.noise
	ms := Msat.MSlice()
	defer ms.Recycle()
	temp := Temp.MSlice()
	defer temp.Recycle()
	alpha := Alpha.MSlice()
	defer alpha.Recycle()
	for i := 0; i < 3; i++ {
		b.generator.GenerateNormal(uintptr(noise.DevPtr(0)), int64(N), mean, stddev)
		cuda.SetTemperature(dst.Comp(i), noise, k2_VgammaDt, ms, temp, alpha)
	}

	b.step = NSteps
	b.dt = Dt_si
}

func getThermalEnergy() float64 {
	if Temp.isZero() || relaxing {
		return 0
	} else {
		return -cellVolume() * dot(&M_full, &B_therm)
	}
}

// Seeds the thermal noise generator
func thermSeed(seed int) {
	B_therm.seed = int64(seed)
	if B_therm.generator != 0 {
		B_therm.generator.SetSeed(B_therm.seed)
	}
}

func (b *thermField) Mesh() *mesh_old.Mesh   { return GetMesh() }
func (b *thermField) NComp() int             { return 3 }
func (b *thermField) Name() string           { return "Thermal field" }
func (b *thermField) Unit() string           { return "T" }
func (b *thermField) average() []float64     { return qAverageUniverse(b) }
func (b *thermField) EvalTo(dst *data.Slice) { evalTo(b, dst) }
func (b *thermField) Slice() (*data.Slice, bool) {
	b.update()
	return b.noise, false
}
