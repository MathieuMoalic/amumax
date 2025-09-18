package engine

var (
	B_ext        = newExcitation("B_ext", "T", "Externally applied field")
	Edens_zeeman = newScalarField("Edens_Zeeman", "J/m3", "Zeeman energy density", AddEdens_zeeman)
	E_Zeeman     = newScalarValue("E_Zeeman", "J", "Zeeman energy", getZeemanEnergy)
)

var AddEdens_zeeman = makeEdensAdder(B_ext, -1)

func init() {
	registerEnergy(getZeemanEnergy, AddEdens_zeeman)
}

func getZeemanEnergy() float64 {
	return -1 * cellVolume() * dot(&MFull, B_ext)
}
