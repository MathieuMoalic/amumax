package draw

import (
	"image"
	"image/color"
	"math"

	"github.com/golang/freetype/raster"
	"golang.org/x/image/math/fixed"

	"github.com/MathieuMoalic/amumax/src/data"
)

func drawArrows(img *image.RGBA, arr [3][][][]float32, sub int) {
	c := newCanvas(img)

	Na := data.SizeOf(arr[0]) // number of arrows
	h := Na[Y]                // orignal image height
	Na[X] = imax(Na[X]/sub, 1)
	Na[Y] = imax(Na[Y]/sub, 1)
	Na[Z] = 1
	small := data.Downsample(arr[:], Na)
	S := float32(sub)

	for iy := 0; iy < Na[Y]; iy++ {
		Ay := float32(h) - (float32(iy)+0.5)*S
		for ix := 0; ix < Na[X]; ix++ {
			Ax := (float32(ix) + 0.5) * S
			mx := small[X][0][iy][ix]
			my := small[Y][0][iy][ix]
			mz := small[Z][0][iy][ix]
			c.Arrow(Ax, Ay, mx, my, mz, float32(sub))

		}
	}

	c.rasterizer.Rasterize(c.RGBAPainter)
	c.rasterizer.Clear()
}

// A canvas is used to draw on.
type canvas struct {
	*image.RGBA
	*raster.RGBAPainter
	rasterizer *raster.Rasterizer
}

// Make a new canvas of size w x h.
func newCanvas(img *image.RGBA) *canvas {
	c := new(canvas)
	c.RGBA = img
	c.RGBAPainter = raster.NewRGBAPainter(c.RGBA)
	c.rasterizer = raster.NewRasterizer(img.Bounds().Max.X, img.Bounds().Max.Y)
	c.rasterizer.UseNonZeroWinding = true
	c.SetColor(color.RGBA{0, 0, 0, 100})
	return c
}

func (c *canvas) Arrow(x, y, mx, my, mz, size float32) {
	arrlen := 0.4 * size
	arrw := 0.2 * size

	norm := float32(math.Sqrt(float64(mx*mx + my*my + mz*mz)))
	if norm == 0 {
		return
	}
	if norm > 1 {
		norm = 1
	}

	theta := math.Atan2(float64(my), float64(mx))
	cos := float32(math.Cos(theta))
	sin := float32(math.Sin(theta))
	r1 := arrlen * norm * float32(math.Cos(math.Asin(float64(mz))))
	r2 := arrw * norm

	pt1 := fixed.Point26_6{X: fixed.Int26_6((r1 * cos) + x), Y: fixed.Int26_6(-(r1 * sin) + y)}
	pt2 := fixed.Point26_6{X: fixed.Int26_6((r2*sin - r1*cos) + x), Y: fixed.Int26_6(-(-r2*cos - r1*sin) + y)}
	pt3 := fixed.Point26_6{X: fixed.Int26_6((-r2*sin - r1*cos) + x), Y: fixed.Int26_6(-(r2*cos - r1*sin) + y)}
	var path raster.Path
	path.Start(pt1)
	path.Add1(pt2)
	path.Add1(pt3)
	path.Add1(pt1)

	c.rasterizer.AddPath(path)
}

func imax(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}
