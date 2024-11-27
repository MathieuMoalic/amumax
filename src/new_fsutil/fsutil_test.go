package new_fsutil_test

import (
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/MathieuMoalic/amumax/src/new_fsutil"
)

func TestFileSystem(t *testing.T) {
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
	if string(readData) != string(data) {
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
	if string(readData) != string(expectedData) {
		t.Errorf("Expected data %s, got %s", expectedData, readData)
	}

	// Test Create
	writer, err := fs.Create(filePath)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	newData := []byte("New content")
	_, err = writer.Write(newData)
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}
	err = writer.Close()
	if err != nil {
		t.Fatalf("Close failed: %v", err)
	}
	readData, err = fs.Read(filePath)
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}
	if string(readData) != string(newData) {
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
	if string(readData) != string(newData) {
		t.Errorf("Expected data %s, got %s", newData, readData)
	}

	// Test Remove
	err = fs.Remove("testdir")
	if err != nil {
		t.Fatalf("Remove failed: %v", err)
	}
	if fs.Exists("testdir") {
		t.Errorf("Directory testdir should not exist after removal")
	}
}
