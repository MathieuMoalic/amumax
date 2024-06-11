import { writable, get } from "svelte/store";
import { baseURL } from "../websocket";

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

export const tablePlotState = writable<TablePlot>({
    autoSaveInterval: 0,
    columns: [],
    xColumn: 't',
    yColumn: 'mx',
    data: []
});

export async function postTableColumns() {
    let XColumn = get(tablePlotState).xColumn;
    let YColumn = get(tablePlotState).yColumn;
    const response = await fetch(`${get(baseURL)}/table`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ XColumn, YColumn })
    });
    if (!response.ok) {
        throw new Error('Failed to post table columns');
    }
    let parsedData: TablePlot = await response.json();
    tablePlotState.set(parsedData);
}