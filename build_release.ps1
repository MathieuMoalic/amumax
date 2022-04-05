$CUDAVERSION="11.6"
$OS="windows"
$BUILDDIR="./build/amumax_${OS}_${CUDAVERSION}"
$RPATH="lib" 
$CUDA_HOME="C:\Program Files\NVIDIA GPU Computing Toolkit\CUDA\v11.6"

# $NVCC_CCBIN="/usr/bin/gcc-8"
# $NVCCFLAGS="-std=c++03 -ccbin=$NVCC_CCBIN --compiler-options -Werror --compiler-options -Wall -Xptxas -O3 -ptx"
# cd cuda
# rm -f *_wrapper.go *.ptx cuda2go
# go build cuda2go.go
# for file in *.cu; do
#     echo $file
#     bname=$(basename -s .cu $file)
#     for cc in 50 52 53 60 61 62 70 72 75 80; do
#         nvcc $NVCCFLAGS -arch=compute_$cc -code=sm_$cc $file -o $bname\_$cc.ptx
#         ./cuda2go $file
#     done
# done
# cd ..

$CGO_CFLAGS="-I${LD_LIBRARY_PATH}"
$CGO_LDFLAGS="-lcufft -lcurand -lcuda -L${LD_LIBRARY_PATH} -Wl,-rpath -Wl,\$ORIGIN/$RPATH"
$CGO_CFLAGS_ALLOW='(-fno-schedule-insns|-malign-double|-ffast-math)'
$CGO_LDFLAGS="-lcufft -lcurand -lcuda -L${CUDA_HOME}/bin"

Remove-Item -ErrorAction Ignore -Recurse ${BUILDDIR}
Remove-Item -ErrorAction Ignore "${BUILDDIR}.zip"
mkdir ${BUILDDIR} 
go build -v -o ${BUILDDIR}
Move-Item amumax.exe -Destination ${BUILDDIR}/amumax.exe
Copy-Item ${CUDA_HOME}/bin/cufft64*.dll -Destination ${BUILDDIR}
Copy-Item ${CUDA_HOME}/bin/curand64*.dll -Destination ${BUILDDIR}
Compress-Archive -Path ${BUILDDIR}/* -DestinationPath "${BUILDDIR}.zip"