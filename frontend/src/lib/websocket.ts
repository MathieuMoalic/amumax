import { type EngineState } from "./types";
import { get } from "svelte/store";
import { headerState, consoleState, solverState, meshState, parametersState, tablePlotState, paused, displayData } from "./api";
import { redraw } from "./tablePlot";
import { update, parsingTime, init } from "./display";
import { type VectorField } from "./types";

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
    const startTime = performance.now(); // Start timing

    const arrayBuffer = await blobToArrayBuffer(data);
    displayData.set(parseVectorField(arrayBuffer));

    const endTime = performance.now(); // End timing
    parsingTime.set(endTime - startTime);
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
    paused.set(get(headerState).status === 'Paused');
    redraw()
}