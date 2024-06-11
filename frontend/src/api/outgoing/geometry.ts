import { writable } from "svelte/store";
import { get } from "svelte/store";
import { baseURL } from "$api/websocket";

export interface Geometry {
    shape: string;
}

export const geometry = writable<Geometry>({
    shape: "",
});

export async function postGeometry() {
    let t = get(geometry);
    let shape = t.shape;
    const response = await fetch(`${get(baseURL)}/geometry`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ shape })
    });
    if (!response.ok) {
        throw new Error('Failed to post geometry');
    }
    let parsedData: Geometry = await response.json();
    geometry.set(parsedData);
}