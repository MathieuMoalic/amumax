package torque

import (
	"github.com/MathieuMoalic/amumax/src/cuda"
	"github.com/MathieuMoalic/amumax/src/log"
	"github.com/MathieuMoalic/amumax/src/magnetization"
	"github.com/MathieuMoalic/amumax/src/parameter"
	"github.com/MathieuMoalic/amumax/src/quantity"
	"github.com/MathieuMoalic/amumax/src/slice"
)

// TODO START
type VectorField struct {
	quantity.Quantity
}

func newVectorField(_, _, _ string, _ func(dst *slice.Slice)) *VectorField {
	return &VectorField{}
}

type Excitation struct {
	quantity.Quantity
}

func newExcitation(_, _, _ string) *Excitation {
	return &Excitation{}
}

type ScalarValue struct {
	quantity.Quantity
}

func newScalarValue(_, _, _ string, _ func() float64) *ScalarValue {
	return &ScalarValue{}
}

type fixedLayerPositionType int

// type flposition struct{}

// func (*flposition) Eval() interface{} {
// 	// return fixedLayerPosition
// 	return nil
// }
// func (*flposition) SetValue(v interface{}) {
// 	// drainOutput()
// 	// fixedLayerPosition = v.(fixedLayerPositionType)
// }
// func (*flposition) Type() reflect.Type { return reflect.TypeOf(fixedLayerPositionType(FIXEDLAYER_TOP)) }

const (
	FIXEDLAYER_TOP fixedLayerPositionType = iota + 1
	FIXEDLAYER_BOTTOM
)

// var (
// 	currentSignFromFixedLayerPosition = map[fixedLayerPositionType]float64{
// 		FIXEDLAYER_TOP:    1.0,
// 		FIXEDLAYER_BOTTOM: -1.0,
// 	}
// )

// TODO END

type Torque1 struct {
	mag                      *magnetization.Magnetization
	log                      *log.Logs
	Alpha                    *parameter.ScalarParam
	Xi                       *parameter.ScalarParam
	Pol                      *parameter.ScalarParam
	Lambda                   *parameter.ScalarParam
	EpsilonPrime             *parameter.ScalarParam
	FrozenSpins              *parameter.ScalarParam
	FreeLayerThickness       *parameter.ScalarParam
	FixedLayer               *Excitation
	Torque                   *VectorField
	LLTorque                 *VectorField
	STTorque                 *VectorField
	J                        *Excitation
	MaxTorque                *ScalarValue
	gammaLL                  float64
	precess                  bool
	disableZhangLiTorque     bool
	disableSlonczewskiTorque bool
	fixedLayerPosition       fixedLayerPositionType
}

func (t *Torque1) Init(log *log.Logs, mag *magnetization.Magnetization) {
	t.log = log
	t.mag = mag
	t.Alpha = parameter.NewScalarParam("alpha", "", "Landau-Lifshitz damping constant")
	t.Xi = parameter.NewScalarParam("xi", "", "Non-adiabaticity of spin-transfer-torque")
	t.Pol = parameter.NewScalarParam("Pol", "", "Electrical current polarization")
	t.Lambda = parameter.NewScalarParam("Lambda", "", "Slonczewski Λ parameter")
	t.EpsilonPrime = parameter.NewScalarParam("EpsilonPrime", "", "Slonczewski secondairy STT term ε'")
	t.FrozenSpins = parameter.NewScalarParam("frozenspins", "", "Defines spins that should be fixed") // 1 - frozen, 0 - free. TODO: check if it only contains 0/1 values
	t.FreeLayerThickness = parameter.NewScalarParam("FreeLayerThickness", "m", "Slonczewski free layer thickness")
	t.FixedLayer = newExcitation("FixedLayer", "", "Slonczewski fixed layer polarization")
	t.Torque = newVectorField("torque", "T", "Total torque/γ0", t.setTorque)
	t.LLTorque = newVectorField("LLtorque", "T", "Landau-Lifshitz torque/γ0", t.setLLTorque)
	// t.STTorque = newVectorField("STTorque", "T", "Spin-transfer torque/γ0", t.addSTTorque)
	t.J = newExcitation("J", "A/m2", "Electrical current density")
	t.MaxTorque = newScalarValue("maxTorque", "T", "Maximum torque/γ0, over all cells", t.getMaxTorque)
	t.gammaLL = 1.7595e11 // Gyromagnetic ratio of spins, in rad/Ts
	t.precess = true
	t.disableZhangLiTorque = false
	t.disableSlonczewskiTorque = false
	t.fixedLayerPosition = FIXEDLAYER_TOP // instructs mumax3 how free and fixed layers are stacked along +z direction

	t.Pol.SetUniform([]float64{1}) // default spin polarization
	t.Lambda.Set(1)                // sensible default value (?).
	// declROnly("FIXEDLAYER_TOP", FIXEDLAYER_TOP, "FixedLayerPosition = FIXEDLAYER_TOP instructs mumax3 that fixed layer is on top of the free layer")
	// declROnly("FIXEDLAYER_BOTTOM", FIXEDLAYER_BOTTOM, "FixedLayerPosition = FIXEDLAYER_BOTTOM instructs mumax3 that fixed layer is underneath of the free layer")
}

// Sets dst to the current total torque
func (t *Torque1) setTorque(dst *slice.Slice) {
	t.setLLTorque(dst)
	// t.addSTTorque(dst)
	// t.freezeSpins(dst)
}

// Sets dst to the current Landau-Lifshitz torque
func (t *Torque1) setLLTorque(dst *slice.Slice) {
	// setEffectiveField(dst) // calc and store B_eff
	alpha := t.Alpha.MSlice()
	defer alpha.Recycle()
	if t.precess {
		cuda.LLTorque(dst, t.mag.Slice, dst, alpha) // overwrite dst with torque
	} else {
		cuda.LLNoPrecess(dst, t.mag.Slice, dst)
	}
}

// // Adds the current spin transfer torque to dst
// func (t *Torque1) addSTTorque(dst *slice.Slice) {
// 	if t.J.isZero() {
// 		return
// 	}
// 	log.AssertMsg(!t.Pol.isZero(), "spin polarization should not be 0")
// 	jspin, rec := t.J.Slice()
// 	if rec {
// 		defer cuda.Recycle(jspin)
// 	}
// 	fl, rec := t.FixedLayer.Slice()
// 	if rec {
// 		defer cuda.Recycle(fl)
// 	}
// 	if !disableZhangLiTorque {
// 		msat := Msat.MSlice()
// 		defer msat.Recycle()
// 		j := J.MSlice()
// 		defer j.Recycle()
// 		alpha := Alpha.MSlice()
// 		defer alpha.Recycle()
// 		xi := Xi.MSlice()
// 		defer xi.Recycle()
// 		pol := Pol.MSlice()
// 		defer pol.Recycle()
// 		cuda.AddZhangLiTorque(dst, NormMag.Buffer(), msat, j, alpha, xi, pol, GetMesh())
// 	}
// 	if !disableSlonczewskiTorque && !FixedLayer.isZero() {
// 		msat := Msat.MSlice()
// 		defer msat.Recycle()
// 		j := J.MSlice()
// 		defer j.Recycle()
// 		fixedP := FixedLayer.MSlice()
// 		defer fixedP.Recycle()
// 		alpha := Alpha.MSlice()
// 		defer alpha.Recycle()
// 		pol := Pol.MSlice()
// 		defer pol.Recycle()
// 		lambda := Lambda.MSlice()
// 		defer lambda.Recycle()
// 		epsPrime := EpsilonPrime.MSlice()
// 		defer epsPrime.Recycle()
// 		thickness := FreeLayerThickness.MSlice()
// 		defer thickness.Recycle()
// 		cuda.AddSlonczewskiTorque2(dst, NormMag.Buffer(),
// 			msat, j, fixedP, alpha, pol, lambda, epsPrime,
// 			thickness,
// 			currentSignFromFixedLayerPosition[fixedLayerPosition],
// 			GetMesh())
// 	}
// }

// func (t *Torque1) freezeSpins(dst *slice.Slice) {
// 	if !t.FrozenSpins.isZero() {
// 		cuda.ZeroMask(dst, FrozenSpins.gpuLUT1(), Regions.Gpu())
// 	}
// }

func (t *Torque1) getMaxTorque() float64 {
	return cuda.MaxVecNorm(t.Torque.Value())
}
