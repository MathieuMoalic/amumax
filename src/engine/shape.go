package engine

import (
	"image"
	"math"

	"github.com/MathieuMoalic/amumax/src/utils"
)

type shapeList struct {
	e *engineState
}

func newShape(engineState *engineState) *shapeList {
	s := &shapeList{e: engineState}

	s.e.script.RegisterFunction("Wave", s.wave)
	s.e.script.RegisterFunction("Ellipsoid", s.ellipsoid)
	s.e.script.RegisterFunction("Ellipse", s.ellipse)
	s.e.script.RegisterFunction("Cone", s.cone)
	s.e.script.RegisterFunction("Circle", s.circle)
	s.e.script.RegisterFunction("Cylinder", s.cylinder)
	s.e.script.RegisterFunction("Cuboid", s.cuboid)
	s.e.script.RegisterFunction("Rect", s.rect)
	s.e.script.RegisterFunction("Triangle", s.triangle)
	s.e.script.RegisterFunction("RTriangle", s.rTriangle)
	s.e.script.RegisterFunction("Hexagon", s.hexagon)
	s.e.script.RegisterFunction("Diamond", s.diamond)
	s.e.script.RegisterFunction("Squircle", s.squircle)
	s.e.script.RegisterFunction("Square", s.square)
	s.e.script.RegisterFunction("XRange", s.xRange)
	s.e.script.RegisterFunction("YRange", s.yRange)
	s.e.script.RegisterFunction("ZRange", s.zRange)
	s.e.script.RegisterFunction("Universe", s.universe)
	s.e.script.RegisterFunction("ImageShape", s.imageShape)
	s.e.script.RegisterFunction("GrainRoughness", s.grainRoughness)
	s.e.script.RegisterFunction("Layers", s.layers)
	s.e.script.RegisterFunction("Layer", s.layer)
	s.e.script.RegisterFunction("Cell", s.cell)

	return s
}

// geometrical shape for setting sample geometry
type shape func(x, y, z float64) bool

// wave with given diameters
func (s *shapeList) wave(period, amin, amax float64) shape {
	return func(x, y, z float64) bool {
		wavex := (math.Cos(x/period*2*math.Pi)/2 - 0.5) * (amax - amin) / 2
		return y > wavex-amin/2 && y < -wavex+amin/2
	}
}

// ellipsoid with given diameters
func (s *shapeList) ellipsoid(diamx, diamy, diamz float64) shape {
	return func(x, y, z float64) bool {
		return utils.Sqr64(x/diamx)+utils.Sqr64(y/diamy)+utils.Sqr64(z/diamz) <= 0.25
	}
}

func (s *shapeList) ellipse(diamx, diamy float64) shape {
	return s.ellipsoid(diamx, diamy, math.Inf(1))
}

// 3D cone with base at z=0 and vertex at z=height.
func (s *shapeList) cone(diam, height float64) shape {
	return func(x, y, z float64) bool {
		return (height-z)*z >= 0 && utils.Sqr64(x/diam)+utils.Sqr64(y/diam) <= 0.25*utils.Sqr64(1-z/height)
	}
}

func (s *shapeList) circle(diam float64) shape {
	return s.cylinder(diam, math.Inf(1))
}

// cylinder along z.
func (s *shapeList) cylinder(diam, height float64) shape {
	return func(x, y, z float64) bool {
		return z <= height/2 && z >= -height/2 &&
			utils.Sqr64(x/diam)+utils.Sqr64(y/diam) <= 0.25
	}
}

// 3D Rectangular slab with given sides.
func (s *shapeList) cuboid(sidex, sidey, sidez float64) shape {
	return func(x, y, z float64) bool {
		rx, ry, rz := sidex/2, sidey/2, sidez/2
		return x < rx && x > -rx && y < ry && y > -ry && z < rz && z > -rz
	}
}

// 2D Rectangle with given sides.
func (s *shapeList) rect(sidex, sidey float64) shape {
	return func(x, y, z float64) bool {
		rx, ry := sidex/2, sidey/2
		return x < rx && x > -rx && y < ry && y > -ry
	}
}

// Equilateral triangle with given sides.
func (s *shapeList) triangle(side float64) shape {
	return func(x, y, z float64) bool {
		c := math.Sqrt(3)
		return y > -side/(2*c) && y < x*c+side/c && y < -x*c+side/c
	}
}

// Rounded Equilateral triangle with given sides.
func (s *shapeList) rTriangle(side, diam float64) shape {
	return func(x, y, z float64) bool {
		c := math.Sqrt(3)
		return y > -side/(2*c) && y < x*c+side/c && y < -x*c+side/c && math.Sqrt(utils.Sqr64(x)+utils.Sqr64(y)) < diam/2
	}
}

// hexagon with given sides.
func (s *shapeList) hexagon(side float64) shape {
	return func(x, y, z float64) bool {
		a, b := math.Sqrt(3), math.Sqrt(3)*side
		return y < b/2 && y < -a*x+b && y > a*x-b && y > -b/2 && y > -a*x-b && y < a*x+b
	}
}

// diamond with given sides.
func (s *shapeList) diamond(sidex, sidey float64) shape {
	return func(x, y, z float64) bool {
		a, b := sidey/sidex, sidey/2
		return y < a*x+b && y < -a*x+b && y > a*x-b && y > -a*x-b
	}
}

// squircle creates a 3D rounded rectangle (a generalized squircle) with specified side lengths and thickness.
func (s *shapeList) squircle(sidex, sidey, sidez, a float64) shape {
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
func (s *shapeList) square(side float64) shape {
	return s.rect(side, side)
}

// All cells with x-coordinate between a and b
func (s *shapeList) xRange(a, b float64) shape {
	return func(x, y, z float64) bool {
		return x >= a && x < b
	}
}

// All cells with y-coordinate between a and b
func (s *shapeList) yRange(a, b float64) shape {
	return func(x, y, z float64) bool {
		return y >= a && y < b
	}
}

// All cells with z-coordinate between a and b
func (s *shapeList) zRange(a, b float64) shape {
	return func(x, y, z float64) bool {
		return z >= a && z < b
	}
}

// Cell layers #a (inclusive) up to #b (exclusive).
func (s *shapeList) layers(a, b int) shape {
	Nz := s.e.mesh.Nz
	if a < 0 || a > Nz || b < 0 || b < a {
		s.e.log.ErrAndExit("layers %d:%d out of bounds (0 - %d)", a, b, Nz)
	}
	dz := s.e.mesh.Dz
	z1 := s.e.mesh.Index2Coord(0, 0, a)[Z] - dz/2
	z2 := s.e.mesh.Index2Coord(0, 0, b)[Z] - dz/2
	return s.zRange(z1, z2)
}

func (s *shapeList) layer(index int) shape {
	return s.layers(index, index+1)
}

// Single cell with given index
func (s *shapeList) cell(ix, iy, iz int) shape {
	dx, dy, dz := s.e.mesh.GetDi()
	pos := s.e.mesh.Index2Coord(ix, iy, iz)
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

func (s *shapeList) universe() shape {
	return s.universeInner
}

// The entire space.
func (s *shapeList) universeInner(x, y, z float64) bool {
	return true
}

func (s *shapeList) imageShape(fname string) shape {
	r, err1 := s.e.fs.Open(fname)
	if err1 != nil {
		s.e.log.ErrAndExit("Error opening image file: %s: %s", fname, err1)
	}
	defer r.Close()
	img, _, err2 := image.Decode(r)
	if err2 != nil {
		s.e.log.ErrAndExit("Error decoding image file: %s: %s", fname, err2)
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
	dx, dy, _ := s.e.mesh.GetDi()
	Nx, Ny, _ := s.e.mesh.GetNi()
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

func (s *shapeList) grainRoughness(grainsize, zmin, zmax float64, seed int) shape {
	s.e.grains.voronoi(grainsize, 0, 256, seed)
	return func(x, y, z float64) bool {
		if z <= zmin {
			return true
		}
		if z >= zmax {
			return false
		}
		r := s.e.grains.getRegion(x, y, z)
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
		return s(utils.Fmod(x, periodX), utils.Fmod(y, periodY), utils.Fmod(z, periodZ))
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
