import { writable } from "svelte/store";

export interface TablePlot {
    autoSaveInterval: number;
    columns: string[];
    xColumn: string;
    yColumn: string;
    xColumnUnit: string;
    yColumnUnit: string;
    data: number[][];
    min: number;
    max: number;
}

export const tablePlotState = writable<TablePlot>({
    autoSaveInterval: 0,
    columns: [],
    xColumn: 't',
    yColumn: 'mx',
    xColumnUnit: 's',
    yColumnUnit: '',
    data: [],
    min: 0,
    max: 0
});
