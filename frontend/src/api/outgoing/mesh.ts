import { get } from "svelte/store";
import { post } from "$api/post";
import { meshState } from "$api/incoming/mesh";

export function postMesh() {
    post('mesh', { mesh: get(meshState) });
}