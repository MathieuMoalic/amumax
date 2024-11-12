package engine

// Add a (pointer to) variable to the script world
func declVar(name string, value interface{}, doc string) {
	World.Var(name, value, doc)
	addQuantity(name, value, doc)
}

// Hack for fixing the closure caveat:
// Defines "t", the time variable, handled specially by Fix()
func declTVar(name string, value interface{}, doc string) {
	World.TVar(name, value, doc)
	addQuantity(name, value, doc)
}

func init() {
	declTVar("t", &Time, "Total simulated time (s)")

	declVar("EnableDemag", &EnableDemag, "Enables/disables demag (default=true)")
	declVar("DemagAccuracy", &DemagAccuracy, "Controls accuracy of demag kernel")

	declVar("step", &NSteps, "Total number of time steps taken")
	declVar("MinDt", &MinDt, "Minimum time step the solver can take (s)")
	declVar("MaxDt", &MaxDt, "Maximum time step the solver can take (s)")
	declVar("MaxErr", &MaxErr, "Maximum error per step the solver can tolerate (default = 1e-5)")
	declVar("Headroom", &Headroom, "Solver headroom (default = 0.8)")
	declVar("FixDt", &FixDt, "Set a fixed time step, 0 disables fixed step (which is the default)")
	declVar("OpenBC", &OpenBC, "Use open boundary conditions (default=false)")
	declVar("ext_BubbleMz", &BubbleMz, "Center magnetization 1.0 or -1.0  (default = 1.0)")
	declVar("EdgeSmooth", &edgeSmooth, "Geometry edge smoothing with edgeSmooth^3 samples per cell, 0=staircase, ~8=very smooth")

	declVar("AutoMeshx", &AutoMeshx, "")
	declVar("AutoMeshy", &AutoMeshy, "")
	declVar("AutoMeshz", &AutoMeshz, "")
	declVar("Tx", &Tx, "")
	declVar("Ty", &Ty, "")
	declVar("Tz", &Tz, "")
	declVar("Nx", &Nx, "")
	declVar("Ny", &Ny, "")
	declVar("Nz", &Nz, "")
	declVar("dx", &Dx, "")
	declVar("dy", &Dy, "")
	declVar("dz", &Dz, "")
	declVar("PBCx", &PBCx, "")
	declVar("PBCy", &PBCy, "")
	declVar("PBCz", &PBCz, "")
	declVar("MinimizerStop", &stopMaxDm, "Stopping max dM for Minimize")
	declVar("MinimizerSamples", &dmSamples, "Number of max dM to collect for Minimize convergence check.")
	declVar("MinimizeMaxSteps", &minimizeMaxSteps, "")
	declVar("MinimizeMaxTimeSeconds", &minimizeMaxTimeSeconds, "")
	declVar("RelaxTorqueThreshold", &relaxTorqueThreshold, "MaxTorque threshold for relax(). If set to -1 (default), relax() will stop when the average torque is steady or increasing.")
	declVar("SnapshotFormat", &snapshotFormat, "Image format for snapshots: jpg, png or gif.")

	declVar("ShiftMagL", &shiftMagL, "Upon shift, insert this magnetization from the left")
	declVar("ShiftMagR", &shiftMagR, "Upon shift, insert this magnetization from the right")
	declVar("ShiftMagU", &shiftMagU, "Upon shift, insert this magnetization from the top")
	declVar("ShiftMagD", &shiftMagD, "Upon shift, insert this magnetization from the bottom")
	declVar("ShiftM", &shiftM, "Whether Shift() acts on magnetization")
	declVar("ShiftGeom", &shiftGeom, "Whether Shift() acts on geometry")
	declVar("ShiftRegions", &shiftRegions, "Whether Shift() acts on regions")
	declVar("TotalShift", &totalShift, "Amount by which the simulation has been shifted (m).")
	declVar("EdgeCarryShift", &EdgeCarryShift, "Whether to use the current magnetization at the border for the cells inserted by Shift")

	declVar("GammaLL", &gammaLL, "Gyromagnetic ratio in rad/Ts")
	declVar("DisableZhangLiTorque", &disableZhangLiTorque, "Disables Zhang-Li torque (default=false)")
	declVar("DisableSlonczewskiTorque", &disableSlonczewskiTorque, "Disables Slonczewski torque (default=false)")
	declVar("DoPrecess", &precess, "Enables LL precession (default=true)")

	declVar("PreviewXDataPoints", &PreviewXDataPoints, "Number of data points in the x direction for the 2D/3D preview")
	declVar("PreviewYDataPoints", &PreviewYDataPoints, "Number of data points in the y direction for the 2D/3D preview")
}

var PreviewXDataPoints = 100
var PreviewYDataPoints = 100
