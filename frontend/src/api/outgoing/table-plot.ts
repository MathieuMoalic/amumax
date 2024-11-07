import { post } from '$api/post';

export function postXColumn(XColumn: string) {
	post('tableplot/xcolumn', { XColumn });
}

export function postYColumn(YColumn: string) {
	post('tableplot/ycolumn', { YColumn });
}

export function postAutoSaveInterval(autoSaveInterval: string) {
	post('tableplot/autosave-interval', { autoSaveInterval });
}

export function postMaxPoints(inputFieldValue: string) {
	post('tableplot/maxpoints', { maxPoints: parseInt(inputFieldValue) });
}

export function postStep(inputFieldValue: string) {
	post('tableplot/step', { step: parseFloat(inputFieldValue) });
}
