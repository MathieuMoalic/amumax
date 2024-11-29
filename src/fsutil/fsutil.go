package fsutil

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"time"

	"github.com/fatih/color"
)

// FileSystem represents a file system with a working directory.
// It includes functionality for asynchronous file operations.
type FileSystem struct {
	Wd       string      // Working directory
	bufSize  int         // Buffer size for buffered writer
	filePerm os.FileMode // File permissions
	dirPerm  os.FileMode // Directory permissions

	// Fields for asynchronous operations
	saveQue chan func() // Channel for queued functions
	queLen  atom        // Number of tasks in the queue (atomic)
}

// atom is an atomic int32 used for counting queued tasks.
type atom int32

func (a *atom) Add(v int32) {
	atomic.AddInt32((*int32)(a), v)
}

func (a *atom) Load() int32 {
	return atomic.LoadInt32((*int32)(a))
}

// NewFileSystem creates a new FileSystem with the specified working directory.
// If scriptPath is "", a temporary directory is created.
// outputDir is "" by default, in which case the working directory is created folling the script path.
// If skipExists is true, the directory is skipped if it already exists. Default to false
// If forceClean is true, the directory is removed if it already exists. Default to false
func NewFileSystem(scriptPath string, outputDir string, skipExists, forceClean bool) (*FileSystem, string, error) {
	fs := &FileSystem{
		bufSize:  16 * 1024,              // Default buffer size for buffered writer (16 KB)
		filePerm: 0644,                   // Default file permissions
		dirPerm:  0755,                   // Default directory permissions
		saveQue:  make(chan func(), 100), // Initialize asynchronous operations
	}
	fs.Wd = fs.getWD(scriptPath, outputDir)

	if fs.IsDir("") {
		// if directory exists and --skip-exist flag is set, skip the directory
		if skipExists {
			warn := fmt.Sprintf("Directory `%s` exists, skipping `%s` because of --skip-exist flag.", fs.Wd, scriptPath)
			// os.Exit(0)
			return nil, warn, nil
			// if directory exists and --force-clean flag is set, remove the directory
		} else if forceClean {
			color.Yellow(fmt.Sprintf("Cleaning `%s`", fs.Wd))
			err := fs.Remove("")
			if err != nil {
				return nil, "", fmt.Errorf("error removing directory `%s`: %v", fs.Wd, err)
			}
			err = fs.Mkdir("")
			if err != nil {
				return nil, "", fmt.Errorf("error creating directory `%s`: %v", fs.Wd, err)
			}
		}
	} else {
		err := fs.Mkdir("")
		if err != nil {
			return nil, "", fmt.Errorf("error creating directory `%s`: %v", fs.Wd, err)
		}
	}
	err := fs.CreateZarrGroup("")
	if err != nil {
		return nil, "", fmt.Errorf("error creating zarr group `%s`: %v", fs.Wd, err)
	}
	return fs, "", nil
}

func (fs *FileSystem) getWD(scriptPath, outputDir string) string {
	zarrPath := ""
	if outputDir != "" {
		zarrPath = outputDir
	} else {
		if scriptPath == "" {
			now := time.Now()
			zarrPath = fmt.Sprintf("/tmp/amumax-%v-%02d-%02d_%02dh%02d.zarr", now.Year(), int(now.Month()), now.Day(), now.Hour(), now.Minute())
		} else {
			zarrPath = strings.TrimSuffix(scriptPath, ".mx3") + ".zarr"
		}
	}
	if !strings.HasSuffix(zarrPath, "/") {
		zarrPath += "/"
	}

	absZarrPath, err := filepath.Abs(zarrPath)
	if err != nil {
		absZarrPath = zarrPath // Use provided zarrPath if Abs fails
	}
	if !filepath.IsAbs(absZarrPath) {
		panic("working directory must be an absolute path")
	}
	return absZarrPath
}

// // run continuously executes tasks from the saveQue channel.
// func (fs *FileSystem) run() {
// 	for f := range fs.saveQue {
// 		f()
// 		fs.queLen.Add(-1)
// 	}
// }

// QueueOutput queues a function for asynchronous execution.
func (fs *FileSystem) QueueOutput(f func()) {
	fs.queLen.Add(1)
	fs.saveQue <- f
}

// Drain waits until all queued asynchronous operations are completed.
func (fs *FileSystem) Drain() {
	for fs.queLen.Load() > 0 {
		select {
		default:
			time.Sleep(1 * time.Millisecond)
		case f := <-fs.saveQue:
			f()
			fs.queLen.Add(-1)
		}
	}
}

// Mkdir creates a directory at the specified path, including any necessary parents.
// It does not return an error if the directory already exists.
func (fs *FileSystem) Mkdir(p string) error {
	p = fs.addWorkDir(p)
	return os.MkdirAll(p, fs.dirPerm)
}

// Touch creates an empty file at the specified path or updates its modification time if it exists.
func (fs *FileSystem) Touch(p string) error {
	p = fs.addWorkDir(p)
	currentTime := time.Now()

	// Check if file exists
	if _, err := os.Stat(p); os.IsNotExist(err) {
		// Create the file
		f, err := os.OpenFile(p, os.O_CREATE|os.O_WRONLY, fs.filePerm)
		if err != nil {
			return err
		}
		return f.Close()
	}
	// Update the modification and access times
	return os.Chtimes(p, currentTime, currentTime)
}

// Remove deletes the file or directory at the path, including any children.
func (fs *FileSystem) Remove(p string) error {
	p = fs.addWorkDir(p)
	return os.RemoveAll(p)
}

// Read reads the entire file and returns its contents.
func (fs *FileSystem) Read(p string) ([]byte, error) {
	p = fs.addWorkDir(p)
	return os.ReadFile(p)
}

// Exists checks if the file or directory at the path exists.
func (fs *FileSystem) Exists(p string) bool {
	p = fs.addWorkDir(p)
	_, err := os.Stat(p)
	return !os.IsNotExist(err)
}

// IsDir checks if the path is a directory.
func (fs *FileSystem) IsDir(p string) bool {
	p = fs.addWorkDir(p)
	fi, err := os.Stat(p)
	if err != nil {
		return false
	}
	return fi.IsDir()
}

// Append appends data to the file, creating it if it does not exist.
func (fs *FileSystem) Append(p string, data []byte) error {
	p = fs.addWorkDir(p)
	err := os.MkdirAll(filepath.Dir(p), fs.dirPerm)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(p, os.O_APPEND|os.O_WRONLY|os.O_CREATE, fs.filePerm)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(data)
	return err
}

// AsyncAppend queues a file append operation to append data to the specified path asynchronously.
func (fs *FileSystem) AsyncAppend(p string, data []byte) error {
	p = fs.addWorkDir(p)
	err := os.MkdirAll(filepath.Dir(p), fs.dirPerm)
	if err != nil {
		return err
	}
	fs.QueueOutput(func() {
		f, err := os.OpenFile(p, os.O_APPEND|os.O_WRONLY|os.O_CREATE, fs.filePerm)
		if err != nil {
			return
		}
		defer f.Close()
		_, _ = f.Write(data)
	})
	return nil
}

// Put creates a file at the path and writes data to it, overwriting if it exists.
func (fs *FileSystem) Put(p string, data []byte) error {
	p = fs.addWorkDir(p)
	err := os.MkdirAll(filepath.Dir(p), fs.dirPerm)
	if err != nil {
		return err
	}
	return os.WriteFile(p, data, fs.filePerm)
}

// AsyncPut queues a file write operation to write data to the specified path asynchronously.
func (fs *FileSystem) AsyncPut(path string, data []byte) error {
	path = fs.addWorkDir(path)
	err := os.MkdirAll(filepath.Dir(path), fs.dirPerm)
	if err != nil {
		return err
	}
	fs.QueueOutput(func() {
		if err := os.WriteFile(path, data, fs.filePerm); err != nil {
			color.Red(fmt.Sprintf("// Error writing file `%s`: %v", path, err))
		}
	})
	return nil
}

// Create opens a file for writing, truncating it if it exists.
func (fs *FileSystem) Create(p string) (*bufio.Writer, *os.File, error) {
	p = fs.addWorkDir(p)
	err := os.MkdirAll(filepath.Dir(p), fs.dirPerm)
	if err != nil {
		return nil, nil, err
	}
	f, err := os.OpenFile(p, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, fs.filePerm)
	if err != nil {
		return nil, nil, err
	}
	writer := bufio.NewWriterSize(f, fs.bufSize)
	return writer, f, nil
}

// Open opens a file for reading.
func (fs *FileSystem) Open(p string) (io.ReadCloser, error) {
	p = fs.addWorkDir(p)
	return os.Open(p)
}

// addWorkDir adds the working directory to the path if it's relative.
func (fs *FileSystem) addWorkDir(p string) string {
	if !filepath.IsAbs(p) {
		return filepath.Join(fs.Wd, p)
	}
	return p
}

// GetWD returns the working directory.
func (fs *FileSystem) GetWD() string {
	return fs.Wd
}
