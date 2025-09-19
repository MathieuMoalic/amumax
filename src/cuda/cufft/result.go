// Package cufft provides bindings for the cuFFT library.
package cufft

//#include <cufft.h>
import "C"

import (
	"fmt"
)

// Result FFT result
type Result int

// NOWORKSPACE PARSEERROR INVALIDDEVICE INCOMPLETEPARAMETERLIST UNALIGNEDDATA INVALIDSIZE SETUPFAILED EXECFAILED INTERNALERROR INVALIDVALUE INVALIDTYPE ALLOCFAILED INVALIDPLAN SUCCESS FFT result value
const (
	SUCCESS                 Result = C.CUFFT_SUCCESS
	INVALIDPLAN             Result = C.CUFFT_INVALID_PLAN
	ALLOCFAILED             Result = C.CUFFT_ALLOC_FAILED
	INVALIDTYPE             Result = C.CUFFT_INVALID_TYPE
	INVALIDVALUE            Result = C.CUFFT_INVALID_VALUE
	INTERNALERROR           Result = C.CUFFT_INTERNAL_ERROR
	EXECFAILED              Result = C.CUFFT_EXEC_FAILED
	SETUPFAILED             Result = C.CUFFT_SETUP_FAILED
	INVALIDSIZE             Result = C.CUFFT_INVALID_SIZE
	UNALIGNEDDATA           Result = C.CUFFT_UNALIGNED_DATA
	INCOMPLETEPARAMETERLIST Result = 0xA // cuda6 values copied to avoid dependency on cuda6/cufft.h
	INVALIDDEVICE           Result = 0xB
	PARSEERROR              Result = 0xC
	NOWORKSPACE             Result = 0xD
)

func (r Result) String() string {
	if str, ok := resultString[r]; ok {
		return str
	}
	return fmt.Sprint("CUFFT Result with unknown error number:", int(r))
}

var resultString map[Result]string = map[Result]string{
	SUCCESS:                 "CUFFT_SUCCESS",
	INVALIDPLAN:             "CUFFT_INVALID_PLAN",
	ALLOCFAILED:             "CUFFT_ALLOC_FAILED",
	INVALIDTYPE:             "CUFFT_INVALID_TYPE",
	INVALIDVALUE:            "CUFFT_INVALID_VALUE",
	INTERNALERROR:           "CUFFT_INTERNAL_ERROR",
	EXECFAILED:              "CUFFT_EXEC_FAILED",
	SETUPFAILED:             "CUFFT_SETUP_FAILED",
	INVALIDSIZE:             "CUFFT_INVALID_SIZE",
	UNALIGNEDDATA:           "CUFFT_UNALIGNED_DATA",
	INCOMPLETEPARAMETERLIST: "CUFFT_INCOMPLETE_PARAMETER_LIST",
	INVALIDDEVICE:           "CUFFT_INVALID_DEVICE",
	PARSEERROR:              "CUFFT_PARSE_ERROR",
	NOWORKSPACE:             "CUFFT_NO_WORKSPACE",
}
