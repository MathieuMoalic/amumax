import { get } from 'svelte/store';
import { post } from '$api/post';

import { tablePlotState } from '$api/incoming/table-plot';

export function postXColumn(XColumn: string) {
	post('tableplot/xcolumn', { XColumn });
}

export function postYColumn(YColumn: string) {
	post('tableplot/ycolumn', { YColumn });
}

export function postAutoSaveInterval(inputFieldValue: string) {
	post('tableplot/autosave-interval', { autoSaveInterval: parseFloat(inputFieldValue) });
}

export function postMaxPoints(inputFieldValue: string) {
	post('tableplot/maxpoints', { maxPoints: parseInt(inputFieldValue) });
}

export function postStep(inputFieldValue: string) {
	post('tableplot/step', { step: parseFloat(inputFieldValue) });
}
