package engine

// Add an LValue to the script world.
// Assign to LValue invokes SetValue()
func declLValue(name string, value lValue, doc string) {
	addParameter(name, value, doc)
	World.LValue(name, newLValueWrapper(name, value), doc)
	addQuantity(name, value, doc)
}

func init() {
	declLValue("MFMLift", &mfmLift, "MFM lift height")
	declLValue("MFMDipole", &mfmTipSize, "Height of vertically magnetized part of MFM tip")
	declLValue("m", &NormMag, `Reduced magnetization (unit length)`)

	declLValue("FilenameFormat", &fformat{}, "printf formatting string for output filenames.")
	declLValue("OutputFormat", &oformat{}, "Format for data files: OVF1_TEXT, OVF1_BINARY, OVF2_TEXT or OVF2_BINARY")

	declLValue("FixedLayerPosition", &flposition{}, "Position of the fixed layer: FIXEDLAYER_TOP, FIXEDLAYER_BOTTOM (default=FIXEDLAYER_TOP)")
}
