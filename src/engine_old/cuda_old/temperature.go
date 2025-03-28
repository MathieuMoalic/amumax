package cuda_old

import (
	"github.com/MathieuMoalic/amumax/src/engine_old/data_old"
	"github.com/MathieuMoalic/amumax/src/engine_old/log_old"
)

// Set Bth to thermal noise (Brown).
// see temperature.cu
func SetTemperature(Bth, noise *data_old.Slice, k2mu0_Mu0VgammaDt float64, Msat, Temp, Alpha MSlice) {
	log_old.AssertMsg(Bth.NComp() == 1 && noise.NComp() == 1,
		"Component mismatch: Bth and noise must both have 1 component in SetTemperature")

	N := Bth.Len()
	cfg := make1DConf(N)

	k_settemperature2_async(Bth.DevPtr(0), noise.DevPtr(0), float32(k2mu0_Mu0VgammaDt),
		Msat.DevPtr(0), Msat.Mul(0),
		Temp.DevPtr(0), Temp.Mul(0),
		Alpha.DevPtr(0), Alpha.Mul(0),
		N, cfg)
}
