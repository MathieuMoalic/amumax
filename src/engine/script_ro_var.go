package engine

// Add a read-only variable to the script world.
// It can be changed, but not by the user.
func declROnly(name string, value any, doc string) {
	World.ROnly(name, value, doc)
	addQuantity(name, value, doc)
}

func init() {
	declROnly("regions", &Regions, "Outputs the region index for each cell")
}
