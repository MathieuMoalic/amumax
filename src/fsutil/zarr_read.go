package fsutil

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"io"
	"math"
	"path"
	"time"

	"github.com/DataDog/zstd"
	"github.com/MathieuMoalic/amumax/src/engine_old/data_old"
)

func (fs *FileSystem) ReadZarr(binaryPath string, od string, logInfo func(string, ...interface{})) (*data_old.Slice, error) {

	// Wait until all files are saved because we might be reading them now
	fs.waitForSave(logInfo)

	// Read and parse the .zarray file
	zarray, err := fs.readZarrayFile(binaryPath, logInfo)
	if err != nil {
		return nil, err
	}

	// Validate the compressor
	if zarray.Compressor.ID != "zstd" {
		return nil, errors.New("LoadFile: Only the Zstd compressor is supported")
	}

	// Create the data.Slice
	array := data_old.NewSlice(zarray.Chunks[4], [3]int{zarray.Chunks[3], zarray.Chunks[2], zarray.Chunks[1]})
	tensors := array.Tensors()

	// Read and decompress data with retries
	dataBytes, err := fs.readAndDecompressData(binaryPath, logInfo)
	if err != nil {
		return nil, err
	}

	// Reconstruct the tensors from the decompressed data
	err = fs.reconstructTensors(dataBytes, tensors)
	if err != nil {
		return nil, err
	}

	return array, nil
}

// waitForSave waits until IsSaving is false
func (fs *FileSystem) waitForSave(logInfo func(string, ...interface{})) {
	msg_sent := false
	for fs.queLen > 0 {
		if !msg_sent {
			logInfo("Waiting for all the files to be saved before reading...")
			msg_sent = true
		}
		time.Sleep(100 * time.Millisecond)
	}
}

// readZarrayFile reads and parses the .zarray file
func (fs *FileSystem) readZarrayFile(binaryPath string, logInfo func(string, ...interface{})) (*zarrayFile, error) {
	zarrayPath := path.Join(path.Dir(binaryPath), ".zarray")
	logInfo("Reading:  %v", binaryPath)

	content, err := fs.Read(zarrayPath)
	if err != nil {
		return nil, err
	}

	var zarray zarrayFile
	err = json.Unmarshal(content, &zarray)
	if err != nil {
		return nil, err
	}

	return &zarray, nil
}

// readAndDecompressData reads and decompresses data with retry logic
func (fs *FileSystem) readAndDecompressData(binaryPath string, logInfo func(string, ...interface{})) ([]byte, error) {
	const maxRetries = 5
	const retryDelay = 1 * time.Second

	var dataBytes []byte
	var lastErr error

	for retries := 0; retries < maxRetries; retries++ {
		// Open the file
		ioReader, err := fs.Open(binaryPath)
		if err != nil {
			lastErr = err
			logInfo("Error opening file: %v, retrying in %v...", err, retryDelay)
			time.Sleep(retryDelay)
			continue
		}

		// Read the compressed data
		compressedData, err := io.ReadAll(ioReader)
		ioReader.Close()
		if err != nil {
			lastErr = err
			logInfo("Error reading file: %v, retrying in %v...", err, retryDelay)
			time.Sleep(retryDelay)
			continue
		}

		// Check if compressedData is empty
		if len(compressedData) == 0 {
			lastErr = errors.New("compressed data is empty")
			logInfo("File is empty, retrying in %v...", retryDelay)
			time.Sleep(retryDelay)
			continue
		}

		// Decompress the data
		dataBytes, err = zstd.Decompress(nil, compressedData)
		if err != nil {
			lastErr = err
			logInfo("Decompression error: %v, retrying in %v...", err, retryDelay)
			time.Sleep(retryDelay)
			continue
		}

		// Successful decompression
		break
	}

	if dataBytes == nil {
		return nil, errors.New("Failed to read and decompress data after retries: " + lastErr.Error())
	}

	return dataBytes, nil
}

// reconstructTensors reconstructs the tensors from the decompressed data
func (fs *FileSystem) reconstructTensors(dataBytes []byte, tensors [][][][]float32) error {
	ncomp := len(tensors)
	sizez := len(tensors[0])
	sizey := len(tensors[0][0])
	sizex := len(tensors[0][0][0])

	expectedSize := sizex * sizey * sizez * ncomp * 4 // float32 is 4 bytes
	if len(dataBytes) != expectedSize {
		return errors.New("decompressed data size mismatch")
	}

	count := 0
	for iz := 0; iz < sizez; iz++ {
		for iy := 0; iy < sizey; iy++ {
			for ix := 0; ix < sizex; ix++ {
				for c := 0; c < ncomp; c++ {
					start := count * 4
					end := start + 4
					if end > len(dataBytes) {
						return errors.New("index out of range while reconstructing tensors")
					}
					tensors[c][iz][iy][ix] = bytesToFloat32(dataBytes[start:end])
					count++
				}
			}
		}
	}
	return nil
}

func bytesToFloat32(bytes []byte) float32 {
	bits := binary.LittleEndian.Uint32(bytes)
	float := math.Float32frombits(bits)
	return float
}
