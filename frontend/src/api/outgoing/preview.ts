import { previewState } from "$api/incoming/preview";
import { post } from "$api/post";
import { clearDisplay, display } from "$lib/preview/plot-vector-field";
import { get } from "svelte/store";

export function postComponent() {
    post('preview-component', { component: get(previewState).component });
}

export function postQuantity() {
    post('preview-quantity', { quantity: get(previewState).quantity });
}

export function postLayer() {
    post('preview-layer', { layer: get(previewState).layer });
}

export function postMaxPoints() {
    clearDisplay();
    post('preview-maxpoints', { maxPoints: get(previewState).maxPoints });
}

export function postRefresh() {
    post('preview-refresh', {});
}