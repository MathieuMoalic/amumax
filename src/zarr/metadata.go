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

func (m *Metadata) Init(currentDir string, startTime time.Time, gpuInfo string) {
	if m.Fields == nil {
		m.Fields = make(map[string]any)
	}
	m.Add("start_time", startTime.Format(time.UnixDate))
	m.Add("gpu", gpuInfo)
	m.Path = currentDir + ".zattrs"
	m.startTime = startTime
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
	log.Log.Info("Total simulation time: %s", time.Since(m.startTime))
	m.Save()
}

const metadataSaveInterval = 5 * time.Second

func (m *Metadata) NeedSave() bool {
	// save once every metadataSaveInterval
	if time.Since(m.lastSave) > metadataSaveInterval {
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

func (m *Metadata) AddMesh(meshData *mesh.Mesh) {
	m.Add("dx", meshData.Dx)
	m.Add("dy", meshData.Dy)
	m.Add("dz", meshData.Dz)
	m.Add("Nx", meshData.Nx)
	m.Add("Ny", meshData.Ny)
	m.Add("Nz", meshData.Nz)
	m.Add("Tx", meshData.Tx)
	m.Add("Ty", meshData.Ty)
	m.Add("Tz", meshData.Tz)
	m.Add("PBCx", meshData.PBCx)
	m.Add("PBCy", meshData.PBCy)
	m.Add("PBCz", meshData.PBCz)
}
