package engine

var (
	BExt        = newExcitation("B_ext", "T", "Externally applied field")
	EdensZeeman = newScalarField("Edens_Zeeman", "J/m3", "Zeeman energy density", AddEdensZeeman)
	EZeeman     = newScalarValue("E_Zeeman", "J", "Zeeman energy", getZeemanEnergy)
)

var AddEdensZeeman = makeEdensAdder(BExt, -1)

func init() {
	registerEnergy(getZeemanEnergy, AddEdensZeeman)
}

func getZeemanEnergy() float64 {
	return -1 * cellVolume() * dot(&MFull, BExt)
}
