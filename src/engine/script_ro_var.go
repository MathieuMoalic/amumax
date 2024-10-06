package engine

// Add a read-only variable to the script world.
// It can be changed, but not by the user.
func DeclROnly(name string, value interface{}, doc string) {
	World.ROnly(name, value, doc)
	AddQuantity(name, value, doc)
}

func init() {
	DeclROnly("regions", &Regions, "Outputs the region index for each cell")
}
