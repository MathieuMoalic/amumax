SHELL = /bin/bash

# When CUDA_HOME is not an environment variable and is not set on the command line, use the nvcc compiler
# from the PATH
ifeq ($(CUDA_HOME),)
	NVCC=nvcc
else 
	NVCC=${CUDA_HOME}/bin/nvcc 
endif

# When CUDA_CC is not an environment variable and is not set on the command line, use compute capability 3.0
ifeq ($(CUDA_CC),)
    CUDA_CC = 30
endif

# The gcc host compiler for nvcc
ifeq ($(NVCC_CCBIN),)
	NVCC_CCBIN=/usr/bin/gcc-8
endif

CUDA_VERSION := $(shell $(NVCC) --version | grep "Cuda compilation" | grep -Eo '[+-]?[0-9]+([.][0-9]+)?' | head -n 1)

NVCC_COMPATIBILITY_FLAGS := -std=c++03
ifneq (,$(filter 7.0 7.5 8.0,$(CUDA_VERSION)))
	NVCC_COMPATIBILITY_FLAGS :=
endif

NVCCFLAGS = $(NVCC_COMPATIBILITY_FLAGS) -ccbin=$(NVCC_CCBIN) --compiler-options -Werror --compiler-options -Wall -Xptxas -O3 -ptx

CUDAFILES := $(wildcard *.cu)
WRAPPERS := $(CUDAFILES:.cu=_wrapper.go)


.PHONY: all wrappers clean realclean


all: wrappers
	@echo "Built with CUDA version ${CUDA_VERSION}"
	go install -v


wrappers: $(WRAPPERS)


%_wrapper.go: %.cu cuda2go
	@ rm -f $(basename $<)*.ptx
	@ for cc in $(CUDA_CC); do \
		echo $(NVCC) $(NVCCFLAGS) -arch=compute_$$cc -code=sm_$$cc $< -o $(basename $<)_$$cc.ptx ;\
		     $(NVCC) $(NVCCFLAGS) -arch=compute_$$cc -code=sm_$$cc $< -o $(basename $<)_$$cc.ptx ;\
	done
	@ ./cuda2go $< > /dev/null
	@ gofmt -w -s -l $@ > /dev/null


cuda2go: cuda2go.go
	go build $<


clean:
	rm -vf *.ptx


realclean:
	rm -vf *_wrapper.go *.ptx cuda2go