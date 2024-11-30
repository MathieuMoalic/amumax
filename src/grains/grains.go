package grains

import (
	"math"
	"math/rand"
)

type addVoronoiToRegionsType func(minRegion int, maxRegion int, getRegion func(float64, float64, float64) int)

type Grains struct {
	addVoronoiToRegions addVoronoiToRegionsType

	grainsize      float64
	tilesize       float64
	tile           int
	minRegion      int
	maxRegion      int
	cache          map[int2][]center
	seed           int64
	rnd            *rand.Rand
	poisson_lambda float64
}

// we pass the registerFunction and addVoronoiToRegions functions to the constructor
// to avoid circular dependencies
func NewGrains(registerFunction func(string, interface{}), addVoronoiToRegions addVoronoiToRegionsType) *Grains {
	TILE := 2 // tile size in grains
	g := &Grains{
		addVoronoiToRegions: addVoronoiToRegions,
		tile:                TILE,
		cache:               make(map[int2][]center),
		rnd:                 rand.New(rand.NewSource(0)),
		poisson_lambda:      float64(TILE * TILE),
	}
	registerFunction("ext_makegrains", g.Voronoi)

	return g
}

// script function
func (g *Grains) Voronoi(grainsize float64, minRegion, maxRegion, seed int) {
	g.grainsize = grainsize
	g.minRegion = minRegion
	g.maxRegion = maxRegion
	g.seed = int64(seed)
	g.tilesize = grainsize * float64(g.tile) // expect 4 grains/block, 36 per 3x3 blocks = safe, relatively round number
	// put this code in the region struct
	g.addVoronoiToRegions(minRegion, maxRegion, g.GetRegion)
}

// integer tile coordinate
type int2 struct{ x, y int }

// Voronoi center info
type center struct {
	x, y   float64 // center position (m)
	region byte    // region for all cells near center
}

// Returns the region of the grain where cell at x,y,z belongs to
func (t *Grains) GetRegion(x, y, z float64) int {
	tile := t.tileOf(x, y) // tile containing x,y

	// look for nearest center in tile + neighbors
	nearest := center{} // dummy initial value, but safe should the infinite impossibility strike.
	mindist := math.Inf(1)
	for tx := tile.x - 1; tx <= tile.x+1; tx++ {
		for ty := tile.y - 1; ty <= tile.y+1; ty++ {
			centers := t.centersInTile(tx, ty)
			for _, c := range centers {
				dist := (x-c.x)*(x-c.x) + (y-c.y)*(y-c.y)
				if dist < mindist {
					nearest = c
					mindist = dist
				}
			}
		}
	}

	//fmt.Println("nearest", x, y, ":", nearest)
	return int(nearest.region)
}

// Returns the list of Voronoi centers in tile(ix, iy), using only ix,iy to seed the random generator
func (t *Grains) centersInTile(tx, ty int) []center {
	pos := int2{tx, ty}
	if c, ok := t.cache[pos]; ok {
		return c
	} else {
		// tile-specific seed that works for positive and negative tx, ty
		seed := (int64(ty)+(1<<24))*(1<<24) + (int64(tx) + (1 << 24))
		t.rnd.Seed(seed ^ t.seed)
		N := t.poisson()
		c := make([]center, N)

		// absolute position of tile (m)
		x0, y0 := float64(tx)*t.tilesize, float64(ty)*t.tilesize

		for i := range c {
			// random position inside tile
			c[i].x = x0 + t.rnd.Float64()*t.tilesize
			c[i].y = y0 + t.rnd.Float64()*t.tilesize
			c[i].region = byte(t.rnd.Intn(t.maxRegion-t.minRegion) + t.minRegion)
		}
		t.cache[pos] = c
		return c
	}
}

func (t *Grains) tileOf(x, y float64) int2 {
	ix := int(math.Floor(x / t.tilesize))
	iy := int(math.Floor(y / t.tilesize))
	return int2{ix, iy}
}

// Generate poisson distributed numbers (according to Knuth)
func (t *Grains) poisson() int {
	L := math.Exp(-t.poisson_lambda)
	k := 1
	p := t.rnd.Float64()
	for p > L {
		k++
		p *= t.rnd.Float64()
	}
	return k - 1
}
