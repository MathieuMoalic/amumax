import { writable } from "svelte/store";
import { get } from "svelte/store";
import { baseURL } from "$api/websocket";

export interface Parameters {
    region: number;
}

export const parameters = writable<Parameters>({
    region: 0,
});

export async function postParameters() {
    let t = get(parameters);
    let region = t.region;
    const response = await fetch(`${get(baseURL)}/parameters`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ region })
    });
    if (!response.ok) {
        throw new Error('Failed to post parameters');
    }
    let parsedData: Parameters = await response.json();
    parameters.set(parsedData);
}