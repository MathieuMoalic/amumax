package new_engine

import (
	"image"
	"math"
)

type Shape struct {
	EngineState *EngineStateStruct
	// shape       func(x, y, z float64) bool
}

func NewShape(engineState *EngineStateStruct) *Shape {
	s := &Shape{EngineState: engineState}

	s.EngineState.world.RegisterFunction("Wave", s.wave)
	s.EngineState.world.RegisterFunction("Ellipsoid", s.ellipsoid)
	s.EngineState.world.RegisterFunction("Ellipse", s.ellipse)
	s.EngineState.world.RegisterFunction("Cone", s.cone)
	s.EngineState.world.RegisterFunction("Circle", s.circle)
	s.EngineState.world.RegisterFunction("Cylinder", s.cylinder)
	s.EngineState.world.RegisterFunction("Cuboid", s.cuboid)
	s.EngineState.world.RegisterFunction("Rect", s.rect)
	s.EngineState.world.RegisterFunction("Triangle", s.triangle)
	s.EngineState.world.RegisterFunction("RTriangle", s.rTriangle)
	s.EngineState.world.RegisterFunction("Hexagon", s.hexagon)
	s.EngineState.world.RegisterFunction("Diamond", s.diamond)
	s.EngineState.world.RegisterFunction("Squircle", s.squircle)
	s.EngineState.world.RegisterFunction("Square", s.square)
	s.EngineState.world.RegisterFunction("XRange", s.xRange)
	s.EngineState.world.RegisterFunction("YRange", s.yRange)
	s.EngineState.world.RegisterFunction("ZRange", s.zRange)
	s.EngineState.world.RegisterFunction("Universe", s.universe)
	s.EngineState.world.RegisterFunction("ImageShape", s.imageShape)
	s.EngineState.world.RegisterFunction("GrainRoughness", s.grainRoughness)
	s.EngineState.world.RegisterFunction("Layers", s.layers)
	s.EngineState.world.RegisterFunction("Layer", s.layer)
	s.EngineState.world.RegisterFunction("Cell", s.cell)

	return s
}

// geometrical shape for setting sample geometry
type shape func(x, y, z float64) bool

func sign(x float64) float64 {
	switch {
	case x > 0:
		return 1
	case x < 0:
		return -1
	default:
		return 0
	}
}

// wave with given diameters
func (s *Shape) wave(period, amin, amax float64) shape {
	return func(x, y, z float64) bool {
		wavex := (math.Cos(x/period*2*math.Pi)/2 - 0.5) * (amax - amin) / 2
		return y > wavex-amin/2 && y < -wavex+amin/2
	}
}

// ellipsoid with given diameters
func (s *Shape) ellipsoid(diamx, diamy, diamz float64) shape {
	return func(x, y, z float64) bool {
		return sqr64(x/diamx)+sqr64(y/diamy)+sqr64(z/diamz) <= 0.25
	}
}

func (s *Shape) ellipse(diamx, diamy float64) shape {
	return s.ellipsoid(diamx, diamy, math.Inf(1))
}

// 3D cone with base at z=0 and vertex at z=height.
func (s *Shape) cone(diam, height float64) shape {
	return func(x, y, z float64) bool {
		return (height-z)*z >= 0 && sqr64(x/diam)+sqr64(y/diam) <= 0.25*sqr64(1-z/height)
	}
}

func (s *Shape) circle(diam float64) shape {
	return s.cylinder(diam, math.Inf(1))
}

// cylinder along z.
func (s *Shape) cylinder(diam, height float64) shape {
	return func(x, y, z float64) bool {
		return z <= height/2 && z >= -height/2 &&
			sqr64(x/diam)+sqr64(y/diam) <= 0.25
	}
}

// 3D Rectangular slab with given sides.
func (s *Shape) cuboid(sidex, sidey, sidez float64) shape {
	return func(x, y, z float64) bool {
		rx, ry, rz := sidex/2, sidey/2, sidez/2
		return x < rx && x > -rx && y < ry && y > -ry && z < rz && z > -rz
	}
}

// 2D Rectangle with given sides.
func (s *Shape) rect(sidex, sidey float64) shape {
	return func(x, y, z float64) bool {
		rx, ry := sidex/2, sidey/2
		return x < rx && x > -rx && y < ry && y > -ry
	}
}

// Equilateral triangle with given sides.
func (s *Shape) triangle(side float64) shape {
	return func(x, y, z float64) bool {
		c := math.Sqrt(3)
		return y > -side/(2*c) && y < x*c+side/c && y < -x*c+side/c
	}
}

// Rounded Equilateral triangle with given sides.
func (s *Shape) rTriangle(side, diam float64) shape {
	return func(x, y, z float64) bool {
		c := math.Sqrt(3)
		return y > -side/(2*c) && y < x*c+side/c && y < -x*c+side/c && math.Sqrt(sqr64(x)+sqr64(y)) < diam/2
	}
}

// hexagon with given sides.
func (s *Shape) hexagon(side float64) shape {
	return func(x, y, z float64) bool {
		a, b := math.Sqrt(3), math.Sqrt(3)*side
		return y < b/2 && y < -a*x+b && y > a*x-b && y > -b/2 && y > -a*x-b && y < a*x+b
	}
}

// diamond with given sides.
func (s *Shape) diamond(sidex, sidey float64) shape {
	return func(x, y, z float64) bool {
		a, b := sidey/sidex, sidey/2
		return y < a*x+b && y < -a*x+b && y > a*x-b && y > -a*x-b
	}
}

// squircle creates a 3D rounded rectangle (a generalized squircle) with specified side lengths and thickness.
func (s *Shape) squircle(sidex, sidey, sidez, a float64) shape {
	// r := math.Min(sidex, sidey) / 2
	return func(x, y, z float64) bool {
		normX := x / (sidex / 2)
		normY := y / (sidey / 2)

		value := normX*normX + normY*normY - a*normX*normX*normY*normY

		if math.Abs(x) > sidex/2 && math.Abs(y) > sidey/2 {
			return false
		} else {
			inSquircleXY := value <= 1
			rz := sidez / 2
			inThickness := z >= -rz && z <= rz
			return inSquircleXY && inThickness
		}
	}
}

// 2D square with given side.
func (s *Shape) square(side float64) shape {
	return s.rect(side, side)
}

// All cells with x-coordinate between a and b
func (s *Shape) xRange(a, b float64) shape {
	return func(x, y, z float64) bool {
		return x >= a && x < b
	}
}

// All cells with y-coordinate between a and b
func (s *Shape) yRange(a, b float64) shape {
	return func(x, y, z float64) bool {
		return y >= a && y < b
	}
}

// All cells with z-coordinate between a and b
func (s *Shape) zRange(a, b float64) shape {
	return func(x, y, z float64) bool {
		return z >= a && z < b
	}
}

// Cell layers #a (inclusive) up to #b (exclusive).
func (s *Shape) layers(a, b int) shape {
	Nz := s.EngineState.mesh.Nz
	if a < 0 || a > Nz || b < 0 || b < a {
		s.EngineState.log.ErrAndExit("layers %d:%d out of bounds (0 - %d)", a, b, Nz)
	}
	dz := s.EngineState.mesh.Dz
	z1 := s.EngineState.utils.Index2Coord(0, 0, a)[Z] - dz/2
	z2 := s.EngineState.utils.Index2Coord(0, 0, b)[Z] - dz/2
	return s.zRange(z1, z2)
}

func (s *Shape) layer(index int) shape {
	return s.layers(index, index+1)
}

// Single cell with given index
func (s *Shape) cell(ix, iy, iz int) shape {
	dx, dy, dz := s.EngineState.mesh.GetD()
	pos := s.EngineState.utils.Index2Coord(ix, iy, iz)
	x1 := pos[X] - dx/2
	y1 := pos[Y] - dy/2
	z1 := pos[Z] - dz/2
	x2 := pos[X] + dx/2
	y2 := pos[Y] + dy/2
	z2 := pos[Z] + dz/2
	return func(x, y, z float64) bool {
		return x > x1 && x < x2 &&
			y > y1 && y < y2 &&
			z > z1 && z < z2
	}
}

func (s *Shape) universe() shape {
	return s.universeInner
}

// The entire space.
func (s *Shape) universeInner(x, y, z float64) bool {
	return true
}

func (s *Shape) imageShape(fname string) shape {
	r, err1 := s.EngineState.fs.Open(fname)
	if err1 != nil {
		s.EngineState.log.ErrAndExit("Error opening image file: %s: %s", fname, err1)
	}
	defer r.Close()
	img, _, err2 := image.Decode(r)
	if err2 != nil {
		s.EngineState.log.ErrAndExit("Error decoding image file: %s: %s", fname, err2)
	}

	width := img.Bounds().Max.X
	height := img.Bounds().Max.Y

	// decode image into bool matrix for fast pixel lookup
	inside := make([][]bool, height)
	for iy := range inside {
		inside[iy] = make([]bool, width)
	}
	for iy := 0; iy < height; iy++ {
		for ix := 0; ix < width; ix++ {
			r, g, b, a := img.At(ix, height-1-iy).RGBA()
			if a > 128 && r+g+b < (0xFFFF*3)/2 {
				inside[iy][ix] = true
			}
		}
	}

	// stretch the image onto the gridsize
	dx, dy, _ := s.EngineState.mesh.GetD()
	Nx, Ny, _ := s.EngineState.mesh.GetN()
	w, h := float64(width), float64(height)
	return func(x, y, z float64) bool {
		ix := int((w/float64(Nx))*(x/dx) + 0.5*w)
		iy := int((h/float64(Ny))*(y/dy) + 0.5*h)
		if ix < 0 || ix >= width || iy < 0 || iy >= height {
			return false
		} else {
			return inside[iy][ix]
		}
	}
}

func (s *Shape) grainRoughness(grainsize, zmin, zmax float64, seed int) shape {
	s.EngineState.grains.voronoi(grainsize, 0, 256, seed)
	return func(x, y, z float64) bool {
		if z <= zmin {
			return true
		}
		if z >= zmax {
			return false
		}
		r := s.EngineState.grains.RegionOf(x, y, z)
		return (z-zmin)/(zmax-zmin) < (float64(r) / 256)
	}
}

// Transl returns a translated copy of the shape.
func (s shape) Transl(dx, dy, dz float64) shape {
	return func(x, y, z float64) bool {
		return s(x-dx, y-dy, z-dz)
	}
}

// Infinitely repeats the shape with given period in x, y, z.
// A period of 0 or infinity means no repetition.
func (s shape) Repeat(periodX, periodY, periodZ float64) shape {
	return func(x, y, z float64) bool {
		return s(fmod(x, periodX), fmod(y, periodY), fmod(z, periodZ))
	}
}

func fmod(a, b float64) float64 {
	if b == 0 || math.IsInf(b, 1) {
		return a
	}
	if math.Abs(a) > b/2 {
		return sign(a) * (math.Mod(math.Abs(a+b/2), b) - b/2)
	} else {
		return a
	}
}

// Scale returns a scaled copy of the shape.
func (s shape) Scale(sx, sy, sz float64) shape {
	return func(x, y, z float64) bool {
		return s(x/sx, y/sy, z/sz)
	}
}

// Rotates the shape around the Z-axis, over θ radians.
func (s shape) RotZ(θ float64) shape {
	cos := math.Cos(θ)
	sin := math.Sin(θ)
	return func(x, y, z float64) bool {
		x_ := x*cos + y*sin
		y_ := -x*sin + y*cos
		return s(x_, y_, z)
	}
}

// Rotates the shape around the Y-axis, over θ radians.
func (s shape) RotY(θ float64) shape {
	cos := math.Cos(θ)
	sin := math.Sin(θ)
	return func(x, y, z float64) bool {
		x_ := x*cos - z*sin
		z_ := x*sin + z*cos
		return s(x_, y, z_)
	}
}

// Rotates the shape around the X-axis, over θ radians.
func (s shape) RotX(θ float64) shape {
	cos := math.Cos(θ)
	sin := math.Sin(θ)
	return func(x, y, z float64) bool {
		y_ := y*cos + z*sin
		z_ := -y*sin + z*cos
		return s(x, y_, z_)
	}
}

// Union of shapes a and b (logical OR).
func (a shape) Add(b shape) shape {
	return func(x, y, z float64) bool {
		return a(x, y, z) || b(x, y, z)
	}
}

// Intersection of shapes a and b (logical AND).
func (a shape) Intersect(b shape) shape {
	return func(x, y, z float64) bool {
		return a(x, y, z) && b(x, y, z)
	}
}

// Inverse (outside) of shape (logical NOT).
func (s shape) Inverse() shape {
	return func(x, y, z float64) bool {
		return !s(x, y, z)
	}
}

// Removes b from a (logical a AND NOT b)
func (a shape) Sub(b shape) shape {
	return func(x, y, z float64) bool {
		return a(x, y, z) && !b(x, y, z)
	}
}

// Logical XOR of shapes a and b
func (a shape) Xor(b shape) shape {
	return func(x, y, z float64) bool {
		A, B := a(x, y, z), b(x, y, z)
		return (A || B) && !(A && B)
	}
}

func sqr64(x float64) float64 { return x * x }
