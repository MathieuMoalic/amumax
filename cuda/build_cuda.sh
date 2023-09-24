#!/bin/bash
NVCCFLAGS="-std=c++03 -ccbin=/usr/bin/gcc --compiler-options -Werror --compiler-options -Wall -Xptxas -O3 -ptx"
CC=52
cd $(dirname -- "$0")
rm -f *_wrapper.go *.ptx cuda2go
go build -v cuda2go.go 
for file in *.cu; do
    bname=$(basename -s .cu $file)
    echo Compiling $file ...
    nvcc -I/opt/cuda/include/ $NVCCFLAGS -arch=compute_$CC -code=sm_$CC $file -o $bname\_$CC.ptx
    ./cuda2go $file
done
