import msgpack from 'msgpack-lite';
import { writable } from "svelte/store";
import { plotChart } from "../lib/table-plot/tablePlot";
import { display } from "$lib/preview/plot-vector-field"

import { type Preview, previewState } from "./incoming/preview";
import { type Header, headerState } from "./incoming/header";
import { type Solver, solverState } from "./incoming/solver";
import { type Console, consoleState } from "./incoming/console";
import { type Mesh, meshState } from "./incoming/mesh";
import { type Parameters, parametersState } from "./incoming/parameters";
import { type TablePlot, tablePlotState } from "./incoming/table-plot";
import { plotPreview } from '$lib/preview/preview-logic';

export let ws: WebSocket;
export function initializeWebSocket(ws: WebSocket) {
    ws = new WebSocket('ws://localhost:5001/ws');
    ws.binaryType = 'arraybuffer';

    ws.onmessage = function (event) {
        const startTime = performance.now(); // Start timing

        const msg = msgpack.decode(new Uint8Array(event.data));
        consoleState.set(msg.console as Console);
        headerState.set(msg.header as Header);
        meshState.set(msg.mesh as Mesh);
        parametersState.set(msg.parameters as Parameters);
        previewState.set(msg.preview as Preview);
        solverState.set(msg.solver as Solver);
        tablePlotState.set(msg.tablePlot as TablePlot);
        plotChart();
        plotPreview();

        const endTime = performance.now(); // End timing
        const parsingTime = endTime - startTime;
        display.update(currentDisplay => {
            if (currentDisplay) {
                return { ...currentDisplay, parsingTime: parsingTime };
            }
            return currentDisplay;
        });
    };
}

// function parseVectorField(arrayBuffer: ArrayBuffer): VectorField {
//     const float32Array = new Float32Array(arrayBuffer);
//     const vectors: VectorField = [];

//     for (let i = 0; i < float32Array.length; i += 3) {
//         vectors.push({
//             x: float32Array[i],
//             y: float32Array[i + 1],
//             z: float32Array[i + 2]
//         });
//     }

//     return vectors;
// }