
export interface Header {
    path: string;
    status: string;

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
// type TablePlotData struct {
// 	X float64 `json:"x"`
// 	Y float64 `json:"y"`
// }
// type TablePlot struct {
// 	AutoSaveInterval float64         `json:"autoSaveInterval"`
// 	Columns          []string        `json:"columns"`
// 	XColumn          string          `json:"xColumn"`
// 	YColumn          string          `json:"yColumn"`
// 	Data             []TablePlotData `json:"data"`
// }
export interface TablePlot {
    autoSaveInterval: number;
    columns: string[];
    xColumn: string;
    yColumn: string;
    data: TablePlotData[];
}
export interface TablePlotData {
    x: number;
    y: number;
}

export interface EngineState {
    header: Header;
    solver: Solver;
    console: Console;
    mesh: Mesh;
    parameters: Parameters;
    tablePlot: TablePlot;
}
