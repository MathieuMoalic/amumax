package engine_old

import (
	"reflect"

	"github.com/MathieuMoalic/amumax/src/engine_old/cuda_old"
	"github.com/MathieuMoalic/amumax/src/engine_old/data_old"
	"github.com/MathieuMoalic/amumax/src/engine_old/log_old"
)

var (
	Alpha                            = newScalarParam("alpha", "", "Landau-Lifshitz damping constant")
	Xi                               = newScalarParam("xi", "", "Non-adiabaticity of spin-transfer-torque")
	Pol                              = newScalarParam("Pol", "", "Electrical current polarization")
	Lambda                           = newScalarParam("Lambda", "", "Slonczewski Λ parameter")
	EpsilonPrime                     = newScalarParam("EpsilonPrime", "", "Slonczewski secondairy STT term ε'")
	FrozenSpins                      = newScalarParam("frozenspins", "", "Defines spins that should be fixed") // 1 - frozen, 0 - free. TODO: check if it only contains 0/1 values
	FreeLayerThickness               = newScalarParam("FreeLayerThickness", "m", "Slonczewski free layer thickness")
	FixedLayer                       = newExcitation("FixedLayer", "", "Slonczewski fixed layer polarization")
	Torque                           = newVectorField("torque", "T", "Total torque/γ0", setTorque)
	LLTorque                         = newVectorField("LLtorque", "T", "Landau-Lifshitz torque/γ0", setLLTorque)
	STTorque                         = newVectorField("STTorque", "T", "Spin-transfer torque/γ0", addSTTorque)
	J                                = newExcitation("J", "A/m2", "Electrical current density")
	MaxTorque                        = newScalarValue("maxTorque", "T", "Maximum torque/γ0, over all cells", getMaxTorque)
	gammaLL                  float64 = 1.7595e11 // Gyromagnetic ratio of spins, in rad/Ts
	precess                          = true
	disableZhangLiTorque             = false
	disableSlonczewskiTorque         = false
	fixedLayerPosition               = FIXEDLAYER_TOP // instructs mumax3 how free and fixed layers are stacked along +z direction
)

func init() {
	Pol.setUniform([]float64{1}) // default spin polarization
	Lambda.Set(1)                // sensible default value (?).
	declROnly("FIXEDLAYER_TOP", FIXEDLAYER_TOP, "FixedLayerPosition = FIXEDLAYER_TOP instructs mumax3 that fixed layer is on top of the free layer")
	declROnly("FIXEDLAYER_BOTTOM", FIXEDLAYER_BOTTOM, "FixedLayerPosition = FIXEDLAYER_BOTTOM instructs mumax3 that fixed layer is underneath of the free layer")
}

// Sets dst to the current total torque
func setTorque(dst *data_old.Slice) {
	setLLTorque(dst)
	addSTTorque(dst)
	freezeSpins(dst)
}

// Sets dst to the current Landau-Lifshitz torque
func setLLTorque(dst *data_old.Slice) {
	setEffectiveField(dst) // calc and store B_eff
	alpha := Alpha.MSlice()
	defer alpha.Recycle()
	if precess {
		cuda_old.LLTorque(dst, NormMag.Buffer(), dst, alpha) // overwrite dst with torque
	} else {
		cuda_old.LLNoPrecess(dst, NormMag.Buffer(), dst)
	}
}

// Adds the current spin transfer torque to dst
func addSTTorque(dst *data_old.Slice) {
	if J.isZero() {
		return
	}
	log_old.AssertMsg(!Pol.isZero(), "spin polarization should not be 0")
	jspin, rec := J.Slice()
	if rec {
		defer cuda_old.Recycle(jspin)
	}
	fl, rec := FixedLayer.Slice()
	if rec {
		defer cuda_old.Recycle(fl)
	}
	if !disableZhangLiTorque {
		msat := Msat.MSlice()
		defer msat.Recycle()
		j := J.MSlice()
		defer j.Recycle()
		alpha := Alpha.MSlice()
		defer alpha.Recycle()
		xi := Xi.MSlice()
		defer xi.Recycle()
		pol := Pol.MSlice()
		defer pol.Recycle()
		cuda_old.AddZhangLiTorque(dst, NormMag.Buffer(), msat, j, alpha, xi, pol, GetMesh())
	}
	if !disableSlonczewskiTorque && !FixedLayer.isZero() {
		msat := Msat.MSlice()
		defer msat.Recycle()
		j := J.MSlice()
		defer j.Recycle()
		fixedP := FixedLayer.MSlice()
		defer fixedP.Recycle()
		alpha := Alpha.MSlice()
		defer alpha.Recycle()
		pol := Pol.MSlice()
		defer pol.Recycle()
		lambda := Lambda.MSlice()
		defer lambda.Recycle()
		epsPrime := EpsilonPrime.MSlice()
		defer epsPrime.Recycle()
		thickness := FreeLayerThickness.MSlice()
		defer thickness.Recycle()
		cuda_old.AddSlonczewskiTorque2(dst, NormMag.Buffer(),
			msat, j, fixedP, alpha, pol, lambda, epsPrime,
			thickness,
			currentSignFromFixedLayerPosition[fixedLayerPosition],
			GetMesh())
	}
}

func freezeSpins(dst *data_old.Slice) {
	if !FrozenSpins.isZero() {
		cuda_old.ZeroMask(dst, FrozenSpins.gpuLUT1(), Regions.Gpu())
	}
}

func getMaxTorque() float64 {
	torque := ValueOf(Torque)
	defer cuda_old.Recycle(torque)
	return cuda_old.MaxVecNorm(torque)
}

type fixedLayerPositionType int

const (
	FIXEDLAYER_TOP fixedLayerPositionType = iota + 1
	FIXEDLAYER_BOTTOM
)

var (
	currentSignFromFixedLayerPosition = map[fixedLayerPositionType]float64{
		FIXEDLAYER_TOP:    1.0,
		FIXEDLAYER_BOTTOM: -1.0,
	}
)

type flposition struct{}

func (*flposition) Eval() interface{} { return fixedLayerPosition }
func (*flposition) SetValue(v interface{}) {
	drainOutput()
	fixedLayerPosition = v.(fixedLayerPositionType)
}
func (*flposition) Type() reflect.Type { return reflect.TypeOf(fixedLayerPositionType(FIXEDLAYER_TOP)) }
