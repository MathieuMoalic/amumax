
import { writable, get } from "svelte/store";
import type { Console, Solver, Mesh, Parameters, TablePlot } from "./types";

export async function refreshTablePlot() {
    let encodedQuery2 = encodeURIComponent(get(tablePlotx));
    let encodedQuery3 = encodeURIComponent(get(tablePloty));
    TablePlotUrl.set(`${get(baseURL)}/tableplot?x=${encodedQuery2}&y=${encodedQuery3}&` + Math.random())
}
export async function refreshImage() {
    let encodedQuery = encodeURIComponent(get(renderQuant));
    imageUrl.set(`${get(baseURL)}/image?quantity=${encodedQuery}&` + Math.random())
}

export const tablePlotx = writable('t');
export const tablePloty = writable('mx');
export const imageUrl = writable('http://localhost:5001/image?quantity=m');
export const TablePlotUrl = writable('http://localhost:5001/tableplot');
export const baseURL = writable('http://localhost:5001');
export const renderQuant = writable('m');
export const renderComp = writable('All');
export const renderLayer = writable(0);
export const paused = writable(false);

export const headerState = writable({
    Path: '',
    Progress: 0,
    Status: ''
});
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
export const consoleState = writable<Console>({ hist: '' });
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
export const parametersState = writable<Parameters>({
    aex: 0,
    alpha: 0,
    anisC1: 0,
    anisC2: 0,
    anisU: 0,
    b1: 0,
    b2: 0,
    bExt: 0,
    dbulk: 0,
    dind: 0,
    epsilonPrime: 0,
    exx: 0,
    exy: 0,
    exz: 0,
    eyy: 0,
    eyz: 0,
    ezz: 0,
    fixedLayer: 0,
    freeLayerThickness: 0,
    frozenspins: 0,
    j: 0,
    kc1: 0,
    kc2: 0,
    kc3: 0,
    ku1: 0,
    ku2: 0,
    lambda: 0,
    mfmDipole: 0,
    mfmLift: 0,
    msat: 0,
    noDemagSpins: 0,
    pol: 0,
    temp: 0,
    xi: 0
});
export const tablePlotState = writable<TablePlot>({
    autoSaveInterval: 0,
    columns: [],
});

export async function postSolverType(type: string) {
    const response = await fetch(`${get(baseURL)}/solver`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ type })
    });
    if (!response.ok) {
        throw new Error('Failed to post solver type');
    }
}
export async function postRun(duration: number) {
    const response = await fetch(`${get(baseURL)}/run`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ duration })
    });
    if (!response.ok) {
        throw new Error('Failed to post run');
    }
}
export async function postSteps(steps: number) {
    const response = await fetch(`${get(baseURL)}/steps`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ steps })
    });
    if (!response.ok) {
        throw new Error('Failed to post steps');
    }
}
export async function postRelax() {
    const response = await fetch(`${get(baseURL)}/relax`, {
        method: 'POST'
    });
    if (!response.ok) {
        throw new Error('Failed to post relax');
    }
}
export async function postBreak() {
    const response = await fetch(`${get(baseURL)}/break`, {
        method: 'POST'
    });
    if (!response.ok) {
        throw new Error('Failed to post break');
    }
}

export async function postCommand(command: string) {
    const response = await fetch(`${get(baseURL)}/console`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ command })
    });
    if (!response.ok) {
        throw new Error('Failed to post command');
    }
}