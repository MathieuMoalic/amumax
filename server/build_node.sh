#!/bin/bash
cd cuda
rm -f *_wrapper.go *.ptx cuda2go
go build cuda2go.go
for file in *.cu; do
    bname=$(basename -s .cu $file)
    for cc in 50 52 53 60 61 62 70 72 75 80; do
        nvcc $NVCCFLAGS -arch=compute_$cc -code=sm_$cc $file -o $bname\_$cc.ptx
        ./cuda2go $file
    done
done
cd .. 
go build -v