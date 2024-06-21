#!/bin/sh
set -e

DEST=$1

if [ -z "$DEST" ]; then
  printf "Where to install amumax? [Default=$HOME/.local/bin]: "
  read DEST
  if [ -z "$DEST" ]; then
    DEST="$HOME/.local/bin"
  fi
fi

mkdir -p "$DEST"
DEST=$(realpath "$DEST")
case ":$PATH:" in
  *:"$DEST":*) ;;
  *) echo && echo " !!! WARNING !!! '$DEST' not in PATH !" && echo >&2 ;;
esac

cd $DEST
echo Downloading amumax from https://github.com/mathieumoalic ...
curl -Ls https://github.com/mathieumoalic/amumax/releases/latest/download/amumax > amumax

echo Downloading libcufft.so.10 from https://developer.download.nvidia.com/compute/cuda/redist/ ...
curl -Ls https://developer.download.nvidia.com/compute/cuda/redist/libcufft/linux-x86_64/libcufft-linux-x86_64-10.9.0.58-archive.tar.xz > tmp
tar -xvf tmp > /dev/null
cp -L libcufft-linux-x86_64-10.9.0.58-archive/lib/libcufft.so.10 .

echo Downloading libcurand.so.10 from https://developer.download.nvidia.com/compute/cuda/redist/ ...
curl -Ls https://developer.download.nvidia.com/compute/cuda/redist/libcurand/linux-x86_64/libcurand-linux-x86_64-10.3.6.39-archive.tar.xz > tmp
tar -xvf tmp > /dev/null
cp -L libcurand-linux-x86_64-10.3.6.39-archive/lib/libcurand.so.10 .

rm -r libcufft-linux-x86_64-10.9.0.58-archive
rm -r libcurand-linux-x86_64-10.3.6.39-archive
rm tmp

echo "Setting amumax as executable"
chmod +x amumax

echo "You can now use 'amumax'"
