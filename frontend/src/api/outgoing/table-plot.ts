import { writable } from "svelte/store";
import { get } from "svelte/store";
import { baseURL } from "$api/websocket";

export interface TablePlot {
    autoSaveInterval: number;
    xColumn: string;
    yColumn: string;
}

export const tablePlot = writable<TablePlot>({
    autoSaveInterval: 0,
    xColumn: 't',
    yColumn: 'mx',
});

export async function postTableColumns() {
    let t = get(tablePlot);
    let autoSaveInterval = t.autoSaveInterval;
    let XColumn = t.xColumn;
    let YColumn = t.yColumn;
    const response = await fetch(`${get(baseURL)}/tableplot`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ XColumn, YColumn, autoSaveInterval })
    });
    if (!response.ok) {
        throw new Error('Failed to post table columns');
    }
    let parsedData: TablePlot = await response.json();
    tablePlot.set(parsedData);
}