import { type EngineState } from "./types";
import { get } from "svelte/store";
import { headerState, consoleState, solverState, meshState, parametersState, tablePlotState, paused, refreshImage, refreshTablePlot } from "./api";

let ws: WebSocket;
let messages: any[] = [];
let input = "";

// Connect to WebSocket
export function wsConnect() {
    ws = new WebSocket('ws://localhost:5001/ws');

    ws.onopen = () => {
        console.log('Connected to the server');
    };

    ws.onerror = (error) => {
        // console.error('WebSocket Error:', error);
    };

    ws.onmessage = (e) => {
        messages = [...messages, e.data];
        parseMessage(e.data);
        refreshImage();
        refreshTablePlot();
    };

    ws.onclose = (e) => {
        console.log('Socket is closed. Reconnect will be attempted in 1 second.', e.reason);
        setTimeout(() => {
            wsConnect();
        }, 1000);
    };
}

function parseMessage(message: string) {
    let engineState: EngineState = JSON.parse(message);
    headerState.set(engineState.header);
    solverState.set(engineState.solver);
    consoleState.set(engineState.console);
    meshState.set(engineState.mesh);
    parametersState.set(engineState.parameters);
    tablePlotState.set(engineState.tablePlot);
    paused.set(get(headerState).status === 'Paused');
}