package mag_config

// Utilities for setting magnetic configurations.

import (
	"math"
	"math/rand"

	"github.com/MathieuMoalic/amumax/src/mesh"
	"github.com/MathieuMoalic/amumax/src/utils"
	"github.com/MathieuMoalic/amumax/src/vector"
)

// Magnetic configuration returns m vector for position (x,y,z)
type Config func(x, y, z float64) vector.Vector

type ConfigList struct {
	mesh *mesh.Mesh
}

func (c *ConfigList) Init(mesh *mesh.Mesh) {
	c.mesh = mesh
}

// Random initial magnetization.
func (c *ConfigList) RandomMag() Config {
	return c.RandomMagSeed(0)
}

// Random initial magnetization,
// generated from random seed.
func (c *ConfigList) RandomMagSeed(seed int) Config {
	rng := rand.New(rand.NewSource(int64(seed)))
	return func(x, y, z float64) vector.Vector {
		return c.RandomDir(rng)
	}
}

// generate anisotropic random unit vector
func (c *ConfigList) RandomDir(rng *rand.Rand) vector.Vector {
	theta := 2 * rng.Float64() * math.Pi
	z := 2 * (rng.Float64() - 0.5)
	b := math.Sqrt(1 - z*z)
	x := b * math.Cos(theta)
	y := b * math.Sin(theta)
	return vector.Vector{x, y, z}
}

// Returns a Uniform magnetization state. E.g.:
//
//	M = Uniform(1, 0, 0)) // saturated along 0
func (c *ConfigList) Uniform(mx, my, mz float64) Config {
	return func(x, y, z float64) vector.Vector {
		return vector.Vector{mx, my, mz}
	}
}

// Make a Vortex magnetization with given circulation and core polarization (+1 or -1).
// The core is smoothed over a few exchange lengths and should easily relax to its ground state.
func (c *ConfigList) Vortex(circ, pol int) Config {
	diam2 := 2 * utils.Sqr64(c.mesh.CellSize()[0])
	return func(x, y, z float64) vector.Vector {
		r2 := x*x + y*y
		r := math.Sqrt(r2)
		mx := -y * float64(circ) / r
		my := x * float64(circ) / r
		mz := 1.5 * float64(pol) * math.Exp(-r2/diam2)
		return c.noNaN(vector.Vector{mx, my, mz}, pol)
	}
}

func (c *ConfigList) NeelSkyrmion(charge, pol int) Config {
	w := 8 * c.mesh.CellSize()[0]
	w2 := w * w
	return func(x, y, z float64) vector.Vector {
		r2 := x*x + y*y
		r := math.Sqrt(r2)
		mz := 2 * float64(pol) * (math.Exp(-r2/w2) - 0.5)
		mx := (x * float64(charge) / r) * (1 - math.Abs(mz))
		my := (y * float64(charge) / r) * (1 - math.Abs(mz))
		return c.noNaN(vector.Vector{mx, my, mz}, pol)
	}
}

func (c *ConfigList) BlochSkyrmion(charge, pol int) Config {
	w := 8 * c.mesh.CellSize()[0]
	w2 := w * w
	return func(x, y, z float64) vector.Vector {
		r2 := x*x + y*y
		r := math.Sqrt(r2)
		mz := 2 * float64(pol) * (math.Exp(-r2/w2) - 0.5)
		mx := (-y * float64(charge) / r) * (1 - math.Abs(mz))
		my := (x * float64(charge) / r) * (1 - math.Abs(mz))
		return c.noNaN(vector.Vector{mx, my, mz}, pol)
	}
}

func (c *ConfigList) AntiVortex(circ, pol int) Config {
	diam2 := 2 * utils.Sqr64(c.mesh.CellSize()[0])
	return func(x, y, z float64) vector.Vector {
		r2 := x*x + y*y
		r := math.Sqrt(r2)
		mx := -x * float64(circ) / r
		my := y * float64(circ) / r
		mz := 1.5 * float64(pol) * math.Exp(-r2/diam2)
		return c.noNaN(vector.Vector{mx, my, mz}, pol)
	}
}

func (c *ConfigList) Radial(charge, pol int) Config {
	return func(x, y, z float64) vector.Vector {
		r2 := x*x + y*y
		r := math.Sqrt(r2)
		mz := 0.0
		mx := (x * float64(charge) / r)
		my := (y * float64(charge) / r)
		return c.noNaN(vector.Vector{mx, my, mz}, pol)
	}
}

// Make a vortex wall configuration.
func (c *ConfigList) VortexWall(mleft, mright float64, circ, pol int) Config {
	h := c.mesh.WorldSize()[1]
	v := c.Vortex(circ, pol)
	return func(x, y, z float64) vector.Vector {
		if x < -h/2 {
			return vector.Vector{mleft, 0, 0}
		}
		if x > h/2 {
			return vector.Vector{mright, 0, 0}
		}
		return v(x, y, z)
	}
}

func (c *ConfigList) noNaN(v vector.Vector, pol int) vector.Vector {
	if math.IsNaN(v[0]) || math.IsNaN(v[1]) || math.IsNaN(v[2]) {
		return vector.Vector{0, 0, float64(pol)}
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
//	TwoDomain(1,0,0,  0,1,0,  -1,0,0) // head-to-head domains with transverse (Néel) wall
//	TwoDomain(1,0,0,  0,0,1,  -1,0,0) // head-to-head domains with perpendicular (Bloch) wall
//	TwoDomain(0,0,1,  1,0,0,   0,0,-1)// up-down domains with Bloch wall
func (c *ConfigList) TwoDomain(mx1, my1, mz1, mxwall, mywall, mzwall, mx2, my2, mz2 float64) Config {
	ww := 2 * c.mesh.CellSize()[0] // wall width in cells
	return func(x, y, z float64) vector.Vector {
		var m vector.Vector
		if x < 0 {
			m = vector.Vector{mx1, my1, mz1}
		} else {
			m = vector.Vector{mx2, my2, mz2}
		}
		gauss := math.Exp(-utils.Sqr64(x / ww))
		m[0] = (1-gauss)*m[0] + gauss*mxwall
		m[1] = (1-gauss)*m[1] + gauss*mywall
		m[2] = (1-gauss)*m[2] + gauss*mzwall
		return m
	}
}

// Conical magnetization configuration.
// The magnetization rotates on a cone defined by coneAngle and coneDirection.
// q is the wave vector of the Conical magnetization configuration.
// The magnetization is
//
//	m = u*cos(coneAngle) + sin(coneAngle)*( ua*cos(q*r) + ub*sin(q*r) )
//
// with ua and ub unit vectors perpendicular to u (normalized coneDirection)
func (c *ConfigList) Conical(q, coneDirection vector.Vector, coneAngle float64) Config {
	u := coneDirection.Div(coneDirection.Len())
	// two unit vectors perpendicular to each other and to the cone direction u
	p := math.Sqrt(1 - u[2]*u[2])
	ua := vector.Vector{u[0] * u[2], u[1] * u[2], u[2]*u[2] - 1}.Div(p)
	ub := vector.Vector{-u[1], u[0], 0}.Div(p)
	// cone direction along z direction? -> oops divided by zero, let's fix this
	if u[2]*u[2] == 1 {
		ua = vector.Vector{1, 0, 0}
		ub = vector.Vector{0, 1, 0}
	}
	sina, cosa := math.Sincos(coneAngle)
	return func(x, y, z float64) vector.Vector {
		sinqr, cosqr := math.Sincos(q[0]*x + q[1]*y + q[2]*z)
		return u.Mul(cosa).MAdd(sina*cosqr, ua).MAdd(sina*sinqr, ub)
	}
}

func (c *ConfigList) Helical(q vector.Vector) Config {
	return c.Conical(q, q, math.Pi/2)
}

// Transl returns a translated copy of configuration c. E.g.:
//
//	M = Vortex(1, 1).Transl(100e-9, 0, 0)  // vortex with center at x=100nm
func (c Config) Transl(dx, dy, dz float64) Config {
	return func(x, y, z float64) vector.Vector {
		return c(x-dx, y-dy, z-dz)
	}
}

// Scale returns a scaled copy of configuration c.
func (c Config) Scale(sx, sy, sz float64) Config {
	return func(x, y, z float64) vector.Vector {
		return c(x/sx, y/sy, z/sz)
	}
}

// Rotates the configuration around the 2-axis, over θ radians.
func (c Config) RotZ(θ float64) Config {
	cos := math.Cos(θ)
	sin := math.Sin(θ)
	return func(x, y, z float64) vector.Vector {
		x_ := x*cos + y*sin
		y_ := -x*sin + y*cos
		m := c(x_, y_, z)
		mx_ := m[0]*cos - m[1]*sin
		my_ := m[0]*sin + m[1]*cos
		return vector.Vector{mx_, my_, m[2]}
	}
}

// Returns a new magnetization equal to c + weight * other.
// E.g.:
//
//	Uniform(1, 0, 0).Add(0.2, RandomMag())
//
// for a uniform state with 20% random distortion.
func (c Config) Add(weight float64, other Config) Config {
	return func(x, y, z float64) vector.Vector {
		return c(x, y, z).MAdd(weight, other(x, y, z))
	}
}
