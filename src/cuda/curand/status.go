package curand

//#include <curand.h>
import "C"

import (
	"fmt"
)

type Status int

const (
	SUCCESS              Status = C.CURAND_STATUS_SUCCESS               // No errors
	VersionMismatch      Status = C.CURAND_STATUS_VERSION_MISMATCH      // Header file and linked library version do not match
	NotInitialized       Status = C.CURAND_STATUS_NOT_INITIALIZED       // Generator not initialized
	AllocationFailed     Status = C.CURAND_STATUS_ALLOCATION_FAILED     // Memory allocation failed
	TypeError            Status = C.CURAND_STATUS_TYPE_ERROR            // Generator is wrong type
	OutOfRange           Status = C.CURAND_STATUS_OUT_OF_RANGE          // Argument out of range
	LengthNotMultiple    Status = C.CURAND_STATUS_LENGTH_NOT_MULTIPLE   // Length requested is not a multple of dimension
	LaunchFailure        Status = C.CURAND_STATUS_LAUNCH_FAILURE        // Kernel launch failure
	PreExistingFailure   Status = C.CURAND_STATUS_PREEXISTING_FAILURE   // Preexisting failure on library entry
	InitializationFailed Status = C.CURAND_STATUS_INITIALIZATION_FAILED // Initialization of CUDA failed
	ArchMismatch         Status = C.CURAND_STATUS_ARCH_MISMATCH         // Architecture mismatch, GPU does not support requested feature
	InternalError        Status = C.CURAND_STATUS_INTERNAL_ERROR        // Internal library error
)

func (s Status) String() string {
	if str, ok := statusStr[s]; ok {
		return str
	}
	return fmt.Sprint("CURAND ERROR NUMBER ", int(s))
}

var statusStr = map[Status]string{
	SUCCESS:              "CURAND_STATUS_SUCCESS",
	VersionMismatch:      "CURAND_STATUS_VERSION_MISMATCH",
	NotInitialized:       "CURAND_STATUS_NOT_INITIALIZED",
	AllocationFailed:     "CURAND_STATUS_ALLOCATION_FAILED",
	TypeError:            "CURAND_STATUS_TYPE_ERROR",
	OutOfRange:           "CURAND_STATUS_OUT_OF_RANGE",
	LengthNotMultiple:    "CURAND_STATUS_LENGTH_NOT_MULTIPLE",
	LaunchFailure:        "CURAND_STATUS_LAUNCH_FAILURE",
	PreExistingFailure:   "CURAND_STATUS_PREEXISTING_FAILURE",
	InitializationFailed: "CURAND_STATUS_INITIALIZATION_FAILED",
	ArchMismatch:         "CURAND_STATUS_ARCH_MISMATCH",
	InternalError:        "CURAND_STATUS_INTERNAL_ERROR",
}

// Documentation was taken from the curand headers.
