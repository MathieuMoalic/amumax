package engine

import "github.com/MathieuMoalic/amumax/src/mag"

// Add a constant to the script world
func declConst(name string, value float64, doc string) {
	World.Const(name, value, doc)
}

func init() {
	declConst("Mu0", mag.Mu0, "Permittivity of vacuum (Tm/A)")
}
