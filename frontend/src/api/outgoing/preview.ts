import { writable } from "svelte/store";
import { get } from "svelte/store";
import { baseURL } from "$api/websocket";

export interface Preview {
    quantity: string;
    component: string;
    layer: number;
    maxPoints: number;
}
export const preview = writable<Preview>({
    quantity: 'm',
    component: 'All',
    layer: 0,
    maxPoints: 1000,
});

export async function postPreview() {

    let data = get(preview);
    const response = await fetch(`${get(baseURL)}/preview`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ ...data })
    });
    if (!response.ok) {
        throw new Error('Failed to post preview');
    }
}
