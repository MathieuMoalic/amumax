package draw

import (
	"bufio"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"path"
	"strings"

	"github.com/MathieuMoalic/amumax/src/data"
)

func RenderFormat(out io.Writer, f *data.Slice, min, max string, arrowSize int, format string, colormap ...ColorMapSpec) error {
	codecs := map[string]codec{".png": pngfull, ".jpg": jpeg100, ".gif": gif256}
	ext := strings.ToLower(path.Ext(format))
	enc := codecs[ext]
	if enc == nil {
		return fmt.Errorf("render: unhandled image type: %s", ext)
	}
	return render(out, f, min, max, arrowSize, enc, colormap...)
}

// encodes an image
type codec func(io.Writer, image.Image) error

// render data and encode with arbitrary codec.
func render(out io.Writer, f *data.Slice, min, max string, arrowSize int, encode codec, colormap ...ColorMapSpec) error {
	img := createRGBAImage(f, min, max, arrowSize, colormap...)
	buf := bufio.NewWriter(out)
	defer buf.Flush()
	return encode(buf, img)
}

// full-quality jpeg codec, passable to Render()
func jpeg100(w io.Writer, img image.Image) error {
	return jpeg.Encode(w, img, &jpeg.Options{Quality: 100})
}

// full quality gif coded, passable to Render()
func gif256(w io.Writer, img image.Image) error {
	return gif.Encode(w, img, &gif.Options{NumColors: 256, Quantizer: nil, Drawer: nil})
}

// png codec, passable to Render()
func pngfull(w io.Writer, img image.Image) error {
	return png.Encode(w, img)
}
