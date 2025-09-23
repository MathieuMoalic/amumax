package zarr

import (
	"encoding/binary"
	"math"
)

// Float64ToBytes converts a float64 to a byte slice in little-endian order.
func Float64ToBytes(f float64) []byte {
	var buf [8]byte
	binary.LittleEndian.PutUint64(buf[:], math.Float64bits(f))
	return buf[:]
}

// BytesToFloat64 converts a byte slice in little-endian order to a float64.
func BytesToFloat64(bytes []byte) float64 {
	bits := binary.LittleEndian.Uint64(bytes)
	float := math.Float64frombits(bits)
	return float
}

// BytesToFloat32 converts a byte slice in little-endian order to a float32.
func BytesToFloat32(bytes []byte) float32 {
	bits := binary.LittleEndian.Uint32(bytes)
	float := math.Float32frombits(bits)
	return float
}

// Float32ToBytes converts a float32 to a byte slice in little-endian order.
func Float32ToBytes(f float32) []byte {
	var buf [4]byte
	binary.LittleEndian.PutUint32(buf[:], math.Float32bits(f))
	return buf[:]
}
