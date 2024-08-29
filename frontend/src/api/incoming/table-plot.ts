import { writable } from "svelte/store";

export interface TablePlot {
    autoSaveInterval: number;
    columns: string[];
    xColumn: string;
    yColumn: string;
    xColumnUnit: string;
    yColumnUnit: string;
    data: number[][];
    xmin: number;
    xmax: number;
    ymin: number;
    ymax: number;
    maxPoints: number;
    step: number;
}

export const tablePlotState = writable<TablePlot>({
    autoSaveInterval: 0,
    columns: [],
    xColumn: 't',
    yColumn: 'mx',
    xColumnUnit: 's',
    yColumnUnit: '',
    data: [],
    xmin: 0,
    xmax: 0,
    ymin: 0,
    ymax: 0,
    maxPoints: 0,
    step: 0,
});
