import { previewState } from "$api/incoming/preview";
import { post } from "$api/post";
import { get } from "svelte/store";

export function postComponent() {
    post('preview-component', get(previewState).component);
}

export function postQuantity() {
    post('preview-quantity', get(previewState).quantity);
}

export function postLayer() {
    post('preview-layer', get(previewState).layer);
}

export function postMaxPoints() {
    post('preview-maxpoints', get(previewState).maxPoints);
}

export function postRefresh() {
    post('preview-refresh', {});
}