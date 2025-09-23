package engine

import (
	"strconv"
	"strings"

	"github.com/MathieuMoalic/amumax/src/fsutil"
	"github.com/MathieuMoalic/amumax/src/log"
)

func init() {
	DeclFunc("readArrayFromFile", readArrayFromFile, "")
	DeclFunc("readArrayFromString", readArrayFromString, "")
}

type InputArray struct {
	data []float64
}

func (a InputArray) Get(i int) float64 {
	if i < 0 || i >= len(a.data) {
		log.Log.ErrAndExit("InputArray.Get: index out of bounds")
		return 0
	}
	return a.data[i]
}

func (a InputArray) Len() int {
	return len(a.data)
}

func readArrayFromFile(filename string) InputArray {
	bytes, err := fsutil.Read(filename)
	if err != nil {
		log.Log.Err("readArrayFromFile: error reading file: %v", err)
		return InputArray{}
	}
	return readArrayFromString(string(bytes))
}

func readArrayFromString(s string) InputArray {
	s = strings.TrimSpace(s)

	// Handle empty input
	if strings.TrimSpace(s) == "" {
		log.Log.Warn("parseFloatArray: empty input")
		return InputArray{}
	}

	// Split by comma
	parts := strings.Split(s, ",")
	result := make([]float64, 0, len(parts))

	for _, part := range parts {
		numStr := strings.TrimSpace(part)
		num, err := strconv.ParseFloat(numStr, 64)
		if err != nil {
			log.Log.Err("parseFloatArray: error parsing '%s': %v", numStr, err)
			return InputArray{}
		}
		result = append(result, num)
	}

	return InputArray{data: result}
}
