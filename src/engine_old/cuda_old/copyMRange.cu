extern "C" __global__ void
copyMRange(float* __restrict__ dst,
          float* __restrict__ src,
          int Nx, int Ny, int Nz,          // full volume dims
          int dx0, int dy0, int dz0,       // dst origin
          int sx0, int sy0, int sz0,       // src origin
          int W, int H, int D,             // box size
          int wrap)                        // 0=clip, 1=wrap
{
    int x = blockIdx.x * blockDim.x + threadIdx.x; // [0..W)
    int y = blockIdx.y * blockDim.y + threadIdx.y; // [0..H)
    int z = blockIdx.z * blockDim.z + threadIdx.z; // [0..D)
    if (x >= W || y >= H || z >= D) return;

    int sx = sx0 + x, sy = sy0 + y, sz = sz0 + z;
    int dx = dx0 + x, dy = dy0 + y, dz = dz0 + z;

    if (wrap) {
        // positive modulo
        sx = (sx % Nx + Nx) % Nx;
        sy = (sy % Ny + Ny) % Ny;
        sz = (sz % Nz + Nz) % Nz;
        dx = (dx % Nx + Nx) % Nx;
        dy = (dy % Ny + Ny) % Ny;
        dz = (dz % Nz + Nz) % Nz;
    } else {
        if (sx < 0 || sy < 0 || sz < 0 || sx >= Nx || sy >= Ny || sz >= Nz) return;
        if (dx < 0 || dy < 0 || dz < 0 || dx >= Nx || dy >= Ny || dz >= Nz) return;
    }

    size_t sidx = ((size_t)sz * Ny + sy) * Nx + sx;
    size_t didx = ((size_t)dz * Ny + dy) * Nx + dx;

    dst[didx] = src[sidx];
}
