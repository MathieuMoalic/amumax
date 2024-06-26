import { writable } from "svelte/store";
import { get } from "svelte/store";
import { baseURL } from "$api/post";

export interface InitialMag {
    mag: string;
}

export const initialMag = writable<InitialMag>({
    mag: "",
});

export async function postInitialMag() {
    let t = get(initialMag);
    let mag = t.mag;
    const response = await fetch(`${get(baseURL)}/initialmag`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ shape: mag })
    });
    if (!response.ok) {
        throw new Error('Failed to post initialmag');
    }
    let parsedData: InitialMag = await response.json();
    initialMag.set(parsedData);
}