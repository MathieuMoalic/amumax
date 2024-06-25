import msgpack from 'msgpack-lite';
import { get, writable } from "svelte/store";
import { redraw } from "../lib/table-plot/tablePlot";
import { update, init, display } from "$lib/preview/plot-vector-field"
import { displayData, type VectorField } from "$api/incoming/preview";

import { type Header, headerState } from "./incoming/header";
import { type Solver, solverState } from "./incoming/solver";
import { type Console, consoleState } from "./incoming/console";
import { type Mesh, meshState } from "./incoming/mesh";
import { type Parameters, parametersState } from "./incoming/parameters";
import { type TablePlot, tablePlotState } from "./incoming/table-plot";

export const baseURL = writable('http://localhost:5001');
export let ws: WebSocket;

export interface EngineState {
    header: Header;
    solver: Solver;
    console: Console;
    mesh: Mesh;
    parameters: Parameters;
    tablePlot: TablePlot;
}

export function initializeWebSocket(ws: WebSocket) {
    ws = new WebSocket('ws://localhost:5001/ws');
    ws.binaryType = 'arraybuffer';

    ws.onmessage = function (event) {
        const startTime = performance.now(); // Start timing
        const msg = msgpack.decode(new Uint8Array(event.data));

        const engineState = msg.engine_state as EngineState;
        headerState.set(engineState.header);
        solverState.set(engineState.solver);
        consoleState.set(engineState.console);
        meshState.set(engineState.mesh);
        parametersState.set(engineState.parameters);
        tablePlotState.set(engineState.tablePlot);
        redraw();

        if (get(displayData).length === 0) {
            displayData.set(parseVectorField(msg.preview_buffer))
            init()
        } else {
            displayData.set(parseVectorField(msg.preview_buffer))
            update();
        }
        console.log(engineState, get(displayData));
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

function parseVectorField(arrayBuffer: ArrayBuffer): VectorField {
    const float32Array = new Float32Array(arrayBuffer);
    const vectors: VectorField = [];

    for (let i = 0; i < float32Array.length; i += 3) {
        vectors.push({
            x: float32Array[i],
            y: float32Array[i + 1],
            z: float32Array[i + 2]
        });
    }

    return vectors;
}