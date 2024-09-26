package zarr

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/MathieuMoalic/amumax/httpfs"
	"github.com/MathieuMoalic/amumax/util"
)

type Metadata struct {
	Fields    map[string]interface{}
	Path      string
	startTime time.Time
}

func (m *Metadata) Init(currentDir string, StartTime time.Time, dx, dy, dz float64, Nx, Ny, Nz int, Tx, Ty, Tz float64, PBCx, PBCy, PBCz int, GPUInfo string) {
	if m.Fields == nil {
		m.Fields = make(map[string]interface{})
	}
	m.Fields["start_time"] = StartTime
	m.Fields["dx"] = dx
	m.Fields["dy"] = dy
	m.Fields["dz"] = dz
	m.Fields["Nx"] = Nx
	m.Fields["Ny"] = Ny
	m.Fields["Nz"] = Nz
	m.Fields["Tx"] = Tx
	m.Fields["Ty"] = Ty
	m.Fields["Tz"] = Tz
	m.Fields["PBCx"] = PBCx
	m.Fields["PBCy"] = PBCy
	m.Fields["PBCz"] = PBCz
	m.Fields["gpu"] = GPUInfo
	m.Path = currentDir + ".zattrs"
	m.startTime = StartTime
	m.Save()
}

func (m *Metadata) Add(key string, val interface{}) {
	if m.Fields == nil {
		m.Fields = make(map[string]interface{})
	}
	val_type := reflect.TypeOf(val).Kind()
	// if val is a float, int or string add it to the metadata
	switch val_type {
	case reflect.Float64, reflect.Int, reflect.String, reflect.Bool:
		m.Fields[key] = val
	case reflect.Pointer:
		ptr_val := reflect.ValueOf(val).Elem()
		val_str := fmt.Sprintf("%v", ptr_val)
		val_str = val_str[1 : len(val_str)-1]
		m.Fields[key] = val_str
		util.Log.Debug("Pointer key %s = %s: %v", key, val_type, val_str)
	case reflect.Array:
		util.Log.Debug("Array key %s = %s: %v", key, val_type, val)
		m.Fields[key] = fmt.Sprintf("%v", val)
	case reflect.Func:
		// ignore functions
		return
	default:
		util.Log.Debug("Metadata key %s has invalid type %s: %v", key, val_type, val)
	}
}

func (m *Metadata) Get(key string) interface{} {
	return m.Fields[key]
}

func (m *Metadata) End() {
	m.Fields["end_time"] = time.Now().Format(time.UnixDate)
	m.Fields["total_time"] = fmt.Sprint(time.Since(m.startTime))
	m.Save()
}

func (m *Metadata) Save() {
	if m.Path != "" {
		zattrs, err := httpfs.Create(m.Path)
		util.Log.PanicIfError(err)
		defer zattrs.Close()
		json_meta, err := json.MarshalIndent(m.Fields, "", "\t")
		util.Log.PanicIfError(err)
		_, err = zattrs.Write([]byte(json_meta))
		util.Log.PanicIfError(err)
	}
}
