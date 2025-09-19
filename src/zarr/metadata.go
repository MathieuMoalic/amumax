package zarr

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/MathieuMoalic/amumax/src/fsutil"
	"github.com/MathieuMoalic/amumax/src/log"
	"github.com/MathieuMoalic/amumax/src/mesh"
)

type Metadata struct {
	Fields    map[string]any
	Path      string
	startTime time.Time
	lastSave  time.Time
}

func (m *Metadata) Init(currentDir string, StartTime time.Time, GPUInfo string) {
	if m.Fields == nil {
		m.Fields = make(map[string]any)
	}
	m.Add("start_time", StartTime.Format(time.UnixDate))
	m.Add("gpu", GPUInfo)
	m.Path = currentDir + ".zattrs"
	m.startTime = StartTime
	m.Save()
	m.lastSave = time.Now()
}

func (m *Metadata) Add(key string, val any) {
	if m.Fields == nil {
		m.Fields = make(map[string]any)
	}
	valType := reflect.TypeOf(val).Kind()
	switch valType {
	case reflect.Float64, reflect.Int, reflect.String, reflect.Bool:
		m.Fields[key] = val
	case reflect.Pointer:
		ptrVal := reflect.ValueOf(val).Elem()
		valStr := fmt.Sprintf("%v", ptrVal)
		valStr = valStr[1 : len(valStr)-1]
		m.Fields[key] = valStr
	case reflect.Array:
		m.Fields[key] = fmt.Sprintf("%v", val)
	case reflect.Func:
		// ignore functions
		return
	default:
		log.Log.Debug("Metadata key %s has invalid type %s: %v", key, valType, val)
	}
}

func (m *Metadata) Get(key string) any {
	return m.Fields[key]
}

func (m *Metadata) End() {
	m.Add("end_time", time.Now().Format(time.UnixDate))
	m.Add("total_time", fmt.Sprint(time.Since(m.startTime)))
	m.Save()
}

func (m *Metadata) NeedSave() bool {
	// save once every 5 seconds
	if time.Since(m.lastSave) > 5*time.Second {
		m.lastSave = time.Now()
		return true
	}
	return false
}

func (m *Metadata) Save() {
	if m.Path != "" {
		zattrs, err := fsutil.Create(m.Path)
		log.Log.PanicIfError(err)
		defer func() {
			cerr := zattrs.Close()
			if cerr != nil {
				log.Log.Err("Error closing zattrs file: %v", cerr)
			}
		}()
		jsonMeta, err := json.MarshalIndent(m.Fields, "", "\t")
		log.Log.PanicIfError(err)
		_, err = zattrs.Write([]byte(jsonMeta))
		log.Log.PanicIfError(err)
	}
}

func (m *Metadata) AddMesh(mesh *mesh.Mesh) {
	m.Add("dx", mesh.Dx)
	m.Add("dy", mesh.Dy)
	m.Add("dz", mesh.Dz)
	m.Add("Nx", mesh.Nx)
	m.Add("Ny", mesh.Ny)
	m.Add("Nz", mesh.Nz)
	m.Add("Tx", mesh.Tx)
	m.Add("Ty", mesh.Ty)
	m.Add("Tz", mesh.Tz)
	m.Add("PBCx", mesh.PBCx)
	m.Add("PBCy", mesh.PBCy)
	m.Add("PBCz", mesh.PBCz)
}
