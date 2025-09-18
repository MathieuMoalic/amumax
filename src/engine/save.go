package engine

import (
	"fmt"
	"path"
	"reflect"
	"strings"

	"github.com/MathieuMoalic/amumax/src/cuda"
	"github.com/MathieuMoalic/amumax/src/data"
	"github.com/MathieuMoalic/amumax/src/draw"
	"github.com/MathieuMoalic/amumax/src/fsutil"
	"github.com/MathieuMoalic/amumax/src/log"
	"github.com/MathieuMoalic/amumax/src/oommf"
)

func init() {
	declROnly("OVF1_BINARY", OVF1Binary, "OutputFormat = OVF1_BINARY sets binary OVF1 output")
	declROnly("OVF2_BINARY", OVF2Binary, "OutputFormat = OVF2_BINARY sets binary OVF2 output")
	declROnly("OVF1_TEXT", OVF1Text, "OutputFormat = OVF1_TEXT sets text OVF1 output")
	declROnly("OVF2_TEXT", OVF2Text, "OutputFormat = OVF2_TEXT sets text OVF2 output")
	declROnly("DUMP", DUMP, "OutputFormat = DUMP sets text DUMP output")
}

var (
	filenameFormat = "%s%06d"   // formatting string for auto filenames.
	snapshotFormat = "jpg"      // user-settable snapshot format
	outputFormat   = OVF2Binary // user-settable output format
)

type fformat struct{}

func (*fformat) Eval() any          { return filenameFormat }
func (*fformat) SetValue(v any)     { drainOutput(); filenameFormat = v.(string) }
func (*fformat) Type() reflect.Type { return reflect.TypeOf("") }

type oformat struct{}

func (*oformat) Eval() any          { return outputFormat }
func (*oformat) SetValue(v any)     { drainOutput(); outputFormat = v.(outputFormatType) }
func (*oformat) Type() reflect.Type { return reflect.TypeOf(outputFormatType(OVF2Binary)) }

// saveOVF once, with auto file name
func saveOVF(q Quantity) {
	qname := nameOf(q)
	fname := autoFname(nameOf(q), outputFormat, autonum[qname])
	saveAsOVF(q, fname)
	autonum[qname]++
}

// Save under given file name (transparent async I/O).
func saveAsOVF(q Quantity, fname string) {
	if !strings.HasPrefix(fname, OD()) {
		fname = OD() + fname // don't clean, turns http:// in http:/
	}

	if path.Ext(fname) == "" {
		fname += ("." + StringFromOutputFormat[outputFormat])
	}
	buffer := ValueOf(q) // TODO: check and optimize for Buffer()
	defer cuda.Recycle(buffer)
	info := oommf.Meta{Time: Time, Name: nameOf(q), Unit: unitOf(q), CellSize: MeshOf(q).CellSize()}
	data := buffer.HostCopy() // must be copy (async io)
	queOutput(func() { saveAsSync(fname, data, info, outputFormat) })
}

// Save image once, with auto file name
func snapshot(q Quantity) {
	qname := nameOf(q)
	fname := fmt.Sprintf(OD()+filenameFormat+"."+snapshotFormat, qname, autonum[qname])
	s := ValueOf(q)
	defer cuda.Recycle(s)
	data := s.HostCopy() // must be copy (asyncio)
	queOutput(func() { snapshotSync(fname, data) })
	autonum[qname]++
}

func snapshotAs(q Quantity, fname string) {
	if !strings.HasPrefix(fname, OD()) {
		fname = OD() + fname // don't clean, turns http:// in http:/
	}

	if path.Ext(fname) == "" {
		fname += ("." + StringFromOutputFormat[outputFormat])
	}
	s := ValueOf(q)
	defer cuda.Recycle(s)
	data := s.HostCopy() // must be copy (asyncio)
	queOutput(func() { snapshotSync(fname, data) })
}

// synchronous snapshot
func snapshotSync(fname string, output *data.Slice) {
	f, err := fsutil.Create(fname)
	log.Log.PanicIfError(err)
	defer f.Close()
	arrowSize := 16
	err = draw.RenderFormat(f, output, "auto", "auto", arrowSize, path.Ext(fname))
	if err != nil {
		log.Log.Warn("Error while rendering snapshot: %v", err)
	}
}

// synchronous save
func saveAsSync(fname string, s *data.Slice, info oommf.Meta, format outputFormatType) {
	f, err := fsutil.Create(fname)
	log.Log.PanicIfError(err)
	defer f.Close()

	switch format {
	case OVF1Text:
		oommf.WriteOVF1(f, s, info, "text")
	case OVF1Binary:
		oommf.WriteOVF1(f, s, info, "binary 4")
	case OVF2Text:
		oommf.WriteOVF2(f, s, info, "text")
	case OVF2Binary:
		oommf.WriteOVF2(f, s, info, "binary 4")
	default:
		panic("invalid output format")
	}
}

type outputFormatType int

const (
	OVF1Text outputFormatType = iota + 1
	OVF1Binary
	OVF2Text
	OVF2Binary
	DUMP
)

var StringFromOutputFormat = map[outputFormatType]string{
	OVF1Text:   "ovf",
	OVF1Binary: "ovf",
	OVF2Text:   "ovf",
	OVF2Binary: "ovf",
	DUMP:       "dump",
}
