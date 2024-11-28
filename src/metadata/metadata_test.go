package metadata

// import (
// 	"encoding/json"
// 	"os"
// 	"path/filepath"
// 	"testing"
// 	"time"

// 	"github.com/MathieuMoalic/amumax/src/fsutil"
// 	"github.com/MathieuMoalic/amumax/src/log"
// 	"github.com/MathieuMoalic/amumax/src/mesh"
// )

// func createTestFileSystem(t *testing.T) *fsutil.FileSystem {
// 	tempDir, err := os.MkdirTemp("", "metadata_test")
// 	if err != nil {
// 		t.Fatalf("Failed to create temporary directory: %v", err)
// 	}
// 	return fsutil.NewFileSystem(tempDir)
// }

// func createTestLogger(t *testing.T) *log.Logs {
// 	tempDir, err := os.MkdirTemp("", "log_test")
// 	if err != nil {
// 		t.Fatalf("Failed to create temporary directory for logs: %v", err)
// 	}
// 	fs := fsutil.NewFileSystem(tempDir)
// 	return log.NewLogs(tempDir, fs, true)
// }

// func cleanupTestFileSystem(_ *testing.T, fs *fsutil.FileSystem) {
// 	os.RemoveAll(fs.GetWD())
// }

// func TestNewMetadataInitialization(t *testing.T) {
// 	fs := createTestFileSystem(t)
// 	defer cleanupTestFileSystem(t, fs)

// 	log := createTestLogger(t)
// 	defer log.Close()

// 	meta := NewMetadata(fs, log)

// 	if meta == nil {
// 		t.Fatal("NewMetadata returned nil")
// 	}
// 	if _, ok := meta.Fields["start_time"]; !ok {
// 		t.Errorf("Metadata does not have 'start_time'")
// 	}
// 	if _, ok := meta.Fields["gpu"]; !ok {
// 		t.Errorf("Metadata does not have 'gpu'")
// 	}
// }

// func TestMetadataAddAndGet(t *testing.T) {
// 	fs := createTestFileSystem(t)
// 	defer cleanupTestFileSystem(t, fs)

// 	log := createTestLogger(t)
// 	defer log.Close()

// 	meta := NewMetadata(fs, log)
// 	meta.Add("testKey", "testValue")

// 	val := meta.Get("testKey")
// 	if val != "testValue" {
// 		t.Errorf("Expected 'testValue', got %v", val)
// 	}
// }

// func TestMetadataEnd(t *testing.T) {
// 	fs := createTestFileSystem(t)
// 	defer cleanupTestFileSystem(t, fs)

// 	log := createTestLogger(t)
// 	defer log.Close()

// 	meta := NewMetadata(fs, log)
// 	time.Sleep(1 * time.Second) // Ensure a non-zero total time
// 	meta.End()

// 	if _, ok := meta.Fields["end_time"]; !ok {
// 		t.Errorf("Metadata does not have 'end_time'")
// 	}
// 	if _, ok := meta.Fields["total_time"]; !ok {
// 		t.Errorf("Metadata does not have 'total_time'")
// 	}
// }

// func TestMetadataNeedSave(t *testing.T) {
// 	fs := createTestFileSystem(t)
// 	defer cleanupTestFileSystem(t, fs)

// 	log := createTestLogger(t)
// 	defer log.Close()

// 	meta := NewMetadata(fs, log)

// 	if meta.NeedSave() {
// 		t.Errorf("Metadata incorrectly reports NeedSave as true immediately after creation")
// 	}

// 	time.Sleep(6 * time.Second)
// 	if !meta.NeedSave() {
// 		t.Errorf("Metadata should report NeedSave as true after 6 seconds")
// 	}
// }

// func TestMetadataSave(t *testing.T) {
// 	fs := createTestFileSystem(t)
// 	defer cleanupTestFileSystem(t, fs)

// 	log := createTestLogger(t)
// 	defer log.Close()

// 	meta := NewMetadata(fs, log)
// 	meta.Add("key", "value")
// 	meta.Save()

// 	metaFilePath := filepath.Join(fs.GetWD(), ".zattrs")
// 	if _, err := os.Stat(metaFilePath); os.IsNotExist(err) {
// 		t.Fatalf("Metadata file '.zattrs' not created")
// 	}

// 	data, err := os.ReadFile(metaFilePath)
// 	if err != nil {
// 		t.Fatalf("Failed to read '.zattrs': %v", err)
// 	}

// 	var savedFields map[string]interface{}
// 	if err := json.Unmarshal(data, &savedFields); err != nil {
// 		t.Fatalf("Failed to unmarshal metadata: %v", err)
// 	}

// 	if savedFields["key"] != "value" {
// 		t.Errorf("Expected saved key to have value 'value', got %v", savedFields["key"])
// 	}
// }

// func TestMetadataAddMesh(t *testing.T) {
// 	fs := createTestFileSystem(t)
// 	defer cleanupTestFileSystem(t, fs)

// 	log := createTestLogger(t)
// 	defer log.Close()

// 	meta := NewMetadata(fs, log)

// 	m := &mesh.Mesh{
// 		Nx: 64, Ny: 64, Nz: 1,
// 		Dx: 1e-9, Dy: 1e-9, Dz: 1e-9,
// 		Tx: 0, Ty: 0, Tz: 0,
// 		PBCx: 1, PBCy: 1, PBCz: 0,
// 	}

// 	meta.AddMesh(m)

// 	for _, key := range []string{"Nx", "Ny", "Nz", "dx", "dy", "dz", "Tx", "Ty", "Tz", "PBCx", "PBCy", "PBCz"} {
// 		if _, ok := meta.Fields[key]; !ok {
// 			t.Errorf("Expected key %s to be present in metadata", key)
// 		}
// 	}
// }

// func TestMetadataAddInvalidType(t *testing.T) {
// 	fs := createTestFileSystem(t)
// 	defer cleanupTestFileSystem(t, fs)

// 	log := createTestLogger(t)
// 	defer log.Close()

// 	meta := NewMetadata(fs, log)

// 	ch := make(chan int)
// 	meta.Add("invalid", ch)

// 	if _, ok := meta.Fields["invalid"]; ok {
// 		t.Errorf("Expected invalid key not to be added to metadata")
// 	}
// }
