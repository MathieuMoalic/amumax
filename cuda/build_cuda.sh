#!/bin/bash
NVCCFLAGS="-std=c++03 -ccbin=/usr/bin/gcc --compiler-options -Werror --compiler-options -Wall -Xptxas -O3 -ptx"
rm -f *_wrapper.go *.ptx cuda2go
go build -v cuda2go.go 
for file in *.cu; do
    bname=$(basename -s .cu $file)
    echo Compiling $file ...
    for cc in 50 52 53 60 61 62 70 72 75 80; do
        nvcc -I/opt/cuda/include/ $NVCCFLAGS -arch=compute_$cc -code=sm_$cc $file -o $bname\_$cc.ptx
        ./cuda2go $file
    done
done