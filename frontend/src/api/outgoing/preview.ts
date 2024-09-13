import { previewState } from '$api/incoming/preview';
import { post } from '$api/post';
import { get } from 'svelte/store';

export function postComponent(component: string) {
	post('preview/component', { component });
}

export function postQuantity(quantity: string) {
	post('preview/quantity', { quantity });
}

export function postLayer() {
	post('preview/layer', { layer: get(previewState).layer });
}

export function postMaxPoints(maxPoints: string) {
	post('preview/maxpoints', { maxPoints });
}

export function postRefresh() {
	post('preview/refresh', {});
}
