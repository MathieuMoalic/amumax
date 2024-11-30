package shape

import (
	"image"
	"math"

	"github.com/MathieuMoalic/amumax/src/fsutil"
	"github.com/MathieuMoalic/amumax/src/grains"
	"github.com/MathieuMoalic/amumax/src/log"
	"github.com/MathieuMoalic/amumax/src/mesh"
	"github.com/MathieuMoalic/amumax/src/script"
	"github.com/MathieuMoalic/amumax/src/utils"
)

type ShapeList struct {
	fs     *fsutil.FileSystem
	grains *grains.Grains
	mesh   *mesh.Mesh
	log    *log.Logs
	script *script.ScriptParser
}

func NewShape(mesh *mesh.Mesh, log *log.Logs, script *script.ScriptParser, fs *fsutil.FileSystem, grains *grains.Grains) *ShapeList {
	s := &ShapeList{mesh: mesh, log: log, script: script, fs: fs, grains: grains}

	s.script.RegisterFunction("Wave", s.wave)
	s.script.RegisterFunction("Ellipsoid", s.ellipsoid)
	s.script.RegisterFunction("Ellipse", s.ellipse)
	s.script.RegisterFunction("Cone", s.cone)
	s.script.RegisterFunction("Circle", s.circle)
	s.script.RegisterFunction("Cylinder", s.cylinder)
	s.script.RegisterFunction("Cuboid", s.cuboid)
	s.script.RegisterFunction("Rect", s.rect)
	s.script.RegisterFunction("Triangle", s.triangle)
	s.script.RegisterFunction("RTriangle", s.rTriangle)
	s.script.RegisterFunction("Hexagon", s.hexagon)
	s.script.RegisterFunction("Diamond", s.diamond)
	s.script.RegisterFunction("Squircle", s.squircle)
	s.script.RegisterFunction("Square", s.square)
	s.script.RegisterFunction("XRange", s.xRange)
	s.script.RegisterFunction("YRange", s.yRange)
	s.script.RegisterFunction("ZRange", s.zRange)
	s.script.RegisterFunction("Universe", s.universe)
	s.script.RegisterFunction("ImageShape", s.imageShape)
	s.script.RegisterFunction("GrainRoughness", s.grainRoughness)
	s.script.RegisterFunction("Layers", s.layers)
	s.script.RegisterFunction("Layer", s.layer)
	s.script.RegisterFunction("Cell", s.cell)

	return s
}

// geometrical Shape for setting sample geometry
type Shape func(x, y, z float64) bool

// wave with given diameters
func (s *ShapeList) wave(period, amin, amax float64) Shape {
	return func(x, y, z float64) bool {
		wavex := (math.Cos(x/period*2*math.Pi)/2 - 0.5) * (amax - amin) / 2
		return y > wavex-amin/2 && y < -wavex+amin/2
	}
}

// ellipsoid with given diameters
func (s *ShapeList) ellipsoid(diamx, diamy, diamz float64) Shape {
	return func(x, y, z float64) bool {
		return utils.Sqr64(x/diamx)+utils.Sqr64(y/diamy)+utils.Sqr64(z/diamz) <= 0.25
	}
}

func (s *ShapeList) ellipse(diamx, diamy float64) Shape {
	return s.ellipsoid(diamx, diamy, math.Inf(1))
}

// 3D cone with base at z=0 and vertex at z=height.
func (s *ShapeList) cone(diam, height float64) Shape {
	return func(x, y, z float64) bool {
		return (height-z)*z >= 0 && utils.Sqr64(x/diam)+utils.Sqr64(y/diam) <= 0.25*utils.Sqr64(1-z/height)
	}
}

func (s *ShapeList) circle(diam float64) Shape {
	return s.cylinder(diam, math.Inf(1))
}

// cylinder along z.
func (s *ShapeList) cylinder(diam, height float64) Shape {
	return func(x, y, z float64) bool {
		return z <= height/2 && z >= -height/2 &&
			utils.Sqr64(x/diam)+utils.Sqr64(y/diam) <= 0.25
	}
}

// 3D Rectangular slab with given sides.
func (s *ShapeList) cuboid(sidex, sidey, sidez float64) Shape {
	return func(x, y, z float64) bool {
		rx, ry, rz := sidex/2, sidey/2, sidez/2
		return x < rx && x > -rx && y < ry && y > -ry && z < rz && z > -rz
	}
}

// 2D Rectangle with given sides.
func (s *ShapeList) rect(sidex, sidey float64) Shape {
	return func(x, y, z float64) bool {
		rx, ry := sidex/2, sidey/2
		return x < rx && x > -rx && y < ry && y > -ry
	}
}

// Equilateral triangle with given sides.
func (s *ShapeList) triangle(side float64) Shape {
	return func(x, y, z float64) bool {
		c := math.Sqrt(3)
		return y > -side/(2*c) && y < x*c+side/c && y < -x*c+side/c
	}
}

// Rounded Equilateral triangle with given sides.
func (s *ShapeList) rTriangle(side, diam float64) Shape {
	return func(x, y, z float64) bool {
		c := math.Sqrt(3)
		return y > -side/(2*c) && y < x*c+side/c && y < -x*c+side/c && math.Sqrt(utils.Sqr64(x)+utils.Sqr64(y)) < diam/2
	}
}

// hexagon with given sides.
func (s *ShapeList) hexagon(side float64) Shape {
	return func(x, y, z float64) bool {
		a, b := math.Sqrt(3), math.Sqrt(3)*side
		return y < b/2 && y < -a*x+b && y > a*x-b && y > -b/2 && y > -a*x-b && y < a*x+b
	}
}

// diamond with given sides.
func (s *ShapeList) diamond(sidex, sidey float64) Shape {
	return func(x, y, z float64) bool {
		a, b := sidey/sidex, sidey/2
		return y < a*x+b && y < -a*x+b && y > a*x-b && y > -a*x-b
	}
}

// squircle creates a 3D rounded rectangle (a generalized squircle) with specified side lengths and thickness.
func (s *ShapeList) squircle(sidex, sidey, sidez, a float64) Shape {
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
func (s *ShapeList) square(side float64) Shape {
	return s.rect(side, side)
}

// All cells with x-coordinate between a and b
func (s *ShapeList) xRange(a, b float64) Shape {
	return func(x, y, z float64) bool {
		return x >= a && x < b
	}
}

// All cells with y-coordinate between a and b
func (s *ShapeList) yRange(a, b float64) Shape {
	return func(x, y, z float64) bool {
		return y >= a && y < b
	}
}

// All cells with z-coordinate between a and b
func (s *ShapeList) zRange(a, b float64) Shape {
	return func(x, y, z float64) bool {
		return z >= a && z < b
	}
}

// Cell layers #a (inclusive) up to #b (exclusive).
func (s *ShapeList) layers(a, b int) Shape {
	Nz := s.mesh.Nz
	if a < 0 || a > Nz || b < 0 || b < a {
		s.log.ErrAndExit("layers %d:%d out of bounds (0 - %d)", a, b, Nz)
	}
	dz := s.mesh.Dz
	z1 := s.mesh.Index2Coord(0, 0, a)[2] - dz/2
	z2 := s.mesh.Index2Coord(0, 0, b)[2] - dz/2
	return s.zRange(z1, z2)
}

func (s *ShapeList) layer(index int) Shape {
	return s.layers(index, index+1)
}

// Single cell with given index
func (s *ShapeList) cell(ix, iy, iz int) Shape {
	dx, dy, dz := s.mesh.GetDi()
	pos := s.mesh.Index2Coord(ix, iy, iz)
	x1 := pos[0] - dx/2
	y1 := pos[1] - dy/2
	z1 := pos[2] - dz/2
	x2 := pos[0] + dx/2
	y2 := pos[1] + dy/2
	z2 := pos[2] + dz/2
	return func(x, y, z float64) bool {
		return x > x1 && x < x2 &&
			y > y1 && y < y2 &&
			z > z1 && z < z2
	}
}

func (s *ShapeList) universe() Shape {
	return Universe
}

func (s *ShapeList) imageShape(fname string) Shape {
	r, err1 := s.fs.Open(fname)
	if err1 != nil {
		s.log.ErrAndExit("Error opening image file: %s: %s", fname, err1)
	}
	defer r.Close()
	img, _, err2 := image.Decode(r)
	if err2 != nil {
		s.log.ErrAndExit("Error decoding image file: %s: %s", fname, err2)
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
	dx, dy, _ := s.mesh.GetDi()
	Nx, Ny, _ := s.mesh.GetNi()
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

func (s *ShapeList) grainRoughness(grainsize, zmin, zmax float64, seed int) Shape {
	s.grains.Voronoi(grainsize, 0, 256, seed)
	return func(x, y, z float64) bool {
		if z <= zmin {
			return true
		}
		if z >= zmax {
			return false
		}
		r := s.grains.GetRegion(x, y, z)
		return (z-zmin)/(zmax-zmin) < (float64(r) / 256)
	}
}

// The entire space.
func Universe(x, y, z float64) bool {
	return true
}

// Transl returns a translated copy of the shape.
func (s Shape) Transl(dx, dy, dz float64) Shape {
	return func(x, y, z float64) bool {
		return s(x-dx, y-dy, z-dz)
	}
}

// Infinitely repeats the shape with given period in x, y, z.
// A period of 0 or infinity means no repetition.
func (s Shape) Repeat(periodX, periodY, periodZ float64) Shape {
	return func(x, y, z float64) bool {
		return s(utils.Fmod(x, periodX), utils.Fmod(y, periodY), utils.Fmod(z, periodZ))
	}
}

// Scale returns a scaled copy of the shape.
func (s Shape) Scale(sx, sy, sz float64) Shape {
	return func(x, y, z float64) bool {
		return s(x/sx, y/sy, z/sz)
	}
}

// Rotates the shape around the 2-axis, over θ radians.
func (s Shape) RotZ(θ float64) Shape {
	cos := math.Cos(θ)
	sin := math.Sin(θ)
	return func(x, y, z float64) bool {
		x_ := x*cos + y*sin
		y_ := -x*sin + y*cos
		return s(x_, y_, z)
	}
}

// Rotates the shape around the 1-axis, over θ radians.
func (s Shape) RotY(θ float64) Shape {
	cos := math.Cos(θ)
	sin := math.Sin(θ)
	return func(x, y, z float64) bool {
		x_ := x*cos - z*sin
		z_ := x*sin + z*cos
		return s(x_, y, z_)
	}
}

// Rotates the shape around the 0-axis, over θ radians.
func (s Shape) RotX(θ float64) Shape {
	cos := math.Cos(θ)
	sin := math.Sin(θ)
	return func(x, y, z float64) bool {
		y_ := y*cos + z*sin
		z_ := -y*sin + z*cos
		return s(x, y_, z_)
	}
}

// Union of shapes a and b (logical OR).
func (a Shape) Add(b Shape) Shape {
	return func(x, y, z float64) bool {
		return a(x, y, z) || b(x, y, z)
	}
}

// Intersection of shapes a and b (logical AND).
func (a Shape) Intersect(b Shape) Shape {
	return func(x, y, z float64) bool {
		return a(x, y, z) && b(x, y, z)
	}
}

// Inverse (outside) of shape (logical NOT).
func (s Shape) Inverse() Shape {
	return func(x, y, z float64) bool {
		return !s(x, y, z)
	}
}

// Removes b from a (logical a AND NOT b)
func (a Shape) Sub(b Shape) Shape {
	return func(x, y, z float64) bool {
		return a(x, y, z) && !b(x, y, z)
	}
}

// Logical XOR of shapes a and b
func (a Shape) Xor(b Shape) Shape {
	return func(x, y, z float64) bool {
		A, B := a(x, y, z), b(x, y, z)
		return (A || B) && !(A && B)
	}
}
