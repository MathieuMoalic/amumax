#!/bin/bash
NVCCFLAGS="-std=c++03 -ccbin=/usr/bin/gcc --compiler-options -Werror --compiler-options -Wall -Xptxas -O3 -ptx"
cd $(dirname -- "$0")
rm -f *_wrapper.go *.ptx cuda2go
go build -v cuda2go.go 
for file in *.cu; do
    bname=$(basename -s .cu $file)
    echo Compiling $file ...
    for cc in 52; do
        nvcc -I/opt/cuda/include/ $NVCCFLAGS -arch=compute_$cc -code=sm_$cc $file -o $bname\_$cc.ptx
        ./cuda2go $file
    done
done
