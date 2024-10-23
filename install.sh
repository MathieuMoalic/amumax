#!/bin/sh
set -e

# Check for necessary commands
for cmd in curl tar xz; do
  if ! command -v $cmd > /dev/null 2>&1; then
    echo "Error: '$cmd' is not installed. Please install it and try again." >&2
    exit 1
  fi
done

DEST=$1

# Prompt for installation path if not provided
if [ -z "$DEST" ]; then
  printf "Where to install amumax? [Default=$HOME/.local/bin]: "
  read DEST
  if [ -z "$DEST" ]; then
    DEST="$HOME/.local/bin"
  fi
fi

mkdir -p "$DEST"
DEST=$(realpath "$DEST")

# Warn if DEST is not in PATH
case ":$PATH:" in
  *:"$DEST":*) ;;
  *) 
    echo && echo " !!! WARNING !!! '$DEST' not in PATH!"
    echo "Consider adding '$DEST' to your PATH." >&2
    ;;
esac

# Download and install amumax
cd $DEST
echo "Downloading amumax from GitHub..."
curl -Ls https://github.com/mathieumoalic/amumax/releases/latest/download/amumax -o amumax

# Download and extract necessary libraries
echo "Downloading and extracting libcufft.so.10..."
curl -Ls https://developer.download.nvidia.com/compute/cuda/redist/libcufft/linux-x86_64/libcufft-linux-x86_64-10.9.0.58-archive.tar.xz -o tmp
tar -xJf tmp > /dev/null
cp -L libcufft-linux-x86_64-10.9.0.58-archive/lib/libcufft.so.10 .

echo "Downloading and extracting libcurand.so.10..."
curl -Ls https://developer.download.nvidia.com/compute/cuda/redist/libcurand/linux-x86_64/libcurand-linux-x86_64-10.3.6.39-archive.tar.xz -o tmp
tar -xJf tmp > /dev/null
cp -L libcurand-linux-x86_64-10.3.6.39-archive/lib/libcurand.so.10 .

# Clean up
rm -r libcufft-linux-x86_64-10.9.0.58-archive
rm -r libcurand-linux-x86_64-10.3.6.39-archive
rm tmp

# Make amumax executable
echo "Setting amumax as executable"
chmod +x amumax

# Completion message
echo "Installation complete. You can now use 'amumax'."
