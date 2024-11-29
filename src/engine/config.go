package engine

// Utilities for setting magnetic configurations.

import (
	"math"
	"math/rand"

	"github.com/MathieuMoalic/amumax/src/data"
	"github.com/MathieuMoalic/amumax/src/mesh"
	"github.com/MathieuMoalic/amumax/src/script"
	"github.com/MathieuMoalic/amumax/src/utils"
)

// Magnetic configuration returns m vector for position (x,y,z)
type config func(x, y, z float64) data.Vector

type configList struct {
	mesh *mesh.Mesh
}

func newConfigList(mesh *mesh.Mesh, s *script.ScriptParser) *configList {
	c := &configList{mesh: mesh}
	s.RegisterFunction("RandomMag", c.randomMag)
	s.RegisterFunction("RandomMagSeed", c.randomMagSeed)
	s.RegisterFunction("Uniform", c.uniform)
	s.RegisterFunction("Vortex", c.vortex)
	s.RegisterFunction("NeelSkyrmion", c.neelSkyrmion)
	s.RegisterFunction("BlochSkyrmion", c.blochSkyrmion)
	s.RegisterFunction("AntiVortex", c.antiVortex)
	s.RegisterFunction("Radial", c.radial)
	s.RegisterFunction("VortexWall", c.vortexWall)
	s.RegisterFunction("TwoDomain", c.twoDomain)
	s.RegisterFunction("Conical", c.conical)
	s.RegisterFunction("Helical", c.helical)
	return c
}

// Random initial magnetization.
func (c *configList) randomMag() config {
	return c.randomMagSeed(0)
}

// Random initial magnetization,
// generated from random seed.
func (c *configList) randomMagSeed(seed int) config {
	rng := rand.New(rand.NewSource(int64(seed)))
	return func(x, y, z float64) data.Vector {
		return c.randomDir(rng)
	}
}

// generate anisotropic random unit vector
func (c *configList) randomDir(rng *rand.Rand) data.Vector {
	theta := 2 * rng.Float64() * math.Pi
	z := 2 * (rng.Float64() - 0.5)
	b := math.Sqrt(1 - z*z)
	x := b * math.Cos(theta)
	y := b * math.Sin(theta)
	return data.Vector{x, y, z}
}

// Returns a uniform magnetization state. E.g.:
//
//	M = uniform(1, 0, 0)) // saturated along X
func (c *configList) uniform(mx, my, mz float64) config {
	return func(x, y, z float64) data.Vector {
		return data.Vector{mx, my, mz}
	}
}

// Make a vortex magnetization with given circulation and core polarization (+1 or -1).
// The core is smoothed over a few exchange lengths and should easily relax to its ground state.
func (c *configList) vortex(circ, pol int) config {
	diam2 := 2 * utils.Sqr64(c.mesh.CellSize()[X])
	return func(x, y, z float64) data.Vector {
		r2 := x*x + y*y
		r := math.Sqrt(r2)
		mx := -y * float64(circ) / r
		my := x * float64(circ) / r
		mz := 1.5 * float64(pol) * math.Exp(-r2/diam2)
		return c.noNaN(data.Vector{mx, my, mz}, pol)
	}
}

func (c *configList) neelSkyrmion(charge, pol int) config {
	w := 8 * c.mesh.CellSize()[X]
	w2 := w * w
	return func(x, y, z float64) data.Vector {
		r2 := x*x + y*y
		r := math.Sqrt(r2)
		mz := 2 * float64(pol) * (math.Exp(-r2/w2) - 0.5)
		mx := (x * float64(charge) / r) * (1 - math.Abs(mz))
		my := (y * float64(charge) / r) * (1 - math.Abs(mz))
		return c.noNaN(data.Vector{mx, my, mz}, pol)
	}
}

func (c *configList) blochSkyrmion(charge, pol int) config {
	w := 8 * c.mesh.CellSize()[X]
	w2 := w * w
	return func(x, y, z float64) data.Vector {
		r2 := x*x + y*y
		r := math.Sqrt(r2)
		mz := 2 * float64(pol) * (math.Exp(-r2/w2) - 0.5)
		mx := (-y * float64(charge) / r) * (1 - math.Abs(mz))
		my := (x * float64(charge) / r) * (1 - math.Abs(mz))
		return c.noNaN(data.Vector{mx, my, mz}, pol)
	}
}

func (c *configList) antiVortex(circ, pol int) config {
	diam2 := 2 * utils.Sqr64(c.mesh.CellSize()[X])
	return func(x, y, z float64) data.Vector {
		r2 := x*x + y*y
		r := math.Sqrt(r2)
		mx := -x * float64(circ) / r
		my := y * float64(circ) / r
		mz := 1.5 * float64(pol) * math.Exp(-r2/diam2)
		return c.noNaN(data.Vector{mx, my, mz}, pol)
	}
}

func (c *configList) radial(charge, pol int) config {
	return func(x, y, z float64) data.Vector {
		r2 := x*x + y*y
		r := math.Sqrt(r2)
		mz := 0.0
		mx := (x * float64(charge) / r)
		my := (y * float64(charge) / r)
		return c.noNaN(data.Vector{mx, my, mz}, pol)
	}
}

// Make a vortex wall configuration.
func (c *configList) vortexWall(mleft, mright float64, circ, pol int) config {
	h := c.mesh.WorldSize()[Y]
	v := c.vortex(circ, pol)
	return func(x, y, z float64) data.Vector {
		if x < -h/2 {
			return data.Vector{mleft, 0, 0}
		}
		if x > h/2 {
			return data.Vector{mright, 0, 0}
		}
		return v(x, y, z)
	}
}

func (c *configList) noNaN(v data.Vector, pol int) data.Vector {
	if math.IsNaN(v[X]) || math.IsNaN(v[Y]) || math.IsNaN(v[Z]) {
		return data.Vector{0, 0, float64(pol)}
	} else {
		return v
	}
}

// Make a 2-domain configuration with domain wall.
// (mx1, my1, mz1) and (mx2, my2, mz2) are the magnetizations in the left and right domain, respectively.
// (mxwall, mywall, mzwall) is the magnetization in the wall. The wall is smoothed over a few cells so it will
// easily relax to its ground state.
// E.g.:
//
//	twoDomain(1,0,0,  0,1,0,  -1,0,0) // head-to-head domains with transverse (Néel) wall
//	twoDomain(1,0,0,  0,0,1,  -1,0,0) // head-to-head domains with perpendicular (Bloch) wall
//	twoDomain(0,0,1,  1,0,0,   0,0,-1)// up-down domains with Bloch wall
func (c *configList) twoDomain(mx1, my1, mz1, mxwall, mywall, mzwall, mx2, my2, mz2 float64) config {
	ww := 2 * c.mesh.CellSize()[X] // wall width in cells
	return func(x, y, z float64) data.Vector {
		var m data.Vector
		if x < 0 {
			m = data.Vector{mx1, my1, mz1}
		} else {
			m = data.Vector{mx2, my2, mz2}
		}
		gauss := math.Exp(-utils.Sqr64(x / ww))
		m[X] = (1-gauss)*m[X] + gauss*mxwall
		m[Y] = (1-gauss)*m[Y] + gauss*mywall
		m[Z] = (1-gauss)*m[Z] + gauss*mzwall
		return m
	}
}

// conical magnetization configuration.
// The magnetization rotates on a cone defined by coneAngle and coneDirection.
// q is the wave vector of the conical magnetization configuration.
// The magnetization is
//
//	m = u*cos(coneAngle) + sin(coneAngle)*( ua*cos(q*r) + ub*sin(q*r) )
//
// with ua and ub unit vectors perpendicular to u (normalized coneDirection)
func (c *configList) conical(q, coneDirection data.Vector, coneAngle float64) config {
	u := coneDirection.Div(coneDirection.Len())
	// two unit vectors perpendicular to each other and to the cone direction u
	p := math.Sqrt(1 - u[Z]*u[Z])
	ua := data.Vector{u[X] * u[Z], u[Y] * u[Z], u[Z]*u[Z] - 1}.Div(p)
	ub := data.Vector{-u[Y], u[X], 0}.Div(p)
	// cone direction along z direction? -> oops divided by zero, let's fix this
	if u[Z]*u[Z] == 1 {
		ua = data.Vector{1, 0, 0}
		ub = data.Vector{0, 1, 0}
	}
	sina, cosa := math.Sincos(coneAngle)
	return func(x, y, z float64) data.Vector {
		sinqr, cosqr := math.Sincos(q[X]*x + q[Y]*y + q[Z]*z)
		return u.Mul(cosa).MAdd(sina*cosqr, ua).MAdd(sina*sinqr, ub)
	}
}

func (c *configList) helical(q data.Vector) config {
	return c.conical(q, q, math.Pi/2)
}

// Transl returns a translated copy of configuration c. E.g.:
//
//	M = Vortex(1, 1).Transl(100e-9, 0, 0)  // vortex with center at x=100nm
func (c config) Transl(dx, dy, dz float64) config {
	return func(x, y, z float64) data.Vector {
		return c(x-dx, y-dy, z-dz)
	}
}

// Scale returns a scaled copy of configuration c.
func (c config) Scale(sx, sy, sz float64) config {
	return func(x, y, z float64) data.Vector {
		return c(x/sx, y/sy, z/sz)
	}
}

// Rotates the configuration around the Z-axis, over θ radians.
func (c config) RotZ(θ float64) config {
	cos := math.Cos(θ)
	sin := math.Sin(θ)
	return func(x, y, z float64) data.Vector {
		x_ := x*cos + y*sin
		y_ := -x*sin + y*cos
		m := c(x_, y_, z)
		mx_ := m[X]*cos - m[Y]*sin
		my_ := m[X]*sin + m[Y]*cos
		return data.Vector{mx_, my_, m[Z]}
	}
}

// Returns a new magnetization equal to c + weight * other.
// E.g.:
//
//	Uniform(1, 0, 0).Add(0.2, RandomMag())
//
// for a uniform state with 20% random distortion.
func (c config) Add(weight float64, other config) config {
	return func(x, y, z float64) data.Vector {
		return c(x, y, z).MAdd(weight, other(x, y, z))
	}
}
