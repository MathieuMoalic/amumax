# Amumax
Personnal fork of [mumax3](https://github.com/mumax/3)

The solvers are unchanged, this is just list of massive quality of life improvements making working with the output data much more efficient and convenient.

## Changelog
- Add saving as zarr
- Remove support for OVF1, OVF2, dump, anything that's not zarr
- Add progress bar for `run`
- Reorder GUI elements
- Dark mode GUI
- Check and warns the user for unoptimized mesh
- `autokernel` will optimize the shape for you (this function slightly changes the size and number of cells while keeping the total size of the system constant)
- Add chunking support as per the [zarr](https://zarr.readthedocs.io/en/stable/) documentation with the function `chunkxyzc`
- The graph plot in the GUI is probably broken
- Add the `ShapeFromRegion` function
- Add new shapes : `squircle`, `triangle`, `rtriangle`, `diamond` and `hexagon`
- `m` is not added by default in the table anymore, only `t` is
- Add the `autosaveas` function
- Add the `round` function from the math library
- Add metadata saving : root_path, start_time, Dx, Dy, Dz, Nx, Ny, Nz, Tx, Ty, Tz, StartTime, EndTime, TotalTime, PBC, Gpu, Host
- Everytime the function `save` is used (from `autosave` for example), the current simulation time `t` is saved too as a zarray attribute
- Save compressed arrays (zstd) by default
