import msgpack from 'msgpack-lite';

import { type Preview, previewState } from './incoming/preview';
import { type Header, headerState } from './incoming/header';
import { type Solver, solverState } from './incoming/solver';
import { type Console, consoleState } from './incoming/console';
import { type Mesh, meshState } from './incoming/mesh';
import { type Parameters, parametersState, sortFieldsByName } from './incoming/parameters';
import { type TablePlot, tablePlotState } from './incoming/table-plot';
import { get, writable } from 'svelte/store';
import { preview3D } from '$lib/preview/preview3D';
import { preview2D } from '$lib/preview/preview2D';
import { plotTable } from '$lib/table-plot/table-plot';

export let connected = writable(false);
let wsUrl = import.meta.env.VITE_WS_URL || '/ws';

export function initializeWebSocket() {
	let retryInterval = 1000;
	let ws: WebSocket | null = null;

	function connect() {
		console.debug('Connecting to WebSocket server at', wsUrl);
		ws = new WebSocket(wsUrl);
		ws.binaryType = 'arraybuffer';

		ws.onopen = function () {
			console.debug('WebSocket connection established');
			connected.set(true);
		};

		ws.onmessage = function (event) {
			parseMsgpack(event.data);
			ws?.send('ok');
			connected.set(true);
		};

		ws.onclose = function () {
			connected.set(false);
			console.debug(
				'WebSocket closed. Attempting to reconnect in ' + retryInterval / 1000 + ' seconds...'
			);
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
		console.debug('Attempting WebSocket connection...');
		try {
			connect();
		} catch (err) {
			console.error(
				'WebSocket connection failed:',
				err,
				'Retrying in ' + retryInterval / 1000 + ' seconds...'
			);
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

	solverState.set(msg.solver as Solver);

	tablePlotState.set(msg.tablePlot as TablePlot);
	plotTable();

	previewState.set(msg.preview as Preview);
	if (get(previewState).type === '3D') {
		preview3D();
	} else {
		preview2D();
	}
}
