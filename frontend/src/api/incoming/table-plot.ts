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
