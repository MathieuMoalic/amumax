import { get } from "svelte/store";
import { post } from "$api/post";

import { tablePlotState } from "$api/incoming/table-plot";

export function postXColumn() {
    post('tableplot/xcolumn', { XColumn: get(tablePlotState).xColumn });
}

export function postYColumn() {
    post('tableplot/ycolumn', { YColumn: get(tablePlotState).yColumn });
}

export function postAutoSaveInterval() {
    post('tableplot/autosave-interval', { autoSaveInterval: get(tablePlotState).autoSaveInterval });
}