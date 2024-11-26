package zarr

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/MathieuMoalic/amumax/src/data"
	"github.com/MathieuMoalic/amumax/src/fsutil"
	"github.com/MathieuMoalic/amumax/src/log"
)

type Metadata struct {
	Fields    map[string]interface{}
	Path      string
	startTime time.Time
	lastSave  time.Time
}

func (m *Metadata) Init(currentDir string, StartTime time.Time, GPUInfo string) {
	if m.Fields == nil {
		m.Fields = make(map[string]interface{})
	}
	m.Add("start_time", StartTime.Format(time.UnixDate))
	m.Add("gpu", GPUInfo)
	m.Path = currentDir + ".zattrs"
	m.startTime = StartTime
	m.Save()
	m.lastSave = time.Now()
}

func (m *Metadata) Add(key string, val interface{}) {
	if m.Fields == nil {
		m.Fields = make(map[string]interface{})
	}
	val_type := reflect.TypeOf(val).Kind()
	switch val_type {
	case reflect.Float64, reflect.Int, reflect.String, reflect.Bool:
		m.Fields[key] = val
	case reflect.Pointer:
		ptr_val := reflect.ValueOf(val).Elem()
		val_str := fmt.Sprintf("%v", ptr_val)
		val_str = val_str[1 : len(val_str)-1]
		m.Fields[key] = val_str
	case reflect.Array:
		m.Fields[key] = fmt.Sprintf("%v", val)
	case reflect.Func:
		// ignore functions
		return
	default:
		log.Log.Debug("Metadata key %s has invalid type %s: %v", key, val_type, val)
	}
}

func (m *Metadata) Get(key string) interface{} {
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
	} else {
		return false
	}
}
func (m *Metadata) Save() {
	if m.Path != "" {
		zattrs, err := fsutil.Create(m.Path)
		log.Log.PanicIfError(err)
		defer zattrs.Close()
		json_meta, err := json.MarshalIndent(m.Fields, "", "\t")
		log.Log.PanicIfError(err)
		_, err = zattrs.Write([]byte(json_meta))
		log.Log.PanicIfError(err)
	}
}

func (m *Metadata) AddMesh(mesh *data.MeshType) {
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
