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
curl -Ls https://github.com/mathieumoalic/amumax/releases/latest/download/libcufft.so.11 > libcufft.so.11
echo Downloading CUDA rand
curl -Ls https://github.com/mathieumoalic/amumax/releases/latest/download/libcurand.so.10 > libcurand.so.10

echo "Setting amumax as executable"
chmod +x amumax

echo "You can now use 'amumax'"
