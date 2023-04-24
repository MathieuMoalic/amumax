# Amumax
fork of [mumax3](https://github.com/mumax/3) meant to increase the integration with a python processing workflow. I made my own wrapper around zarr called [pyzfn](https://github.com/MathieuMoalic/pyzfn) which leverages the mumax data in the zarr format.

The solvers ( and results ) are unchanged, this is just list of massive quality of life improvements making working with the output data much more efficient and convenient.

It's not 100% compatible with the original `.mx3` files. See changes below.

## Installation
### Linux
#### Install script
Don't just run an script on the internet. Read it, check what it does and then you can run this command to install amumax: 

```sh -c "$(curl -fsSL https://raw.githubusercontent.com/MathieuMoalic/amumax/main/install.sh)```

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
`Nx`,`Ny`,`Nz`,`dx`,`dy`,`dz`,`PBCx`,`PBCy`,`PBCz` are now predefined variables. You define the Mesh through them. You don't need to call a function to initiate the Mesh, it is automatically done the first time you run a solver but you can't redefine the Mesh after that !


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
PBCx = 32
PBcy = 32
PBCz = 0
```

### Other changes
- Remove the Google trackers in the GUI.
- Add saving as zarr
- Mostly remove support for OVF1, OVF2, dump, anything that's not zarr.
- Add progress bar for `run`
- Reorder GUI elements
- Dark mode GUI
- Check and warns the user for unoptimized mesh
- `AutoMesh = True` will optimize the shape for you (this function slightly changes the size and number of cells while keeping the total size of the system constant)
- Add chunking support as per the [zarr](https://zarr.readthedocs.io/en/stable/) documentation with the functions:
    - `SaveAsChunk(q Quantity, name string, rchunks RequestedChunking)`
    - `AutoSaveAsChunk(q Quantity, name string, period float64, rchunks RequestedChunking)`
- `Chunk(x, y, z, c int) -> RequestedChunking` chunks must fit an integer number of times along the axes. The chunks will be modified to be valid and as closed as the chunks you requested
- Add the `ShapeFromRegion` function
- Add new shapes : `squircle`, `triangle`, `rtriangle`, `diamond` and `hexagon`
- Add the `AutoSaveAs` function
- Add the `Round` function from the math library
- Add metadata saving : root_path, start_time, dx, dy, dz, Nx, Ny, Nz, Tx, Ty, Tz, StartTime, EndTime, TotalTime, PBC, Gpu, Host
- Everytime the function `Save` is used (from `AutoSave` for example), the current simulation time `t` is saved too as a zarray attribute
- Save compressed arrays (zstd) by default
- `ext_makegrains` now also takes a new argument `minRegion`. ext_makegrains(grainsize, minRegion, maxRegion, seed)

## Contribution
I'm happy to consider any feature request. Don't hesitate to submit issues or PRs.
