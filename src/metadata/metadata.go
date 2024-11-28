package metadata

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/MathieuMoalic/amumax/src/cuda"
	"github.com/MathieuMoalic/amumax/src/mesh"
	"github.com/MathieuMoalic/amumax/src/new_fsutil"
	"github.com/MathieuMoalic/amumax/src/new_log"
)

type Metadata struct {
	Fields    map[string]interface{}
	startTime time.Time
	lastSave  time.Time
	fs        *new_fsutil.FileSystem
	log       *new_log.Logs
}

func NewMetadata(fs *new_fsutil.FileSystem, log *new_log.Logs) *Metadata {
	m := &Metadata{}
	m.Fields = make(map[string]interface{})
	m.startTime = time.Now()
	m.lastSave = time.Now()
	m.fs = fs
	m.log = log
	m.Add("start_time", m.startTime.Format(time.UnixDate))
	m.Add("gpu", cuda.GPUInfo)
	m.Save()
	return m
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
		m.log.Debug("Metadata key %s has invalid type %s: %v", key, val_type, val)
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
	writer, file, err := m.fs.Create(".zattrs")
	m.log.PanicIfError(err)
	defer file.Close()
	json_meta, err := json.MarshalIndent(m.Fields, "", "\t")
	m.log.Debug("Saving metadata: %s", json_meta)
	m.log.Debug("File: %s", file.Name())
	m.log.PanicIfError(err)
	_, err = writer.Write(json_meta)
	writer.Flush()
	m.log.PanicIfError(err)
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
