package log

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/MathieuMoalic/amumax/src/fsutil"
)

func setupFileSystem(t *testing.T, tempDir string) *fsutil.FileSystem {
	fs := fsutil.NewFileSystem(tempDir)
	if !fs.IsDir(tempDir) {
		t.Fatalf("Failed to set up test file system")
	}
	return fs
}

func createTempDir(t *testing.T) string {
	tempDir, err := os.MkdirTemp("", "logtest")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	return tempDir
}

func cleanupTempDir(t *testing.T, tempDir string) {
	err := os.RemoveAll(tempDir)
	if err != nil {
		t.Fatalf("Failed to clean up temporary directory: %v", err)
	}
}

func TestNewLogs(t *testing.T) {
	tempDir := createTempDir(t)
	defer cleanupTempDir(t, tempDir)

	fs := setupFileSystem(t, tempDir)
	logPath := filepath.Join(tempDir, "logs")
	logs := NewLogs(logPath, fs, true)

	if logs == nil || logs.file == nil || logs.writer == nil {
		t.Fatal("Failed to initialize Logs")
	}
	logs.Close()
}

func TestLogWriting(t *testing.T) {
	tempDir := createTempDir(t)
	defer cleanupTempDir(t, tempDir)

	fs := setupFileSystem(t, tempDir)
	logPath := filepath.Join(tempDir, "logs")
	logs := NewLogs(logPath, fs, false)
	defer logs.Close()

	// Log a command
	logs.Command("Test command")
	logs.FlushToFile()

	// Verify the log file contains the expected content
	logFilePath := filepath.Join(tempDir, "logs/log.txt")
	content, err := os.ReadFile(logFilePath)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}
	if !strings.Contains(string(content), "Test command") {
		t.Errorf("Log file does not contain expected content: %s", content)
	}
}

func TestAutoFlush(t *testing.T) {
	tempDir := createTempDir(t)
	defer cleanupTempDir(t, tempDir)

	fs := setupFileSystem(t, tempDir)
	logPath := filepath.Join(tempDir, "logs")
	logs := NewLogs(logPath, fs, false)
	defer logs.Close()

	// Run AutoFlush in a separate goroutine
	go logs.AutoFlushToFile()

	// Log a message and wait briefly
	logs.Command("AutoFlush test")
	time.Sleep(6 * time.Second) // Ensure AutoFlush is triggered

	// Verify the log file content
	logFilePath := filepath.Join(tempDir, "logs/log.txt")
	content, err := os.ReadFile(logFilePath)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}
	if !strings.Contains(string(content), "AutoFlush test") {
		t.Errorf("Log file does not contain expected content: %s", content)
	}
}

func TestErrorLogging(t *testing.T) {
	tempDir := createTempDir(t)
	defer cleanupTempDir(t, tempDir)

	fs := setupFileSystem(t, tempDir)
	logPath := filepath.Join(tempDir, "logs")
	logs := NewLogs(logPath, fs, false)
	defer logs.Close()

	logs.Err("Test error: %d", 42)
	logs.FlushToFile()

	// Verify the log file contains the error message
	logFilePath := filepath.Join(tempDir, "logs/log.txt")
	content, err := os.ReadFile(logFilePath)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}
	if !strings.Contains(string(content), "Test error: 42") {
		t.Errorf("Log file does not contain expected error message: %s", content)
	}
}

func TestDebugLogging(t *testing.T) {
	tempDir := createTempDir(t)
	defer cleanupTempDir(t, tempDir)

	fs := setupFileSystem(t, tempDir)
	logPath := filepath.Join(tempDir, "logs")
	logs := NewLogs(logPath, fs, true) // Debug mode enabled
	defer logs.Close()

	logs.Debug("Debugging: %s", "test")
	logs.FlushToFile()

	// Verify debug message is logged
	logFilePath := filepath.Join(tempDir, "logs/log.txt")
	content, err := os.ReadFile(logFilePath)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}
	if !strings.Contains(string(content), "Debugging: test") {
		t.Errorf("Log file does not contain expected debug message: %s", content)
	}
}

func TestPanicIfError(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("PanicIfError did not panic as expected")
		}
	}()

	tempDir := createTempDir(t)
	defer cleanupTempDir(t, tempDir)

	fs := setupFileSystem(t, tempDir)
	logPath := filepath.Join(tempDir, "logs")
	logs := NewLogs(logPath, fs, false)
	defer logs.Close()

	logs.PanicIfError(fmt.Errorf("This is a test error"))
}
