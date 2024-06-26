import { writable } from "svelte/store";

export interface Solver {
    type: string;
    steps: number;
    time: number;
    dt: number;
    errPerStep: number;
    maxTorque: number;
    fixdt: number;
    mindt: number;
    maxdt: number;
    maxerr: number;
}

export const solverState = writable<Solver>({
    type: "",
    steps: 0,
    time: 0,
    dt: 0,
    errPerStep: 0,
    maxTorque: 0,
    fixdt: 0,
    mindt: 0,
    maxdt: 0,
    maxerr: 0
});

