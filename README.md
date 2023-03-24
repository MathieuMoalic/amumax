# Amumax
fork of [mumax3](https://github.com/mumax/3)

The solvers ( and results ) are unchanged, this is just list of massive quality of life improvements making working with the output data much more efficient and convenient.

## Changelog
- Add saving as zarr
- `Nx`,`Ny`,`Nz`,`dx`,`dy`,`dz` are now predefined variables. You define the Mesh through them. You can't redefine the Mesh.
- Remove support for OVF1, OVF2, dump, anything that's not zarr
- Add progress bar for `run`
- Reorder GUI elements
- Dark mode GUI
- Check and warns the user for unoptimized mesh
- `AutoMesh = True` will optimize the shape for you (this function slightly changes the size and number of cells while keeping the total size of the system constant)
- Add chunking support as per the [zarr](https://zarr.readthedocs.io/en/stable/) documentation with the functions:
    - `SaveAsChunk(q Quantity, name string, rchunks RequestedChunking)`
    - `AutoSaveAsChunk(q Quantity, name string, period float64, rchunks RequestedChunking)`
- `Chunk(x, y, z, c int) -> RequestedChunking` chunks must fix an integer number of times along the axes. The chunks will be modified to be valid and as closed as the chunks you requested
- The graph plot in the GUI is probably broken
- Add the `ShapeFromRegion` function
- Add new shapes : `squircle`, `triangle`, `rtriangle`, `diamond` and `hexagon`
- `m` is not added by default in the table anymore, only `t` is
- Add the `autosaveas` function
- Add the `round` function from the math library
- Add metadata saving : root_path, start_time, Dx, Dy, Dz, Nx, Ny, Nz, Tx, Ty, Tz, StartTime, EndTime, TotalTime, PBC, Gpu, Host
- Everytime the function `save` is used (from `autosave` for example), the current simulation time `t` is saved too as a zarray attribute
- Save compressed arrays (zstd) by default

# TODO
 - Fix tablesave without tableadd before
 - Re-add the plot
