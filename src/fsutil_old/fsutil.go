package fsutil_old

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var wd = "" // Working directory set by SetWD

const (
	DirPerm  = 0755             // Permissions for new directories
	FilePerm = 0644             // Permissions for new files
	BUFSIZE  = 16 * 1024 * 1024 // Buffer size for buffered writer
)

// SetWD sets the working directory, prefixed to all relative paths.
func SetWD(dir string) {
	if dir != "" && !strings.HasSuffix(dir, string(os.PathSeparator)) {
		dir += string(os.PathSeparator)
	}
	wd = dir
}

// Mkdir creates a directory at the specified path.
func Mkdir(p string) error {
	p = addWorkDir(p)
	return os.Mkdir(p, DirPerm)
}

// Touch creates an empty file at the specified path.
func Touch(p string) error {
	p = addWorkDir(p)
	f, err := os.OpenFile(p, os.O_CREATE|os.O_WRONLY, FilePerm)
	if err == nil {
		f.Close()
	}
	return err
}

// Remove deletes the file or directory at the path, including any children.
func Remove(p string) error {
	p = addWorkDir(p)
	return os.RemoveAll(p)
}

// Read reads the entire file and returns its contents.
func Read(p string) ([]byte, error) {
	p = addWorkDir(p)
	return os.ReadFile(p)
}

// Exists checks if the file or directory at the path exists.
func Exists(p string) bool {
	p = addWorkDir(p)
	_, err := os.Stat(p)
	return !os.IsNotExist(err)
}

// IsDir checks if the path is a directory.
func IsDir(p string) bool {
	p = addWorkDir(p)
	fi, err := os.Stat(p)
	if err != nil {
		return false
	}
	return fi.IsDir()
}

// Append appends data to the file.
func Append(p string, data []byte) error {
	p = addWorkDir(p)
	f, err := os.OpenFile(p, os.O_APPEND|os.O_WRONLY, FilePerm)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(data)
	return err
}

// Put creates a file at the path and writes data to it.
func Put(p string, data []byte) error {
	p = addWorkDir(p)
	err := os.MkdirAll(filepath.Dir(p), DirPerm)
	if err != nil {
		return err
	}
	return os.WriteFile(p, data, FilePerm)
}

// Create opens a file for writing, truncating it if it exists.
func Create(p string) (WriteCloseFlusher, error) {
	p = addWorkDir(p)
	err := os.Remove(p)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	f, err := os.OpenFile(p, os.O_CREATE|os.O_WRONLY, FilePerm)
	if err != nil {
		return nil, err
	}
	writer := &bufWriter{
		buf:  bufio.NewWriterSize(f, BUFSIZE),
		file: f,
	}
	return writer, nil
}

// Open opens a file for reading.
func Open(p string) (io.ReadCloser, error) {
	p = addWorkDir(p)
	return os.Open(p)
}

// WriteCloseFlusher represents a writer that can be flushed and closed.
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

// Helper function to add the working directory to the path if it's relative.
func addWorkDir(p string) string {
	if !filepath.IsAbs(p) {
		return filepath.Join(wd, p)
	}
	return p
}
