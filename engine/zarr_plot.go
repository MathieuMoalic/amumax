package engine

import (
	"bytes"
	"errors"
	"image/color"
	"io"
	"time"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

const DefaultCacheLifetime = 1 * time.Second

var Tableplot *zTablePlot

func init() {
	Tableplot = &zTablePlot{"t", "mx", ImagePlotCache{}}
	Tableplot.cache.expirytime = time.Now().Add(DefaultCacheLifetime)
}

type ImagePlotCache struct {
	img        []byte    // cached output
	err        error     // cached error
	expirytime time.Time // expiration time of the cache
	// lifetime   time.Duration // maximum lifetime of the cache
}

// zTablePlot is a WriterTo which writes a (cached) plot image of table data.
// The internally cached image will be updated when the column indices have been changed,
// or when the cache lifetime is exceeded.
// If another GO routine is updating the image, the cached image will be written.
type zTablePlot struct {
	// updating   bool
	X, Y  string
	cache ImagePlotCache
}

// func zNewPlot(table *zTablePlot) (p *zTablePlot) {
// 	// p = &zTablePlot{table: table, xcol: "", ycol: ""}
// 	// p.cache.lifetime = DefaultCacheLifetime
// 	return
// }

func (p *zTablePlot) NeedSave() bool {
	return time.Now().After(p.cache.expirytime)
}

func (p *zTablePlot) SelectDataColumns(xlabel, ylabel string) error {
	if _, exists := ZTables.Data[xlabel]; exists {
		p.X = xlabel
	} else {
		return errors.New("ylabel doesn't exist")
	}
	if _, exists := ZTables.Data[ylabel]; exists {
		p.Y = ylabel
	} else {
		return errors.New("ylabel doesn't exist")
	}
	return nil
}

func (p *zTablePlot) WriteTo(w io.Writer) (int64, error) {
	p.update()

	if p.cache.err != nil {
		return 0, p.cache.err
	}
	nBytes, err := w.Write(p.cache.img)
	return int64(nBytes), err
}

// Updates the cached image if the cache is expired
// Does nothing if the image is already being updated by another GO process
func (p *zTablePlot) update() {
	xdata := ZTables.Data[p.X]
	ydata := ZTables.Data[p.Y][:len(xdata)]
	points := make(plotter.XYs, len(xdata))
	for i := 0; i < len(xdata); i++ {
		points[i].X = xdata[i]
		points[i].Y = ydata[i]
	}

	line, err := plotter.NewLine(points)
	line.Color = color.RGBA{248, 248, 242, 255}
	if err != nil {
		return
	}
	pl := plot.New()
	pl.X.Label.Text = p.X
	pl.Y.Label.Text = p.Y
	pl.X.Tick.Color = color.RGBA{248, 248, 242, 255}
	pl.Y.Tick.Color = color.RGBA{248, 248, 242, 255}
	pl.X.Tick.Label.Color = color.RGBA{248, 248, 242, 255}
	pl.Y.Tick.Label.Color = color.RGBA{248, 248, 242, 255}
	pl.X.Label.TextStyle.Color = color.RGBA{248, 248, 242, 255}
	pl.Y.Label.TextStyle.Color = color.RGBA{248, 248, 242, 255}
	pl.X.Color = color.RGBA{248, 248, 242, 255}
	pl.Y.Color = color.RGBA{248, 248, 242, 255}
	pl.BackgroundColor = color.RGBA{40, 42, 54, 0}
	pl.Add(line)
	wr, err := pl.WriterTo(8*vg.Inch, 4*vg.Inch, "png")
	if err != nil {
		return
	}

	buf := bytes.NewBuffer(nil)
	_, err = wr.WriteTo(buf)

	if err != nil {
		LogOut("Warning: Couldn't render the plot")
		p.cache.err = err
	} else {
		p.cache.img = buf.Bytes()
		p.cache.expirytime = time.Now().Add(DefaultCacheLifetime)
	}
}

func (p *zTablePlot) Render() (*bytes.Buffer, error) {

	xdata := ZTables.Data[p.X]
	ydata := ZTables.Data[p.Y][:len(xdata)]
	points := make(plotter.XYs, len(xdata))
	for i := 0; i < len(xdata); i++ {
		points[i].X = xdata[i]
		points[i].Y = ydata[i]
	}

	line, err := plotter.NewLine(points)
	line.Color = color.RGBA{248, 248, 242, 255}
	if err != nil {
		return nil, err
	}
	pl := plot.New()
	pl.X.Label.Text = p.X
	pl.Y.Label.Text = p.Y
	pl.X.Tick.Color = color.RGBA{248, 248, 242, 255}
	pl.Y.Tick.Color = color.RGBA{248, 248, 242, 255}
	pl.X.Tick.Label.Color = color.RGBA{248, 248, 242, 255}
	pl.Y.Tick.Label.Color = color.RGBA{248, 248, 242, 255}
	pl.X.Label.TextStyle.Color = color.RGBA{248, 248, 242, 255}
	pl.Y.Label.TextStyle.Color = color.RGBA{248, 248, 242, 255}
	pl.X.Color = color.RGBA{248, 248, 242, 255}
	pl.Y.Color = color.RGBA{248, 248, 242, 255}
	pl.BackgroundColor = color.RGBA{40, 42, 54, 0}
	pl.Add(line)
	wr, err := pl.WriterTo(8*vg.Inch, 4*vg.Inch, "png")
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(nil)
	_, err = wr.WriteTo(buf)
	return buf, err
}
