#!/bin/bash
# This file is used to build release versions in docker
rm -rf $BUILDDIR
mkdir -p $BUILDDIR/$RPATH

# cd cuda
# # rm -f *_wrapper.go *.ptx cuda2go
# go build cuda2go.go
# for file in *.cu; do
#     echo $file
#     bname=$(basename -s .cu $file)
#     for cc in 30 32 35 37 50 52 53 60 61 62 70 72 75 80; do
#         nvcc $NVCCFLAGS -arch=compute_$cc -code=sm_$cc $file -o $bname\_$cc.ptx
#         ./cuda2go $file
#     done
# done
# cd .. 

go build -v
mv ./amumax ./$BUILDDIR
cp $( ldd ${BUILDDIR}/${NAME} | grep libcufft | awk '{print $3}' ) ${BUILDDIR}/${RPATH}
cp $( ldd ${BUILDDIR}/${NAME} | grep libcurand | awk '{print $3}' ) ${BUILDDIR}/${RPATH}
# (cd build && tar -czf ${NAME}_linux.tar.gz ${NAME}_linux)