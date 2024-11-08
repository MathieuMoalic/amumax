import { previewState } from '$api/incoming/preview';
import { post } from '$api/post';
import { get } from 'svelte/store';

export function postComponent(component: string) {
	post('preview/component', { component });
}

export function postQuantity(quantity: string) {
	post('preview/quantity', { quantity });
}

export function postLayer(layer: number) {
	post('preview/layer', { layer });
}

export function postXChosenSize(xChosenSize: number) {
	post('preview/XChosenSize', { xChosenSize });
}
export function postYChosenSize(yChosenSize: number) {
	post('preview/YChosenSize', { yChosenSize });
}

export function postRefresh() {
	post('preview/refresh', {});
}
