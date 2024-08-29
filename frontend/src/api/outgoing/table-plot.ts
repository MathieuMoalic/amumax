import { get } from "svelte/store";
import { post } from "$api/post";

import { tablePlotState } from "$api/incoming/table-plot";

export function postXColumn() {
    post('tableplot/xcolumn', { XColumn: get(tablePlotState).xColumn });
}

export function postYColumn() {
    post('tableplot/ycolumn', { YColumn: get(tablePlotState).yColumn });
}

export function postAutoSaveInterval(inputFieldValue: string) {
    post('tableplot/autosave-interval', { autoSaveInterval:parseFloat(inputFieldValue) });
}

export function postMaxPoints(inputFieldValue: string) {
    post('tableplot/maxpoints', { maxPoints: parseInt(inputFieldValue) });
}

export function postStep(inputFieldValue: string) {
    post('tableplot/step', { step: parseFloat(inputFieldValue) });
}