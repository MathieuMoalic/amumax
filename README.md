# Amumax

fork of [mumax3](https://github.com/mumax/3) meant to increase the integration with a python processing workflow. I made my own wrapper around zarr called [pyzfn](https://github.com/MathieuMoalic/pyzfn) which leverages the mumax data in the zarr format.

The solvers ( and results ) are unchanged, this is just list of massive quality of life improvements making working with the output data much more efficient and convenient.

It's not 100% compatible with the original `.mx3` files. See changes below.

## Installation

### Linux

#### Install script

Don't just run an script on the internet. [Read it](https://raw.githubusercontent.com/MathieuMoalic/amumax/main/install.sh), check what it does and then you can run this command to install amumax:

```bash
sh -c "$(curl -fsSL https://raw.githubusercontent.com/MathieuMoalic/amumax/main/install.sh)"
```

#### Manually

Download [cufft](https://developer.download.nvidia.com/compute/cuda/redist/libcufft/linux-x86_64/) and [curand](https://developer.download.nvidia.com/compute/cuda/redist/libcurand/linux-x86_64/), unpack and add the shared objects to $PATH, or just install the full CUDA suite with your package manager.
Download the latest [release](https://github.com/MathieuMoalic/amumax/releases/):

```bash
curl -L https://github.com/mathieumoalic/amumax/releases/latest/download/amumax > amumax
amumax -v
```

`libcurand.so` and `libcufft.so` must either be in the same folder as `amumax` or on $PATH.

### Windows

- Get linux: https://learn.microsoft.com/en-us/windows/wsl/install
- Follow the steps above.

## Differences from mumax3

### New way to define the mesh

`Tx`,`Ty`,`Tz`,`Nx`,`Ny`,`Nz`,`dx`,`dy`,`dz`,`PBCx`,`PBCy`,`PBCz` are now predefined variables. You define the Mesh through them. You don't need to call a function to initiate the Mesh, it is automatically done the first time you run a solver but you can't redefine the Mesh after that !
`Tx` is the total size along the x axis. `Nx` is the number of cells along the x axis. `dx` is the number of cells along the x axis. Keep in mind that variables in mx3 script files aren't case sensitive so `tx` is like `Tx` for example.

old:

```go
SetGridSize(256,256,10)
SetCellSize(1e-9,1e-9,1e-9)
SetPBC(32,32,0)
```

new:

```go
Nx = 256
Ny = 256
Nz = 10
dx = 1e-9
dy = 1e-9
dz = 1e-9
PBCx = 32 // Optionnal
PBcy = 32 // Optionnal
PBCz = 0 // Optionnal
```

new (alternative but equivalent):

```go
Nx = 256
Ny = 256
Nz = 10
Tx = 256e-9
Ty = 256e-9
Tz = 10e-9
PBCx = 32 // Optionnal
PBcy = 32 // Optionnal
```

### You can add metadata

```go
Metadata("lattice_constant",500e-9)
Metadata("ref paper","X et al. (2023)")
```

You can access it in the file `.zattrs`. Or using [pyzfn](https://github.com/MathieuMoalic/pyzfn):

```python
print(job.lattice_constant)
print(job["ref paper"])
```
### You can save data by chunks

```go
Nx = 16
Ny = 32
Nz = 10
Tx = 16e-9
Ty = 32e-9
Tz = 10e-9
SetGeom(Universe())
sampling_interval = 5e-12
AutoSaveAsChunk(m,"m_chunked", sampling_interval, Chunk(1, 1, Nz, 3))
Run(1e-9)
```
This code will save the magnetization as chunks: 
  - The `x` and `y` dimensions are unchunked ( set as 1 in Chunk() )
  - The `z` dimension will have as many chunks as cells across the thickness ( 10 in  this case )
  - `mz`, `my` and `mz` will be also chunked, saved separately.
Why would you want it? Because it makes loading the data from disk MUCH faster if you chunk in a smart way.
Say you want to calculate the FFT of the top layer of the `y` component of the magnetization, in this case, loading the data from disk will be around 30 times faster.

### Other changes

- Remove the Google trackers in the GUI.
- Add saving as zarr
- Rename the functions to save `.ovf` files as `SaveOvf` and `SaveOvfAs`. You cannot autosave ovf files anymore.
- Add progress bar for `run`, can be turned off with `-magnets=false`
- Reorder GUI elements
- Dark mode GUI
- Check and warns the user for unoptimized mesh
- `AutoMeshx = True`,`AutoMeshy = True`and `AutoMeshz = True` will optimize the corresponding mesh axis for you (this function slightly changes the size and number of cells while keeping the total size of the system constant)
- Add chunking support as per the [zarr](https://zarr.readthedocs.io/en/stable/) documentation with the functions:
  - `SaveAsChunk(q Quantity, name string, rchunks RequestedChunking)`
  - `AutoSaveAsChunk(q Quantity, name string, period float64, rchunks RequestedChunking)`
- `Chunk(x, y, z, c int) -> RequestedChunking` chunks must fit an integer number of times along the axes. The chunks will be modified to be valid and as closed as the chunks you requested
- Add the `ShapeFromRegion` function
- Add new shapes : `squircle`, `triangle`, `rtriangle`, `diamond` and `hexagon`
- Add the `AutoSaveAs` function
- Add the `Round` function from the math library
- Add metadata saving : root_path, start_time, dx, dy, dz, Nx, Ny, Nz, Tx, Ty, Tz, StartTime, EndTime, TotalTime, PBC, Gpu, Host
- Add `MinimizeMaxSteps` and `MinimizeMaxTimeSeconds` to timeout long `Minimize()` runs.
- Everytime the function `Save` is used (from `AutoSave` for example), the current simulation time `t` is saved too as a zarray attribute
- Save compressed arrays (zstd) by default
- `ext_makegrains` now also takes a new argument `minRegion`. ext_makegrains(grainsize, minRegion, maxRegion, seed)
- Add colors for terminal logs

## Building from source

Using `podman` or `docker`:

```bash
git clone https://github.com/MathieuMoalic/amumax
cd amumax
podman build -t matmoa/amumax:build .
podman run --rm -v $PWD:/src matmoa/amumax:build bash cuda/build_cuda.sh
podman run --rm -v $PWD:/src matmoa/amumax:build
./build/amumax -v
```

The amumax binary and cuda libraries are then found in `build`.

## Contribution

I'm happy to consider any feature request. Don't hesitate to submit issues or PRs.

## TODO
- Calculate the mesh when enough arguments are given instead of when `SetGeom` is called
