package metadata

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/MathieuMoalic/amumax/src/cuda"
	"github.com/MathieuMoalic/amumax/src/fsutil"
	"github.com/MathieuMoalic/amumax/src/log"
	"github.com/MathieuMoalic/amumax/src/mesh"
)

type Metadata struct {
	Fields        map[string]interface{}
	startTime     time.Time
	fs            *fsutil.FileSystem
	log           *log.Logs
	lastSavedHash [32]byte // Hash of the last saved Fields
}

func NewMetadata(fs *fsutil.FileSystem, log *log.Logs) *Metadata {
	m := &Metadata{}
	m.Fields = make(map[string]interface{})
	m.startTime = time.Now()
	m.fs = fs
	m.log = log
	m.Add("start_time", m.startTime.Format(time.UnixDate))
	m.Add("gpu", cuda.GPUInfo)
	m.FlushToFile()
	return m
}

func (m *Metadata) Add(key string, val interface{}) {
	if m.Fields == nil {
		m.Fields = make(map[string]interface{})
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
		m.log.Debug("Metadata key %s has invalid type %s: %v", key, valType, val)
	}
}

func (m *Metadata) Get(key string) interface{} {
	return m.Fields[key]
}

func (m *Metadata) Close() {
	m.Add("end_time", time.Now().Format(time.UnixDate))
	m.Add("total_time", fmt.Sprint(time.Since(m.startTime)))
	m.FlushToFile()
	// there is no need to close the file, as it is closed when the file is written
}

func (m *Metadata) FlushToFile() {
	// Compute the hash of the current Fields
	jsonMeta, err := json.Marshal(m.Fields)
	m.log.PanicIfError(err)

	currentHash := sha256.Sum256(jsonMeta)
	if currentHash == m.lastSavedHash {
		// No changes in Fields, skip saving
		return
	}

	// Save the metadata to the file
	writer, file, err := m.fs.Create(".zattrs")
	m.log.PanicIfError(err)
	defer file.Close()

	indentedJsonMeta, err := json.MarshalIndent(m.Fields, "", "\t")
	m.log.PanicIfError(err)

	_, err = writer.Write(indentedJsonMeta)
	writer.Flush()
	m.log.PanicIfError(err)

	// Update the hash after successful save
	m.lastSavedHash = currentHash
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
