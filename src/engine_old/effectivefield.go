package engine_old

// Effective field

import "github.com/MathieuMoalic/amumax/src/engine_old/data_old"

var B_eff = newVectorField("B_eff", "T", "Effective field", setEffectiveField)

// Sets dst to the current effective field, in Tesla.
// This is the sum of all effective field terms,
// like demag, exchange, ...
func setEffectiveField(dst *data_old.Slice) {
	setDemagField(dst)    // set to B_demag...
	addExchangeField(dst) // ...then add other terms
	addAnisotropyField(dst)
	addMagnetoelasticField(dst)
	B_ext.AddTo(dst)
	if !relaxing {
		B_therm.AddTo(dst)
	}
	addCustomField(dst)
}
