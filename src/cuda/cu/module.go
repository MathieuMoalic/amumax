package cu

// This file implements loading of CUDA ptx modules

//#include <cuda.h>
import "C"

import (
	"unsafe"
)

// Module Represents a CUDA CUmodule, a reference to executable device code.
type Module uintptr

// ModuleLoad Loads a compute module from file
func ModuleLoad(fname string) Module {
	// fmt.Fprintln(os.Stderr, "driver.ModuleLoad", fname)
	var mod C.CUmodule
	err := Result(C.cuModuleLoad(&mod, C.CString(fname)))
	if err != SUCCESS {
		panic(err)
	}
	return Module(uintptr(unsafe.Pointer(mod)))
}

// ModuleLoadData Loads a compute module from string
func ModuleLoadData(image string) Module {
	var mod C.CUmodule
	err := Result(C.cuModuleLoadData(&mod, unsafe.Pointer(C.CString(image))))
	if err != SUCCESS {
		panic(err)
	}
	return Module(uintptr(unsafe.Pointer(mod)))
}

// ModuleGetFunction Returns a Function handle.
func ModuleGetFunction(module Module, name string) Function {
	var function C.CUfunction
	err := Result(C.cuModuleGetFunction(
		&function,
		C.CUmodule(unsafe.Pointer(uintptr(module))),
		C.CString(name)))
	if err != SUCCESS {
		panic(err)
	}
	return Function(uintptr(unsafe.Pointer(function)))
}

// GetFunction Returns a Function handle.
func (m Module) GetFunction(name string) Function {
	return ModuleGetFunction(m, name)
}
