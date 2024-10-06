package engine

// Add a (pointer to) variable to the script world
func DeclVar(name string, value interface{}, doc string) {
	World.Var(name, value, doc)
	AddQuantity(name, value, doc)
}

// Hack for fixing the closure caveat:
// Defines "t", the time variable, handled specially by Fix()
func DeclTVar(name string, value interface{}, doc string) {
	World.TVar(name, value, doc)
	AddQuantity(name, value, doc)
}
func init() {
	DeclTVar("t", &Time, "Total simulated time (s)")

	DeclVar("EnableDemag", &EnableDemag, "Enables/disables demag (default=true)")
	DeclVar("DemagAccuracy", &DemagAccuracy, "Controls accuracy of demag kernel")

	DeclVar("step", &NSteps, "Total number of time steps taken")
	DeclVar("MinDt", &MinDt, "Minimum time step the solver can take (s)")
	DeclVar("MaxDt", &MaxDt, "Maximum time step the solver can take (s)")
	DeclVar("MaxErr", &MaxErr, "Maximum error per step the solver can tolerate (default = 1e-5)")
	DeclVar("Headroom", &Headroom, "Solver headroom (default = 0.8)")
	DeclVar("FixDt", &FixDt, "Set a fixed time step, 0 disables fixed step (which is the default)")
	DeclVar("OpenBC", &OpenBC, "Use open boundary conditions (default=false)")
	DeclVar("ext_BubbleMz", &BubbleMz, "Center magnetization 1.0 or -1.0  (default = 1.0)")
	DeclVar("EdgeSmooth", &edgeSmooth, "Geometry edge smoothing with edgeSmooth^3 samples per cell, 0=staircase, ~8=very smooth")

	DeclVar("AutoMeshx", &AutoMeshx, "")
	DeclVar("AutoMeshy", &AutoMeshy, "")
	DeclVar("AutoMeshz", &AutoMeshz, "")
	DeclVar("Tx", &Tx, "")
	DeclVar("Ty", &Ty, "")
	DeclVar("Tz", &Tz, "")
	DeclVar("Nx", &Nx, "")
	DeclVar("Ny", &Ny, "")
	DeclVar("Nz", &Nz, "")
	DeclVar("dx", &Dx, "")
	DeclVar("dy", &Dy, "")
	DeclVar("dz", &Dz, "")
	DeclVar("PBCx", &PBCx, "")
	DeclVar("PBCy", &PBCy, "")
	DeclVar("PBCz", &PBCz, "")
	DeclVar("MinimizerStop", &StopMaxDm, "Stopping max dM for Minimize")
	DeclVar("MinimizerSamples", &DmSamples, "Number of max dM to collect for Minimize convergence check.")
	DeclVar("MinimizeMaxSteps", &MinimizeMaxSteps, "")
	DeclVar("MinimizeMaxTimeSeconds", &MinimizeMaxTimeSeconds, "")
	DeclVar("RelaxTorqueThreshold", &RelaxTorqueThreshold, "MaxTorque threshold for relax(). If set to -1 (default), relax() will stop when the average torque is steady or increasing.")
	DeclVar("SnapshotFormat", &SnapshotFormat, "Image format for snapshots: jpg, png or gif.")

	DeclVar("ShiftMagL", &ShiftMagL, "Upon shift, insert this magnetization from the left")
	DeclVar("ShiftMagR", &ShiftMagR, "Upon shift, insert this magnetization from the right")
	DeclVar("ShiftMagU", &ShiftMagU, "Upon shift, insert this magnetization from the top")
	DeclVar("ShiftMagD", &ShiftMagD, "Upon shift, insert this magnetization from the bottom")
	DeclVar("ShiftM", &ShiftM, "Whether Shift() acts on magnetization")
	DeclVar("ShiftGeom", &ShiftGeom, "Whether Shift() acts on geometry")
	DeclVar("ShiftRegions", &ShiftRegions, "Whether Shift() acts on regions")
	DeclVar("TotalShift", &TotalShift, "Amount by which the simulation has been shifted (m).")

	DeclVar("GammaLL", &GammaLL, "Gyromagnetic ratio in rad/Ts")
	DeclVar("DisableZhangLiTorque", &DisableZhangLiTorque, "Disables Zhang-Li torque (default=false)")
	DeclVar("DisableSlonczewskiTorque", &DisableSlonczewskiTorque, "Disables Slonczewski torque (default=false)")
	DeclVar("DoPrecess", &Precess, "Enables LL precession (default=true)")
}
