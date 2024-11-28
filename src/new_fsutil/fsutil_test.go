package new_fsutil_test

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/MathieuMoalic/amumax/src/new_fsutil"
)

func TestFileSystemAsyncPutAndAppend(t *testing.T) {
	tempDir := t.TempDir()
	fs := new_fsutil.NewFileSystem(tempDir)

	filePath := "testfile_async.txt"
	data1 := []byte("Async first line\n")
	data2 := []byte("Async second line\n")

	// Test AsyncPut
	err := fs.AsyncPut(filePath, data1)
	if err != nil {
		t.Fatalf("AsyncPut failed: %v", err)
	}

	// Test AsyncAppend
	err = fs.AsyncAppend(filePath, data2)
	if err != nil {
		t.Fatalf("AsyncAppend failed: %v", err)
	}

	// Wait for operations to complete
	fs.Drain()

	// Verify file contents
	readData, err := fs.Read(filePath)
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}

	expectedData := append(data1, data2...)
	if !bytes.Equal(readData, expectedData) {
		t.Errorf("Expected data %q, got %q", expectedData, readData)
	}
}

func TestFileSystemMultipleAsyncAppends(t *testing.T) {
	tempDir := t.TempDir()
	fs := new_fsutil.NewFileSystem(tempDir)

	filePath := "testfile_multiple_async.txt"
	data := []byte("Line\n")

	// Queue multiple writes
	for i := 0; i < 50; i++ {
		err := fs.AsyncAppend(filePath, data)
		if err != nil {
			t.Fatalf("AsyncAppend failed: %v", err)
		}
	}

	// Drain the queue
	fs.Drain()

	// Verify file contents
	readData, err := fs.Read(filePath)
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}

	expectedData := bytes.Repeat(data, 50)
	if !bytes.Equal(readData, expectedData) {
		t.Errorf("Expected data of length %d, got %d", len(expectedData), len(readData))
	}
}

func TestFileSystemAsyncPutWithoutDrain(t *testing.T) {
	tempDir := t.TempDir()
	fs := new_fsutil.NewFileSystem(tempDir)

	filePath := "testfile_no_drain.txt"
	data := []byte("Data without drain\n")

	// Test AsyncPut without Drain
	err := fs.AsyncPut(filePath, data)
	if err != nil {
		t.Fatalf("AsyncPut failed: %v", err)
	}

	// Allow time for asynchronous operation
	time.Sleep(100 * time.Millisecond)

	// Attempt to read file
	readData, err := fs.Read(filePath)
	if err != nil {
		t.Logf("File may not exist yet as Drain was not called")
	} else if !bytes.Equal(readData, data) {
		t.Errorf("Expected data %q, got %q", data, readData)
	}

	// Call Drain and verify
	fs.Drain()
	readData, err = fs.Read(filePath)
	if err != nil {
		t.Fatalf("Read failed after Drain: %v", err)
	}
	if !bytes.Equal(readData, data) {
		t.Errorf("Expected data %q after Drain, got %q", data, readData)
	}
}

func TestFileSystemMkdir(t *testing.T) {
	tempDir := t.TempDir()
	fs := new_fsutil.NewFileSystem(tempDir)

	dirPath := "testdir/subdir"
	err := fs.Mkdir(dirPath)
	if err != nil {
		t.Fatalf("Mkdir failed: %v", err)
	}

	if !fs.Exists(dirPath) || !fs.IsDir(dirPath) {
		t.Errorf("Directory %s should exist and be a directory", dirPath)
	}
}

func TestFileSystemTouchAndUpdate(t *testing.T) {
	tempDir := t.TempDir()
	fs := new_fsutil.NewFileSystem(tempDir)

	filePath := "testfile_touch.txt"
	err := fs.Touch(filePath)
	if err != nil {
		t.Fatalf("Touch failed: %v", err)
	}
	if !fs.Exists(filePath) {
		t.Errorf("File %s should exist", filePath)
	}

	infoBefore, err := os.Stat(filepath.Join(tempDir, filePath))
	if err != nil {
		t.Fatalf("Stat failed: %v", err)
	}
	modTimeBefore := infoBefore.ModTime()

	// Wait and touch again
	time.Sleep(1 * time.Second)
	err = fs.Touch(filePath)
	if err != nil {
		t.Fatalf("Touch failed: %v", err)
	}

	infoAfter, err := os.Stat(filepath.Join(tempDir, filePath))
	if err != nil {
		t.Fatalf("Stat failed: %v", err)
	}
	modTimeAfter := infoAfter.ModTime()
	if !modTimeAfter.After(modTimeBefore) {
		t.Errorf("Modification time was not updated")
	}
}

func TestFileSystemPutAndRead(t *testing.T) {
	tempDir := t.TempDir()
	fs := new_fsutil.NewFileSystem(tempDir)

	filePath := "testfile_put_read.txt"
	data := []byte("Hello, World!")

	err := fs.Put(filePath, data)
	if err != nil {
		t.Fatalf("Put failed: %v", err)
	}

	readData, err := fs.Read(filePath)
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}

	if !bytes.Equal(readData, data) {
		t.Errorf("Expected data %s, got %s", data, readData)
	}
}

func TestFileSystemAppend(t *testing.T) {
	tempDir := t.TempDir()
	fs := new_fsutil.NewFileSystem(tempDir)

	filePath := "testfile_append.txt"
	initialData := []byte("Initial line")
	appendData := []byte("\nAppended line")

	err := fs.Put(filePath, initialData)
	if err != nil {
		t.Fatalf("Put failed: %v", err)
	}

	err = fs.Append(filePath, appendData)
	if err != nil {
		t.Fatalf("Append failed: %v", err)
	}

	expectedData := append(initialData, appendData...)
	readData, err := fs.Read(filePath)
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}

	if !bytes.Equal(readData, expectedData) {
		t.Errorf("Expected data %s, got %s", expectedData, readData)
	}
}

func TestFileSystemCreateAndWrite(t *testing.T) {
	tempDir := t.TempDir()
	fs := new_fsutil.NewFileSystem(tempDir)

	filePath := "testfile_create.txt"
	newData := []byte("New content")

	writer, file, err := fs.Create(filePath)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	_, err = writer.Write(newData)
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	err = writer.Flush()
	if err != nil {
		t.Fatalf("Flush failed: %v", err)
	}

	err = file.Close()
	if err != nil {
		t.Fatalf("Close failed: %v", err)
	}

	readData, err := fs.Read(filePath)
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}

	if !bytes.Equal(readData, newData) {
		t.Errorf("Expected data %s, got %s", newData, readData)
	}
}

func TestFileSystemOpen(t *testing.T) {
	tempDir := t.TempDir()
	fs := new_fsutil.NewFileSystem(tempDir)

	filePath := "testfile_open.txt"
	data := []byte("Open test content")

	err := fs.Put(filePath, data)
	if err != nil {
		t.Fatalf("Put failed: %v", err)
	}

	reader, err := fs.Open(filePath)
	if err != nil {
		t.Fatalf("Open failed: %v", err)
	}
	defer reader.Close()

	readData := make([]byte, len(data))
	_, err = reader.Read(readData)
	if err != nil && err != io.EOF {
		t.Fatalf("Read failed: %v", err)
	}

	if !bytes.Equal(readData, data) {
		t.Errorf("Expected data %s, got %s", data, readData)
	}
}
