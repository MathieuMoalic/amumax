import { get, writable } from "svelte/store";
import { redraw } from "../lib/table-plot/tablePlot";
import { update, init, display } from "./incoming/preview";

import { type Header, headerState } from "./incoming/header";
import { type Solver, solverState } from "./incoming/solver";
import { type Console, consoleState } from "./incoming/console";
import { type Mesh, meshState } from "./incoming/mesh";
import { type Parameters, parametersState } from "./incoming/parameters";
import { type TablePlot, tablePlotState } from "./incoming/table-plot";

export const displayData = writable<VectorField>([]);
export const baseURL = writable('http://localhost:5001');

export type VectorField = Array<{ x: number; y: number; z: number }>;

export interface EngineState {
    header: Header;
    solver: Solver;
    console: Console;
    mesh: Mesh;
    parameters: Parameters;
    tablePlot: TablePlot;
}

let ws: WebSocket;
export function wsConnect() {
    ws = new WebSocket('ws://localhost:5001/ws');

    ws.onopen = () => {
        // console.log('Connected to the server');
    };

    ws.onerror = (error) => {
        // console.error('WebSocket Error:', error);
    };

    ws.onmessage = async (e) => {
        const startTime = performance.now(); // Start timing

        if (typeof e.data === 'string') {
            try {
                parseJsonMessage(e.data);
            } catch (error) {
                console.error('Error parsing JSON message:', error);
            }
        } else if (e.data instanceof Blob) {
            try {
                if (get(displayData).length === 0) {
                    await parseBinaryMessage(e.data);
                    init()
                } else {
                    parseBinaryMessage(e.data);
                    update();
                }
            } catch (error) {
                console.error('Error parsing Binary message:', error);
            }
        }
        const endTime = performance.now(); // End timing
        const parsingTime = endTime - startTime;
        display.update(currentDisplay => {
            if (currentDisplay) {
                return { ...currentDisplay, parsingTime: parsingTime };
            }
            return currentDisplay;
        });
    };

    ws.onclose = (e) => {
        console.log('Socket is closed.');
        // console.log('Socket is closed. Reconnect will be attempted in 1 second.', e.reason);
        // setTimeout(() => {
        //     wsConnect();
        // }, 1000);
    };
}

async function parseBinaryMessage(data: Blob) {
    const arrayBuffer = await blobToArrayBuffer(data);
    displayData.set(parseVectorField(arrayBuffer));
}

function blobToArrayBuffer(blob: Blob): Promise<ArrayBuffer> {
    return new Promise((resolve, reject) => {
        const reader = new FileReader();
        reader.onload = () => resolve(reader.result as ArrayBuffer);
        reader.onerror = reject;
        reader.readAsArrayBuffer(blob);
    });
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

function parseJsonMessage(message: string) {
    let engineState: EngineState = JSON.parse(message);
    headerState.set(engineState.header);
    solverState.set(engineState.solver);
    consoleState.set(engineState.console);
    meshState.set(engineState.mesh);
    parametersState.set(engineState.parameters);
    tablePlotState.set(engineState.tablePlot);
    // paused.set(get(headerState).status === 'Paused');
    redraw()
}