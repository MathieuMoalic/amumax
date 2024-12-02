#!/bin/env bash
set -e
cd $(dirname -- "$0")
rm -f *_wrapper.go *.ptx cuda2go
go build -v cuda2go.go 
NVCCFLAGS="-std=c++03 -ccbin=gcc --compiler-options -Werror --compiler-options -Wall -Xptxas -O3 -ptx"
COMPCAP=52 # Oldest supported architecture for compatibility
for file in *.cu; do
    bname=$(basename -s .cu $file)
    echo Compiling $file ...
    nvcc -I/opt/cuda/include/ $NVCCFLAGS -arch=compute_$COMPCAP -code=sm_$COMPCAP $file -o $bname\_$COMPCAP.ptx
    ./cuda2go $file
done
