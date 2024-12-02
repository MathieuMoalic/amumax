package engine_old

import "github.com/MathieuMoalic/amumax/src/engine_old/mag_old"

// Add a constant to the script world
func declConst(name string, value float64, doc string) {
	World.Const(name, value, doc)
}

func init() {
	declConst("Mu0", mag_old.Mu0, "Permittivity of vacuum (Tm/A)")
}
