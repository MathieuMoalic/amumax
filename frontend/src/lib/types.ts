
export interface Header {
    Path: string;
    Progress: number;
    Status: string;

}
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
export interface Console {
    hist: string;
}
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

export interface Parameters {
    aex: number;
    alpha: number;
    anisC1: number;
    anisC2: number;
    anisU: number;
    b1: number;
    b2: number;
    bExt: number;
    dbulk: number;
    dind: number;
    epsilonPrime: number;
    exx: number;
    exy: number;
    exz: number;
    eyy: number;
    eyz: number;
    ezz: number;
    fixedLayer: number;
    freeLayerThickness: number;
    frozenspins: number;
    j: number;
    kc1: number;
    kc2: number;
    kc3: number;
    ku1: number;
    ku2: number;
    lambda: number;
    mfmDipole: number;
    mfmLift: number;
    msat: number;
    noDemagSpins: number;
    pol: number;
    temp: number;
    xi: number;
}

export interface TablePlot {
    autoSaveInterval: number;
    columns: string[];
}

export interface EngineState {
    header: Header;
    solver: Solver;
    console: Console;
    mesh: Mesh;
    parameters: Parameters;
    tablePlot: TablePlot;
}