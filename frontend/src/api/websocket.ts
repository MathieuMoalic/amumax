import msgpack from 'msgpack-lite';
import { plotChart } from "../lib/table-plot/tablePlot";

import { type Preview, previewState } from "./incoming/preview";
import { type Header, headerState } from "./incoming/header";
import { type Solver, solverState } from "./incoming/solver";
import { type Console, consoleState } from "./incoming/console";
import { type Mesh, meshState } from "./incoming/mesh";
import { type Parameters, parametersState, sortFieldsByName } from "./incoming/parameters";
import { type TablePlot, tablePlotState } from "./incoming/table-plot";
import { plotPreview } from '$lib/preview/preview-logic';


export function initializeWebSocket() {
    let retryInterval = 1000;
    let ws: WebSocket | null = null;

    function connect() {
        ws = new WebSocket(`/ws`);
        ws.binaryType = 'arraybuffer';

        ws.onopen = function () {
            console.log('WebSocket connection established');
        };

        ws.onmessage = function (event) {
            parseMsgpack(event.data);
            ws?.send('ok');
        };

        ws.onclose = function () {
            console.log('WebSocket closed. Attempting to reconnect in ' + retryInterval / 1000 + ' seconds...');
            ws = null; // Ensure ws is set to null when it is closed
            setTimeout(connect, retryInterval);
        };

        ws.onerror = function (event) {
            console.error('WebSocket encountered error:', event);
            if (ws) {
                ws.close();
            }
        };
    }

    function tryConnect() {
        console.log('Attempting WebSocket connection...');
        try {
            connect();
        } catch (err) {
            console.error('WebSocket connection failed:', err, 'Retrying in ' + retryInterval / 1000 + ' seconds...');
            setTimeout(tryConnect, retryInterval);
        }
    }

    tryConnect();
}


export function parseMsgpack(data: ArrayBuffer) {
    const msg = msgpack.decode(new Uint8Array(data));
    consoleState.set(msg.console as Console);
    headerState.set(msg.header as Header);
    meshState.set(msg.mesh as Mesh);
    parametersState.set(msg.parameters as Parameters);
    sortFieldsByName();
    previewState.set(msg.preview as Preview);
    solverState.set(msg.solver as Solver);
    tablePlotState.set(msg.tablePlot as TablePlot);
    plotChart();
    plotPreview();
}