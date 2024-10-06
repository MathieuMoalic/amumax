package engine

import "github.com/MathieuMoalic/amumax/src/mag"

// Add a constant to the script world
func DeclConst(name string, value float64, doc string) {
	World.Const(name, value, doc)
}

func init() {
	DeclConst("Mu0", mag.Mu0, "Permittivity of vacuum (Tm/A)")
}
