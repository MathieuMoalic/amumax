package script

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/MathieuMoalic/amumax/src/chunk"
	"github.com/MathieuMoalic/amumax/src/fsutil"
	"github.com/MathieuMoalic/amumax/src/geometry"
	"github.com/MathieuMoalic/amumax/src/grains"
	"github.com/MathieuMoalic/amumax/src/mag_config"
	"github.com/MathieuMoalic/amumax/src/magnetization"
	"github.com/MathieuMoalic/amumax/src/mesh"
	"github.com/MathieuMoalic/amumax/src/metadata"
	"github.com/MathieuMoalic/amumax/src/regions"
	"github.com/MathieuMoalic/amumax/src/saved_quantities"
	"github.com/MathieuMoalic/amumax/src/shape"
	"github.com/MathieuMoalic/amumax/src/solver"
	"github.com/MathieuMoalic/amumax/src/table"
	"github.com/MathieuMoalic/amumax/src/utils"
	"github.com/MathieuMoalic/amumax/src/window_shift"
)

func (p *ScriptParser) AddToScopeAll(
	fs *fsutil.FileSystem,
	mesh *mesh.Mesh,
	geometry *geometry.Geometry,
	grains *grains.Grains,
	mag_config *mag_config.ConfigList,
	magnetization *magnetization.Magnetization,
	metadata *metadata.Metadata,
	regions *regions.Regions,
	saved_quantities *saved_quantities.SavedQuantities,
	solver *solver.Solver,
	table *table.Table,
	window_shift *window_shift.WindowShift,
	shape *shape.ShapeList,
) {
	p.RegisterFunction("Flush", fs.Drain, "Flush all pending output to disk.")
	// p.RegisterFunction("AutoSaveOvf", autoSaveOVF, "Auto save space-dependent quantity every period (s).")
	// p.RegisterFunction("AutoSnapshot", autoSnapshot, "Auto save image of quantity every period (s).")
	p.RegisterFunction("Chunk", chunk.CreateRequestedChunk, "")
	p.RegisterFunction("Uniform", mag_config.Uniform, "Uniform magnetization in given direction")
	p.RegisterFunction("Vortex", mag_config.Vortex, "Vortex magnetization with given circulation and core polarization")
	p.RegisterFunction("Antivortex", mag_config.AntiVortex, "Antivortex magnetization with given circulation and core polarization")
	p.RegisterFunction("Radial", mag_config.Radial, "Radial magnetization with given charge and core polarization")
	p.RegisterFunction("NeelSkyrmion", mag_config.NeelSkyrmion, "Néél skyrmion magnetization with given charge and core polarization")
	p.RegisterFunction("BlochSkyrmion", mag_config.BlochSkyrmion, "Bloch skyrmion magnetization with given chirality and core polarization")
	p.RegisterFunction("TwoDomain", mag_config.TwoDomain, "Twodomain magnetization with with given magnetization in left domain, wall, and right domain")
	p.RegisterFunction("VortexWall", mag_config.VortexWall, "Vortex wall magnetization with given mx in left and right domain and core circulation and polarization")
	p.RegisterFunction("RandomMag", mag_config.RandomMag, "Random magnetization")
	p.RegisterFunction("RandomMagSeed", mag_config.RandomMagSeed, "Random magnetization with given seed")
	p.RegisterFunction("Conical", mag_config.Conical, "Conical state for given wave vector, cone direction, and cone angle")
	p.RegisterFunction("Helical", mag_config.Helical, "Helical state for given wave vector")
	// p.RegisterFunction("Crop", crop, "Crops a quantity to cell ranges [x1,x2[, [y1,y2[, [z1,z2[")
	// p.RegisterFunction("CropX", cropX, "Crops a quantity to cell ranges [x1,x2[")
	// p.RegisterFunction("CropY", cropY, "Crops a quantity to cell ranges [y1,y2[")
	// p.RegisterFunction("CropZ", cropZ, "Crops a quantity to cell ranges [z1,z2[")
	// p.RegisterFunction("CropLayer", cropLayer, "Crops a quantity to a single layer")
	// p.RegisterFunction("CropRegion", cropRegion, "Crops a quantity to a region")

	// p.RegisterFunction("AddFieldTerm", addFieldTerm, "Add an expression to B_eff.")
	// p.RegisterFunction("AddEdensTerm", addEdensTerm, "Add an expression to Edens.")
	// p.RegisterFunction("Add", add, "Add two quantities")
	// p.RegisterFunction("Madd", madd, "Weighted addition: Madd(Q1,Q2,c1,c2) = c1*Q1 + c2*Q2")
	// p.RegisterFunction("Dot", dotProductFunc, "Dot product of two vector quantities")
	// p.RegisterFunction("Cross", cross, "Cross product of two vector quantities")
	// p.RegisterFunction("Mul", mul, "Point-wise product of two quantities")
	// p.RegisterFunction("MulMV", mulMV, "Matrix-Vector product: MulMV(AX, AY, AZ, m) = (AX·m, AY·m, AZ·m). "+
	// 	"The arguments Ax, Ay, Az and m are quantities with 3 componets.")
	// p.RegisterFunction("Div", div, "Point-wise division of two quantities")
	// p.RegisterFunction("Const", constScalar, "Constant, uniform number")
	// p.RegisterFunction("ConstVector", constVector, "Constant, uniform vector")
	// p.RegisterFunction("Shifted", shiftedQuant, "Shifted quantity")
	// p.RegisterFunction("Masked", maskedQuant, "Mask quantity with shape")
	// p.RegisterFunction("Normalized", normalizedQuant, "Normalize quantity")
	// p.RegisterFunction("RemoveCustomFields", removeCustomFields, "Removes all custom fields again")

	// p.RegisterFunction("ext_ScaleExchange", scaleInterExchange, "Re-scales exchange coupling between two regions.")
	// p.RegisterFunction("ext_InterExchange", interExchange, "Sets exchange coupling between two regions.")
	// p.RegisterFunction("ext_ScaleDind", scaleInterDind, "Re-scales Dind coupling between two regions.")
	// p.RegisterFunction("ext_InterDind", interDind, "Sets Dind coupling between two regions.")
	// p.RegisterFunction("ext_centerBubble", centerBubble, "centerBubble shifts m after each step to keep the bubble position close to the center of the window")
	// p.RegisterFunction("ext_centerWall", centerWall, "centerWall(c) shifts m after each step to keep m_c close to zero")
	// p.RegisterFunction("ext_make3dgrains", voronoi3d, "3D Voronoi tesselation over shape (grain size, starting region number, num regions, shape, seed)")
	p.RegisterFunction("ext_makegrains", grains.Voronoi, "Voronoi tesselation (grain size, num regions)")
	// p.RegisterFunction("ext_rmSurfaceCharge", removeLRSurfaceCharge, "Compensate magnetic charges on the left and right sides of an in-plane magnetized wire. Arguments: region, mx on left and right side, resp.")
	p.RegisterFunction("SetGeom", geometry.SetGeom, "Sets the geometry to a given shape")
	// p.RegisterFunction("Minimize", minimize, "Use steepest conjugate gradient method to minimize the total energy")

	// p.RegisterFunction("ReCreateMesh", ReCreateMesh, "")
	p.RegisterFunction("SmoothMesh", mesh.SmoothMesh, "")
	p.RegisterFunction("SetMesh", mesh.SetMesh, "")
	p.RegisterFunction("SetGridSize", mesh.SetGridSize, "")
	p.RegisterFunction("SetCellSize", mesh.SetCellSize, "")
	p.RegisterFunction("SetPBC", mesh.SetPBC, "")
	p.RegisterFunction("SetTotalSize", mesh.SetTotalSize, "")

	// p.RegisterFunction("DefRegion", DefRegion, "Define a material region with given index (0-255) and shape")
	// p.RegisterFunction("RedefRegion", RedefRegion, "Reassign all cells with a given region (first argument) to a new region (second argument)")
	// p.RegisterFunction("ShapeFromRegion", ShapeFromRegion, "")
	// p.RegisterFunction("DefRegionCell", DefRegionCell, "Set a material region (first argument) in one cell "+
	// 	"by the index of the cell (last three arguments)")
	// p.RegisterFunction("Relax", relax, "Try to minimize the total energy")
	// p.RegisterFunction("Run", run, "Run the simulation for a time in seconds")
	// p.RegisterFunction("RunWithoutPrecession", runWithoutPrecession, "Run the simulation for a time in seconds with precession disabled")
	// p.RegisterFunction("Steps", steps, "Run the simulation for a number of time steps")
	// p.RegisterFunction("RunWhile", runWhile, "Run while condition function is true")
	// p.RegisterFunction("SetSolver", setSolver, "Set solver type. 1:Euler, 2:Heun, 3:Bogaki-Shampine, 4: Runge-Kutta (RK45), 5: Dormand-Prince, 6: Fehlberg, -1: Backward Euler")
	// p.RegisterFunction("Exit", Exit, "Exit from the program")
	// p.RegisterFunction("RunShell", runShell, "Run a shell command")

	// p.RegisterFunction("SaveOvf", saveOVF, "Save space-dependent quantity once, with auto filename")
	// p.RegisterFunction("SaveOvfAs", saveAsOVF, "Save space-dependent quantity with custom filename")
	// p.RegisterFunction("Snapshot", snapshot, "Save image of quantity")
	// p.RegisterFunction("SnapshotAs", snapshotAs, "Save image of quantity with custom filename")

	p.RegisterFunction("Ellipsoid", shape.Ellipsoid, "3D Ellipsoid with axes in meter")
	p.RegisterFunction("Ellipse", shape.Ellipse, "2D Ellipse with axes in meter")
	p.RegisterFunction("Cone", shape.Cone, "3D Cone with diameter and height in meter. The base is at z=0. If the height is positive, the tip points in the +z direction.")
	p.RegisterFunction("Cylinder", shape.Cylinder, "3D Cylinder with diameter and height in meter")
	p.RegisterFunction("Circle", shape.Circle, "2D Circle with diameter in meter")
	p.RegisterFunction("Squircle", shape.Squircle, "2D Squircle with diameter in meter")
	p.RegisterFunction("Cuboid", shape.Cuboid, "Cuboid with sides in meter")
	p.RegisterFunction("Rect", shape.Rect, "2D rectangle with size in meter")
	p.RegisterFunction("Wave", shape.Wave, "Wave with (Period, Min amplitude and Max amplitude) in meter")
	p.RegisterFunction("Triangle", shape.Triangle, "Equilateral triangle with side in meter")
	p.RegisterFunction("RTriangle", shape.RTriangle, "Rounded Equilateral triangle with side in meter")
	p.RegisterFunction("Diamond", shape.Diamond, "Diamond with side in meter")
	p.RegisterFunction("Hexagon", shape.Hexagon, "Hexagon with side in meter")
	p.RegisterFunction("Square", shape.Square, "2D square with size in meter")
	p.RegisterFunction("XRange", shape.XRange, "Part of space between x1 (inclusive) and x2 (exclusive), in meter")
	p.RegisterFunction("YRange", shape.YRange, "Part of space between y1 (inclusive) and y2 (exclusive), in meter")
	p.RegisterFunction("ZRange", shape.ZRange, "Part of space between z1 (inclusive) and z2 (exclusive), in meter")
	p.RegisterFunction("Layers", shape.Layers, "Part of space between cell layer1 (inclusive) and layer2 (exclusive), in integer indices")
	p.RegisterFunction("Layer", shape.Layer, "Single layer (along z), by integer index starting from 0")
	p.RegisterFunction("Universe", shape.Universe, "Entire space")
	p.RegisterFunction("Cell", shape.Cell, "Single cell with given integer index (i, j, k)")
	p.RegisterFunction("ImageShape", shape.ImageShape, "Use black/white image as shape")
	p.RegisterFunction("GrainRoughness", shape.GrainRoughness, "Grainy surface with different heights per grain "+
		"with a typical grain size (first argument), minimal height (second argument), and maximal "+
		"height (third argument). The last argument is a seed for the random number generator.")
	p.RegisterFunction("Shift", window_shift.ShiftX, "Shifts the simulation by +1/-1 cells along X")
	p.RegisterFunction("ShiftY", window_shift.ShiftY, "Shifts the simulation by +1/-1 cells along Y")
	// p.RegisterFunction("ThermSeed", thermSeed, "Set a random seed for thermal noise")

	// p.RegisterFunction("Expect", expect, "Used for automated tests: checks if a value is close enough to the expected value")
	// p.RegisterFunction("ExpectV", expectV, "Used for automated tests: checks if a vector is close enough to the expected value")
	// p.RegisterFunction("Fprintln", fprintln, "Print to file")
	// p.RegisterFunction("Sign", sign, "Signum function")
	// p.RegisterFunction("Vector", vector, "Constructs a vector with given components")
	p.RegisterFunction("Print", p.Print, "Print to standard output")
	// p.RegisterFunction("LoadFile", loadFile, "Load a zarr data file")
	// p.RegisterFunction("LoadOvfFile", loadOvfFile, "Load an ovf data file")
	p.RegisterFunction("Index2Coord", mesh.Index2Coord, "Convert cell index to x,y,z coordinate in meter")
	// p.RegisterFunction("NewSlice", newSlice, "Makes a 4D array with a specified number of components (first argument) "+
	// 	"and a specified size nx,ny,nz (remaining arguments)")
	// p.RegisterFunction("NewVectorMask", newVectorMask, "Makes a 3D array of vectors")
	// p.RegisterFunction("NewScalarMask", newScalarMask, "Makes a 3D array of scalars")
	// p.RegisterFunction("RegionFromCoordinate", regionFromCoordinate, "RegionFromCoordinate")

	p.RegisterFunction("AutoSaveAs", saved_quantities.AutoSaveAs, "Auto save space-dependent quantity every period (s) as the zarr standard.")
	p.RegisterFunction("AutoSaveAsChunk", saved_quantities.AutoSaveAsChunk, "Auto save space-dependent quantity every period (s) as the zarr standard.")
	p.RegisterFunction("AutoSave", saved_quantities.AutoSave, "Auto save space-dependent quantity every period (s) as the zarr standard.")
	p.RegisterFunction("SaveAs", saved_quantities.SaveAs, "Save space-dependent quantity as the zarr standard.")
	p.RegisterFunction("SaveAsChunk", saved_quantities.SaveAsChunk, "")
	p.RegisterFunction("Save", saved_quantities.Save, "Save space-dependent quantity as the zarr standard.")

	p.RegisterFunction("TableSave", table.Save, "Save the data table right now.")
	p.RegisterFunction("TableAdd", table.Add, "Save the data table periodically.")
	// p.RegisterFunction("TableAddVar", tableAddVar, "Save the data table periodically.")
	p.RegisterFunction("TableAddAs", table.AddAs, "Save the data table periodically.")
	p.RegisterFunction("TableAutoSave", table.AutoSave, "Save the data table periodically.")

	p.RegisterVariable("Nx", &mesh.Nx, "")
	p.RegisterVariable("Ny", &mesh.Ny, "")
	p.RegisterVariable("Nz", &mesh.Nz, "")
	p.RegisterVariable("dx", &mesh.Dx, "")
	p.RegisterVariable("dy", &mesh.Dy, "")
	p.RegisterVariable("dz", &mesh.Dz, "")
	p.RegisterVariable("Tx", &mesh.Tx, "")
	p.RegisterVariable("Ty", &mesh.Ty, "")
	p.RegisterVariable("Tz", &mesh.Tz, "")
	p.RegisterVariable("PBCx", &mesh.PBCx, "")
	p.RegisterVariable("PBCy", &mesh.PBCy, "")
	p.RegisterVariable("PBCz", &mesh.PBCz, "")

	p.RegisterVariable("geom", geometry, "")
	p.RegisterVariable("m", magnetization, "")
	p.RegisterVariable("t", &solver.Time, "Total simulated time (s)")

	// p.RegisterVariable("EnableDemag", &EnableDemag, "Enables/disables demag (default=true)")
	// p.RegisterVariable("DemagAccuracy", &DemagAccuracy, "Controls accuracy of demag kernel")

	p.RegisterVariable("step", &solver.NSteps, "Total number of time steps taken")
	p.RegisterVariable("MinDt", &solver.MinDt, "Minimum time step the solver can take (s)")
	p.RegisterVariable("MaxDt", &solver.MaxDt, "Maximum time step the solver can take (s)")
	p.RegisterVariable("MaxErr", &solver.MaxErr, "Maximum error per step the solver can tolerate (default = 1e-5)")
	// p.RegisterVariable("Headroom", &Headroom, "Solver headroom (default = 0.8)")
	p.RegisterVariable("FixDt", &solver.FixDt, "Set a fixed time step, 0 disables fixed step (which is the default)")
	// p.RegisterVariable("OpenBC", &OpenBC, "Use open boundary conditions (default=false)")
	// p.RegisterVariable("ext_BubbleMz", &BubbleMz, "Center magnetization 1.0 or -1.0  (default = 1.0)")
	p.RegisterVariable("EdgeSmooth", &geometry.EdgeSmooth, "Geometry edge smoothing with edgeSmooth^3 samples per cell, 0=staircase, ~8=very smooth")

	// p.RegisterVariable("MinimizerStop", &stopMaxDm, "Stopping max dM for Minimize")
	// p.RegisterVariable("MinimizerSamples", &dmSamples, "Number of max dM to collect for Minimize convergence check.")
	// p.RegisterVariable("MinimizeMaxSteps", &minimizeMaxSteps, "")
	// p.RegisterVariable("MinimizeMaxTimeSeconds", &minimizeMaxTimeSeconds, "")
	// p.RegisterVariable("RelaxTorqueThreshold", &relaxTorqueThreshold, "MaxTorque threshold for relax(). If set to -1 (default), relax() will stop when the average torque is steady or increasing.")
	// p.RegisterVariable("SnapshotFormat", &snapshotFormat, "Image format for snapshots: jpg, png or gif.")

	p.RegisterVariable("ShiftMagL", &window_shift.ShiftMagL, "Upon shift, insert this magnetization from the left")
	p.RegisterVariable("ShiftMagR", &window_shift.ShiftMagR, "Upon shift, insert this magnetization from the right")
	p.RegisterVariable("ShiftMagU", &window_shift.ShiftMagU, "Upon shift, insert this magnetization from the top")
	p.RegisterVariable("ShiftMagD", &window_shift.ShiftMagD, "Upon shift, insert this magnetization from the bottom")
	// p.RegisterVariable("ShiftM", &shiftM, "Whether Shift() acts on magnetization")
	// p.RegisterVariable("ShiftGeom", &shiftGeom, "Whether Shift() acts on geometry")
	// p.RegisterVariable("ShiftRegions", &shiftRegions, "Whether Shift() acts on regions")
	// p.RegisterVariable("TotalShift", &totalShift, "Amount by which the simulation has been shifted (m).")
	// p.RegisterVariable("EdgeCarryShift", &EdgeCarryShift, "Whether to use the current magnetization at the border for the cells inserted by Shift")

	// p.RegisterVariable("GammaLL", &gammaLL, "Gyromagnetic ratio in rad/Ts")
	// p.RegisterVariable("DisableZhangLiTorque", &disableZhangLiTorque, "Disables Zhang-Li torque (default=false)")
	// p.RegisterVariable("DisableSlonczewskiTorque", &disableSlonczewskiTorque, "Disables Slonczewski torque (default=false)")
	// p.RegisterVariable("DoPrecess", &precess, "Enables LL precession (default=true)")

	// p.RegisterVariable("PreviewXDataPoints", &PreviewXDataPoints, "Number of data points in the x direction for the 2D/3D preview")
	// p.RegisterVariable("PreviewYDataPoints", &PreviewYDataPoints, "Number of data points in the y direction for the 2D/3D preview")
}
func (p *ScriptParser) Print(msg ...interface{}) {
	p.log.Info("%v", utils.CustomFmt(msg))
}

// RegisterFunction registers a pre-defined function in the world.
func (p *ScriptParser) RegisterFunction(name string, function interface{}, doc string) {
	name = strings.ToLower(name)
	p.functionsScope[name] = p.wrapFunction(function, name)
}

// RegisterVariable registers a pre-defined variable in the world.
func (p *ScriptParser) RegisterVariable(name string, value interface{}, doc string) {
	name = strings.ToLower(name)
	if value == nil {
		p.log.ErrAndExit("Value is nil for variable: %s", name)
	}
	p.variablesScope[name] = value
}

// registerUserVariable registers a user-defined variable in the world.
func (p *ScriptParser) registerUserVariable(name string, value interface{}) {
	name = strings.ToLower(name)
	if existingValue, ok := p.variablesScope[name]; ok {
		switch ptr := existingValue.(type) {
		case *int:
			if v, ok := value.(int); ok {
				*ptr = v
			}
		case *float64:
			if v, ok := value.(float64); ok {
				*ptr = v
			}
		case *bool:
			if v, ok := value.(bool); ok {
				*ptr = v
			}
		default:
			p.log.PanicIfError(fmt.Errorf("unsupported type: %T", ptr))
		}
	} else {
		p.variablesScope[name] = value
	}
	if p.isMeshExpression(name) {
		p.initializeMeshIfReady()
	}
}

// wrapFunction creates a universal wrapper for any function.
func (p *ScriptParser) wrapFunction(fn interface{}, name string) func([]interface{}) (interface{}, error) {
	name = strings.ToLower(name)
	return func(args []interface{}) (interface{}, error) {
		fnValue := reflect.ValueOf(fn)
		fnType := fnValue.Type()

		// Ensure the provided function is callable
		if fnType.Kind() != reflect.Func {
			return nil, fmt.Errorf("provided argument is not a function")
		}

		numIn := fnType.NumIn()
		isVariadic := fnType.IsVariadic()
		numFixedArgs := numIn
		if isVariadic {
			numFixedArgs-- // The last parameter is variadic
		}

		// Check if the number of arguments is sufficient
		if (!isVariadic && len(args) != numIn) || (isVariadic && len(args) < numFixedArgs) {
			expectedArgs := numIn
			if isVariadic {
				expectedArgs = numFixedArgs
				return nil, fmt.Errorf(
					"%s expects at least %d arguments (%s), got %d",
					name,
					expectedArgs,
					p.formatFunctionSignature(fnType, name),
					len(args),
				)
			} else {
				return nil, fmt.Errorf(
					"%s expects %d arguments (%s), got %d",
					name,
					expectedArgs,
					p.formatFunctionSignature(fnType, name),
					len(args),
				)
			}
		}

		// Prepare arguments for the function call
		in := make([]reflect.Value, numFixedArgs)

		// Handle fixed arguments
		for i := 0; i < numFixedArgs; i++ {
			expectedType := fnType.In(i)
			if len(args) <= i {
				return nil, fmt.Errorf(
					"%s: missing argument for parameter %d\nExpected function signature: %s",
					name,
					i+1,
					p.formatFunctionSignature(fnType, name),
				)
			}
			arg := args[i]
			argVal := reflect.ValueOf(arg)

			// Check if the argument is assignable to the expected type
			if !argVal.Type().AssignableTo(expectedType) {
				if expectedType.Kind() == reflect.Interface && argVal.Type().Implements(expectedType) {
					// The argument implements the expected interface; proceed without conversion
				} else if argVal.Type().ConvertibleTo(expectedType) {
					argVal = argVal.Convert(expectedType)
				} else {
					return nil, fmt.Errorf(
						"%s: argument %d (%v) is not assignable to %s\nExpected function signature: %s",
						name,
						i+1,
						argVal.Type(),
						expectedType,
						p.formatFunctionSignature(fnType, name),
					)
				}
			}

			in[i] = argVal
		}

		// Handle variadic arguments
		if isVariadic {
			variadicType := fnType.In(numIn - 1).Elem() // Element type of variadic parameter
			numVariadicArgs := len(args) - numFixedArgs
			variadicSlice := reflect.MakeSlice(reflect.SliceOf(variadicType), numVariadicArgs, numVariadicArgs)
			for i := 0; i < numVariadicArgs; i++ {
				arg := args[numFixedArgs+i]
				argVal := reflect.ValueOf(arg)

				// Check if the argument is assignable to the variadic type
				if !argVal.Type().AssignableTo(variadicType) {
					if variadicType.Kind() == reflect.Interface && argVal.Type().Implements(variadicType) {
						// The argument implements the expected interface; proceed without conversion
					} else if argVal.Type().ConvertibleTo(variadicType) {
						argVal = argVal.Convert(variadicType)
					} else {
						return nil, fmt.Errorf(
							"%s: argument %d (%v) is not assignable to %s\nExpected function signature: %s",
							name,
							numFixedArgs+i+1,
							argVal.Type(),
							variadicType,
							p.formatFunctionSignature(fnType, name),
						)
					}
				}

				variadicSlice.Index(i).Set(argVal)
			}
			// Append the variadic slice to the arguments
			in = append(in, variadicSlice)
		}

		var out []reflect.Value
		// Call the function using reflection
		if isVariadic {
			out = fnValue.CallSlice(in)
		} else {
			out = fnValue.Call(in)
		}

		// Handle the function's return values
		switch len(out) {
		case 0:
			return nil, nil
		case 1:
			if fnType.Out(0).Name() == "error" {
				if !out[0].IsNil() {
					return nil, out[0].Interface().(error)
				}
				return nil, nil
			}
			return out[0].Interface(), nil
		case 2:
			var err error
			if !out[1].IsNil() {
				err = out[1].Interface().(error)
			}
			return out[0].Interface(), err
		default:
			return nil, fmt.Errorf("%s has unsupported number of return values: %d", name, len(out))
		}
	}
}

// formatFunctionSignature returns a string representation of the function's signature.
func (p *ScriptParser) formatFunctionSignature(fnType reflect.Type, name string) string {
	var sb strings.Builder
	sb.WriteString(name)
	sb.WriteString("(")
	numIn := fnType.NumIn()
	isVariadic := fnType.IsVariadic()
	for i := 0; i < numIn; i++ {
		if i > 0 {
			sb.WriteString(", ")
		}
		inType := fnType.In(i)
		if isVariadic && i == numIn-1 {
			sb.WriteString("...")
			sb.WriteString(inType.Elem().String())
		} else {
			sb.WriteString(inType.String())
		}
	}
	sb.WriteString(")")
	if fnType.NumOut() > 0 {
		sb.WriteString(" (")
		for i := 0; i < fnType.NumOut(); i++ {
			if i > 0 {
				sb.WriteString(", ")
			}
			outType := fnType.Out(i)
			sb.WriteString(outType.String())
		}
		sb.WriteString(")")
	}
	return sb.String()
}

func (p *ScriptParser) isMeshExpression(name string) bool {
	namesToCheck := []string{"Nx", "Ny", "Nz", "Dx", "Dy", "Dz", "Tx", "Ty", "Tz"}
	for _, v := range namesToCheck {
		if strings.EqualFold(v, name) {
			return true
		}
	}
	return false
}

func (p *ScriptParser) getVariable(name string) (interface{}, bool) {
	name = strings.ToLower(name)
	value, ok := p.variablesScope[name]
	return value, ok
}

func (p *ScriptParser) getFunction(name string) (interface{}, bool) {
	name = strings.ToLower(name)
	value, ok := p.functionsScope[name]
	return value, ok
}
