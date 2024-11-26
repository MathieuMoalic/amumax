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
	declROnly("OVF1_BINARY", OVF1_BINARY, "OutputFormat = OVF1_BINARY sets binary OVF1 output")
	declROnly("OVF2_BINARY", OVF2_BINARY, "OutputFormat = OVF2_BINARY sets binary OVF2 output")
	declROnly("OVF1_TEXT", OVF1_TEXT, "OutputFormat = OVF1_TEXT sets text OVF1 output")
	declROnly("OVF2_TEXT", OVF2_TEXT, "OutputFormat = OVF2_TEXT sets text OVF2 output")
	declROnly("DUMP", DUMP, "OutputFormat = DUMP sets text DUMP output")
}

var (
	filenameFormat = "%s%06d"    // formatting string for auto filenames.
	snapshotFormat = "jpg"       // user-settable snapshot format
	outputFormat   = OVF2_BINARY // user-settable output format
)

type fformat struct{}

func (*fformat) Eval() interface{}      { return filenameFormat }
func (*fformat) SetValue(v interface{}) { drainOutput(); filenameFormat = v.(string) }
func (*fformat) Type() reflect.Type     { return reflect.TypeOf("") }

type oformat struct{}

func (*oformat) Eval() interface{}      { return outputFormat }
func (*oformat) SetValue(v interface{}) { drainOutput(); outputFormat = v.(outputFormatType) }
func (*oformat) Type() reflect.Type     { return reflect.TypeOf(outputFormatType(OVF2_BINARY)) }

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
	queOutput(func() { saveAs_sync(fname, data, info, outputFormat) })
}

// Save image once, with auto file name
func snapshot(q Quantity) {
	qname := nameOf(q)
	fname := fmt.Sprintf(OD()+filenameFormat+"."+snapshotFormat, qname, autonum[qname])
	s := ValueOf(q)
	defer cuda.Recycle(s)
	data := s.HostCopy() // must be copy (asyncio)
	queOutput(func() { snapshot_sync(fname, data) })
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
	queOutput(func() { snapshot_sync(fname, data) })
}

// synchronous snapshot
func snapshot_sync(fname string, output *data.Slice) {
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
func saveAs_sync(fname string, s *data.Slice, info oommf.Meta, format outputFormatType) {
	f, err := fsutil.Create(fname)
	log.Log.PanicIfError(err)
	defer f.Close()

	switch format {
	case OVF1_TEXT:
		oommf.WriteOVF1(f, s, info, "text")
	case OVF1_BINARY:
		oommf.WriteOVF1(f, s, info, "binary 4")
	case OVF2_TEXT:
		oommf.WriteOVF2(f, s, info, "text")
	case OVF2_BINARY:
		oommf.WriteOVF2(f, s, info, "binary 4")
	default:
		panic("invalid output format")
	}

}

type outputFormatType int

const (
	OVF1_TEXT outputFormatType = iota + 1
	OVF1_BINARY
	OVF2_TEXT
	OVF2_BINARY
	DUMP
)

var (
	StringFromOutputFormat = map[outputFormatType]string{
		OVF1_TEXT:   "ovf",
		OVF1_BINARY: "ovf",
		OVF2_TEXT:   "ovf",
		OVF2_BINARY: "ovf",
		DUMP:        "dump"}
)
