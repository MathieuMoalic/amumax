package zarr_old

import (
	"encoding/binary"
	"math"
)

func Float64ToBytes(f float64) []byte {
	var buf [8]byte
	binary.LittleEndian.PutUint64(buf[:], math.Float64bits(f))
	return buf[:]
}

func BytesToFloat64(bytes []byte) float64 {
	bits := binary.LittleEndian.Uint64(bytes)
	float := math.Float64frombits(bits)
	return float
}

func BytesToFloat32(bytes []byte) float32 {
	bits := binary.LittleEndian.Uint32(bytes)
	float := math.Float32frombits(bits)
	return float
}

func Float32ToBytes(f float32) []byte {
	var buf [4]byte
	binary.LittleEndian.PutUint32(buf[:], math.Float32bits(f))
	return buf[:]
}
