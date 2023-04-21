#!/bin/sh
set -e

printf "Where to install amumax? [Default=$HOME/.local/bin]: "
read DEST
if [ -z $DEST ];then
    DEST="$HOME/.local/bin"
fi
mkdir -p $DEST
DEST=$(realpath $DEST)
case :$PATH:
  in *:$DEST:*) ;;
     *) echo && echo " !!! WARNING !!! '$DEST' not in PATH !" && echo >&2;;
esac
cd $DEST
echo Downloading amumax
curl -Ls https://github.com/mathieumoalic/amumax/releases/latest/download/amumax > amumax
echo Downloading CUDA fft
curl -Ls https://developer.download.nvidia.com/compute/cuda/redist/libcufft/linux-x86_64/libcufft-linux-x86_64-11.0.2.54-archive.tar.xz > libcufft.tar.xz
echo Downloading CUDA rand
curl -Ls https://developer.download.nvidia.com/compute/cuda/redist/libcurand/linux-x86_64/libcurand-linux-x86_64-10.3.2.106-archive.tar.xz > libcurand.tar.xz

echo Extracting CUDA fft
tar xf libcufft.tar.xz  --wildcards "*/libcufft.so*"

echo Extracting CUDA rand
tar xf libcurand.tar.xz  --wildcards "*/libcurand.so*"
mv libcu*-archive/lib/lib* .

echo "Removing artifacts"
rm -rdf libcufft-linux-x86_64-11.0.2.54-archive libcurand-linux-x86_64-10.3.2.106-archive libcurand.tar.xz libcufft.tar.xz

echo "Setting amumax as executable"
chmod +x amumax

echo "You can now use 'amumax'"