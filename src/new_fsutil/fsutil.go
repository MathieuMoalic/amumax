package new_fsutil

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync/atomic"
	"time"

	"github.com/fatih/color"
)

// FileSystem represents a file system with a working directory.
// It includes functionality for asynchronous file operations.
type FileSystem struct {
	wd       string      // Working directory
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
// The working directory should be an absolute path.
func NewFileSystem(wd string) *FileSystem {
	absWd, err := filepath.Abs(wd)
	if err != nil {
		absWd = wd // Use provided wd if Abs fails
	}
	if !filepath.IsAbs(absWd) {
		panic("working directory must be an absolute path")
	}
	fs := &FileSystem{
		wd:       absWd,
		bufSize:  16 * 1024, // Default buffer size for buffered writer (16 KB)
		filePerm: 0644,      // Default file permissions
		dirPerm:  0755,      // Default directory permissions

		// Initialize fields for asynchronous operations
		saveQue: make(chan func(), 100), // Queue capacity of 100
	}
	go fs.run()
	return fs
}

// run continuously executes tasks from the saveQue channel.
func (fs *FileSystem) run() {
	for f := range fs.saveQue {
		f()
		fs.queLen.Add(-1)
	}
}

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

// SetBufferSize sets the buffer size for the buffered writer.
func (fs *FileSystem) SetBufferSize(size int) {
	if size > 0 {
		fs.bufSize = size
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
func (fs *FileSystem) Create(p string) (WriteCloseFlusher, error) {
	p = fs.addWorkDir(p)
	err := os.MkdirAll(filepath.Dir(p), fs.dirPerm)
	if err != nil {
		return nil, err
	}
	f, err := os.OpenFile(p, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, fs.filePerm)
	if err != nil {
		return nil, err
	}
	writer := &bufWriter{
		buf:  bufio.NewWriterSize(f, fs.bufSize),
		file: f,
	}
	return writer, nil
}

// Open opens a file for reading.
func (fs *FileSystem) Open(p string) (io.ReadCloser, error) {
	p = fs.addWorkDir(p)
	return os.Open(p)
}

// WriteCloseFlusher represents a writer that can be flushed and closed.
// It is not safe for concurrent use by multiple goroutines.
type WriteCloseFlusher interface {
	io.WriteCloser
	Flush() error
}

type bufWriter struct {
	buf  *bufio.Writer
	file *os.File
}

func (w *bufWriter) Write(p []byte) (int, error) {
	return w.buf.Write(p)
}

func (w *bufWriter) Flush() error {
	return w.buf.Flush()
}

func (w *bufWriter) Close() error {
	err := w.Flush()
	if err != nil {
		w.file.Close()
		return err
	}
	return w.file.Close()
}

// addWorkDir adds the working directory to the path if it's relative.
func (fs *FileSystem) addWorkDir(p string) string {
	if !filepath.IsAbs(p) {
		return filepath.Join(fs.wd, p)
	}
	return p
}

// GetWD returns the working directory.
func (fs *FileSystem) GetWD() string {
	return fs.wd
}
