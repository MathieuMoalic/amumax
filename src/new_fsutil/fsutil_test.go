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

func TestFileSystemWithAsyncOperations(t *testing.T) {
	tempDir := t.TempDir()
	fs := new_fsutil.NewFileSystem(tempDir)

	// Test AsyncPut and AsyncAppend
	filePath := "testfile_async.txt"
	data1 := []byte("Async first line\n")
	data2 := []byte("Async second line\n")

	err := fs.AsyncPut(filePath, data1)
	if err != nil {
		t.Fatalf("AsyncPut failed: %v", err)
	}

	err = fs.AsyncAppend(filePath, data2)
	if err != nil {
		t.Fatalf("AsyncAppend failed: %v", err)
	}

	// Wait for asynchronous operations to complete
	fs.Drain()

	// Read the file to check the contents
	readData, err := fs.Read(filePath)
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}

	expectedData := append(data1, data2...)
	if !bytes.Equal(readData, expectedData) {
		t.Errorf("Expected data %q, got %q", expectedData, readData)
	}

	// Test multiple asynchronous operations
	filePath = "testfile_multiple_async.txt"
	data := []byte("Line\n")

	// Queue multiple writes
	for i := 0; i < 50; i++ {
		err1 := fs.AsyncAppend(filePath, data)
		if err1 != nil {
			t.Fatalf("AsyncAppend failed: %v", err1)
		}
	}

	// Drain the queue
	fs.Drain()

	// Read the file
	readData, err = fs.Read(filePath)
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}

	expectedData = bytes.Repeat(data, 50)
	if !bytes.Equal(readData, expectedData) {
		t.Errorf("Expected data of length %d, got %d", len(expectedData), len(readData))
	}

	// Test AsyncPut without calling Drain (to demonstrate behavior)
	filePath = "testfile_no_drain.txt"
	data = []byte("Data without drain\n")

	err = fs.AsyncPut(filePath, data)
	if err != nil {
		t.Fatalf("AsyncPut failed: %v", err)
	}

	// Do not call fs.Drain()

	// Allow some time for the asynchronous operation to potentially complete
	time.Sleep(100 * time.Millisecond)

	// Attempt to read the file
	readData, err = fs.Read(filePath)
	if err != nil {
		// The file may not exist yet
		t.Logf("File may not exist yet as Drain was not called")
	} else if !bytes.Equal(readData, data) {
		t.Errorf("Expected data %q, got %q", data, readData)
	}

	// Now call Drain and check again
	fs.Drain()
	readData, err = fs.Read(filePath)
	if err != nil {
		t.Fatalf("Read failed after Drain: %v", err)
	}
	if !bytes.Equal(readData, data) {
		t.Errorf("Expected data %q after Drain, got %q", data, readData)
	}

	// Clean up
	err = fs.Remove("testdir")
	if err != nil {
		t.Fatalf("Remove failed: %v", err)
	}
}

func TestFileSystemBasicOperations(t *testing.T) {
	tempDir := t.TempDir()
	fs := new_fsutil.NewFileSystem(tempDir)

	// Test Mkdir
	dirPath := "testdir/subdir"
	err := fs.Mkdir(dirPath)
	if err != nil {
		t.Fatalf("Mkdir failed: %v", err)
	}
	fullDirPath := filepath.Join(tempDir, dirPath)
	if !fs.Exists(dirPath) || !fs.IsDir(dirPath) {
		t.Errorf("Directory %s should exist and be a directory", fullDirPath)
	}

	// Test Touch
	filePath := "testdir/subdir/testfile.txt"
	err = fs.Touch(filePath)
	if err != nil {
		t.Fatalf("Touch failed: %v", err)
	}
	if !fs.Exists(filePath) {
		t.Errorf("File %s should exist", filePath)
	}

	// Capture modification time
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

	// Check if modification time has been updated
	infoAfter, err := os.Stat(filepath.Join(tempDir, filePath))
	if err != nil {
		t.Fatalf("Stat failed: %v", err)
	}
	modTimeAfter := infoAfter.ModTime()
	if !modTimeAfter.After(modTimeBefore) {
		t.Errorf("Modification time was not updated")
	}

	// Test Put and Read
	data := []byte("Hello, World!")
	err = fs.Put(filePath, data)
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

	// Test Append
	appendData := []byte("\nAppended line")
	err = fs.Append(filePath, appendData)
	if err != nil {
		t.Fatalf("Append failed: %v", err)
	}
	expectedData := append(data, appendData...)
	readData, err = fs.Read(filePath)
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}
	if !bytes.Equal(readData, expectedData) {
		t.Errorf("Expected data %s, got %s", expectedData, readData)
	}

	// Test Create
	writer, file, err := fs.Create(filePath)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	newData := []byte("New content")
	_, err = writer.Write(newData)
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}
	err = file.Close()
	if err != nil {
		t.Fatalf("Close failed: %v", err)
	}
	readData, err = fs.Read(filePath)
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}
	if !bytes.Equal(readData, newData) {
		t.Errorf("Expected data %s, got %s", newData, readData)
	}

	// Test Open
	reader, err := fs.Open(filePath)
	if err != nil {
		t.Fatalf("Open failed: %v", err)
	}
	defer reader.Close()
	readData = make([]byte, len(newData))
	_, err = reader.Read(readData)
	if err != nil && err != io.EOF {
		t.Fatalf("Read failed: %v", err)
	}
	if !bytes.Equal(readData, newData) {
		t.Errorf("Expected data %s, got %s", newData, readData)
	}

	// Clean up
	err = fs.Remove("testdir")
	if err != nil {
		t.Fatalf("Remove failed: %v", err)
	}
}
