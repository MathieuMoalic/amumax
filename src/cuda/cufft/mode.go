package cufft

//#include <cufft.h>
import "C"

import (
	"fmt"
)

// CompatibilityMode CUFFT compatibility mode
type CompatibilityMode int

const (
	CompatibilityFFTWPadding CompatibilityMode = C.CUFFT_COMPATIBILITY_FFTW_PADDING
)

func (t CompatibilityMode) String() string {
	if str, ok := compatibilityModeString[t]; ok {
		return str
	}
	return fmt.Sprint("CUFFT Compatibility mode with unknown number:", int(t))
}

var compatibilityModeString map[CompatibilityMode]string = map[CompatibilityMode]string{
	CompatibilityFFTWPadding: "CUFFT_COMPATIBILITY_FFTW_PADDING",
}
