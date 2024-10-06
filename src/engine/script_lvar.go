package engine

// Add an LValue to the script world.
// Assign to LValue invokes SetValue()
func DeclLValue(name string, value LValue, doc string) {
	AddParameter(name, value, doc)
	World.LValue(name, newLValueWrapper(name, value), doc)
	AddQuantity(name, value, doc)
}

func init() {
	DeclLValue("MFMLift", &MFMLift, "MFM lift height")
	DeclLValue("MFMDipole", &MFMTipSize, "Height of vertically magnetized part of MFM tip")
	DeclLValue("m", &M, `Reduced magnetization (unit length)`)

	DeclLValue("FilenameFormat", &fformat{}, "printf formatting string for output filenames.")
	DeclLValue("OutputFormat", &oformat{}, "Format for data files: OVF1_TEXT, OVF1_BINARY, OVF2_TEXT or OVF2_BINARY")

	DeclLValue("FixedLayerPosition", &flposition{}, "Position of the fixed layer: FIXEDLAYER_TOP, FIXEDLAYER_BOTTOM (default=FIXEDLAYER_TOP)")
}
