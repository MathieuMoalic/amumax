import { writable } from "svelte/store";
import { get } from "svelte/store";
import { baseURL } from "$api/websocket";

export interface Preview {
    Quantity: string;
    Component: number;
    Layer: number;
}

export const preview = writable<Preview>({
    Quantity: 'm',
    Component: 0,
    Layer: 0,
});

export async function postPreview() {
    let value = get(preview);
    const response = await fetch(`${get(baseURL)}/preview`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ value })
    });
    if (!response.ok) {
        throw new Error('Failed to post preview');
    }
}
