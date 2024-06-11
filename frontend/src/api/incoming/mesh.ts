import { writable } from "svelte/store";

export interface Mesh {
    dx: number;
    dy: number;
    dz: number;
    Nx: number;
    Ny: number;
    Nz: number;
    Tx: number;
    Ty: number;
    Tz: number;
    PBCx: number;
    PBCy: number;
    PBCz: number;
}

export const meshState = writable<Mesh>({
    dx: 0,
    dy: 0,
    dz: 0,
    Nx: 0,
    Ny: 0,
    Nz: 0,
    Tx: 0,
    Ty: 0,
    Tz: 0,
    PBCx: 0,
    PBCy: 0,
    PBCz: 0
});