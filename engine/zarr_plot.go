package engine

// import (
// 	"bytes"
// 	"io"
// 	"time"

// 	"gonum.org/v1/plot"
// 	"gonum.org/v1/plot/plotter"
// 	"gonum.org/v1/plot/plotutil"
// 	"gonum.org/v1/plot/vg"
// )

// const DefaultCacheLifetime = 1 * time.Second

// var tableplot *zTablePlot

// // zTablePlot is a WriterTo which writes a (cached) plot image of table data.
// // The internally cached image will be updated when the column indices have been changed,
// // or when the cache lifetime is exceeded.
// // If another GO routine is updating the image, the cached image will be written.
// type zTablePlot struct {
// 	updating   bool
// 	table      *DataTable
// 	xcol, ycol string

// 	cache struct {
// 		img        []byte        // cached output
// 		err        error         // cached error
// 		expirytime time.Time     // expiration time of the cache
// 		lifetime   time.Duration // maximum lifetime of the cache
// 	}
// }

// func zNewPlot(table *DataTable) (p *zTablePlot) {
// 	p = &zTablePlot{table: table, xcol: "", ycol: ""}
// 	p.cache.lifetime = DefaultCacheLifetime
// 	return
// }

// func (p *zTablePlot) SelectDataColumns(xcolidx, ycolidx string) {
// 	if xcolidx != p.xcol || ycolidx != p.ycol {
// 		p.xcol, p.ycol = xcolidx, ycolidx
// 		p.cache.expirytime = time.Time{} // this will trigger an update at the next write
// 	}
// }

// func (p *zTablePlot) WriteTo(w io.Writer) (int64, error) {
// 	p.update()

// 	if p.cache.err != nil {
// 		return 0, p.cache.err
// 	}
// 	nBytes, err := w.Write(p.cache.img)
// 	return int64(nBytes), err
// }

// // Updates the cached image if the cache is expired
// // Does nothing if the image is already being updated by another GO process
// func (p *zTablePlot) update() {
// 	// xcol, ycol := p.xcol, p.ycol
// 	// needupdate := !p.updating && time.Now().After(p.cache.expirytime)
// 	// p.updating = p.updating || needupdate

// 	// if !needupdate {
// 	// 	return
// 	// }
// 	// img, err := zCreatePlot(p.table, xcol, ycol)

// 	// p.cache.img, p.cache.err = img, err
// 	// p.updating = false
// 	// if p.xcol == xcol && p.ycol == ycol {
// 	// 	p.cache.expirytime = time.Now().Add(p.cache.lifetime)
// 	// } else { // column indices have been changed during the update
// 	// 	p.cache.expirytime = time.Time{}
// 	// }
// }

// // Returns a png image plot of table data
// func zCreatePlot(table *DataTable, xcol, ycol string) (img []byte, err error) {
// 	// xdata := ZTableTime.Data
// 	// ydata := ZTables[ycol].Data[0]
// 	xdata := []float64{0}
// 	ydata := []float64{0}
// 	if (len(xdata) == 0) || (len(ydata) == 0) {
// 		xdata = []float64{0}
// 		ydata = []float64{0}
// 	}
// 	pl := plot.New()
// 	pl.X.Label.Text = xcol
// 	pl.Y.Label.Text = ycol
// 	pl.X.Label.Padding = 0.2 * vg.Inch
// 	pl.Y.Label.Padding = 0.2 * vg.Inch

// 	points := make(plotter.XYs, len(xdata))
// 	for i := 0; i < len(xdata); i++ {
// 		points[i].X = xdata[i]
// 		points[i].Y = ydata[i]
// 	}
// 	err = plotutil.AddLinePoints(pl, points, "legend1")
// 	if err != nil {
// 		panic(err)
// 	}

// 	wr, err := pl.WriterTo(8*vg.Inch, 4*vg.Inch, "png")
// 	if err != nil {
// 		return
// 	}

// 	buf := bytes.NewBuffer(nil)
// 	_, err = wr.WriteTo(buf)

// 	if err != nil {
// 		return nil, err
// 	} else {
// 		return buf.Bytes(), nil
// 	}
// }
