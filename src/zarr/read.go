package zarr

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"io"
	"math"
	"os"
	"path"
	"time"

	"github.com/DataDog/zstd"
	"github.com/MathieuMoalic/amumax/src/data"
	"github.com/MathieuMoalic/amumax/src/httpfs"
	"github.com/MathieuMoalic/amumax/src/log"
)

func Read(binaryPath string, od string) (*data.Slice, error) {
	// Resolve the binary path to an absolute path
	binaryPath = resolvePath(binaryPath, od)

	// Wait until all files are saved because we might be reading them now
	waitForSave()

	// Read and parse the .zarray file
	zarray, err := readZarrayFile(binaryPath)
	if err != nil {
		return nil, err
	}

	// Validate the compressor
	if zarray.Compressor.ID != "zstd" {
		return nil, errors.New("LoadFile: Only the Zstd compressor is supported")
	}

	// Create the data.Slice
	array := data.NewSlice(zarray.Chunks[4], [3]int{zarray.Chunks[3], zarray.Chunks[2], zarray.Chunks[1]})
	tensors := array.Tensors()

	// Read and decompress data with retries
	dataBytes, err := readAndDecompressData(binaryPath)
	if err != nil {
		return nil, err
	}

	// Reconstruct the tensors from the decompressed data
	err = reconstructTensors(dataBytes, tensors)
	if err != nil {
		return nil, err
	}

	return array, nil
}

// resolvePath resolves the binary path to an absolute path
func resolvePath(binaryPath string, od string) string {
	if !path.IsAbs(binaryPath) {
		binaryPath = path.Join(path.Dir(od), binaryPath)
	}
	return path.Clean(binaryPath)
}

// waitForSave waits until IsSaving is false
func waitForSave() {
	msg_sent := false
	for IsSaving {
		if !msg_sent {
			log.Log.Info("Waiting for all the files to be saved before reading...")
			msg_sent = true
		}
		time.Sleep(100 * time.Millisecond)
	}
}

// readZarrayFile reads and parses the .zarray file
func readZarrayFile(binaryPath string) (*zarrayFile, error) {
	zarrayPath := path.Join(path.Dir(binaryPath), ".zarray")
	log.Log.Info("Reading:  %v", binaryPath)

	content, err := os.ReadFile(zarrayPath)
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
func readAndDecompressData(binaryPath string) ([]byte, error) {
	const maxRetries = 5
	const retryDelay = 1 * time.Second

	var dataBytes []byte
	var lastErr error

	for retries := 0; retries < maxRetries; retries++ {
		// Open the file
		ioReader, err := httpfs.Open(binaryPath)
		if err != nil {
			lastErr = err
			log.Log.Info("Error opening file: %v, retrying in %v...", err, retryDelay)
			time.Sleep(retryDelay)
			continue
		}

		// Read the compressed data
		compressedData, err := io.ReadAll(ioReader)
		ioReader.Close()
		if err != nil {
			lastErr = err
			log.Log.Info("Error reading file: %v, retrying in %v...", err, retryDelay)
			time.Sleep(retryDelay)
			continue
		}

		// Check if compressedData is empty
		if len(compressedData) == 0 {
			lastErr = errors.New("compressed data is empty")
			log.Log.Info("File is empty, retrying in %v...", retryDelay)
			time.Sleep(retryDelay)
			continue
		}

		// Decompress the data
		dataBytes, err = zstd.Decompress(nil, compressedData)
		if err != nil {
			lastErr = err
			log.Log.Info("Decompression error: %v, retrying in %v...", err, retryDelay)
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
func reconstructTensors(dataBytes []byte, tensors [][][][]float32) error {
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

// bytesToFloat32 converts a 4-byte slice to a float32
func bytesToFloat32(b []byte) float32 {
	return math.Float32frombits(binary.LittleEndian.Uint32(b))
}
