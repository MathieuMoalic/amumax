import { writable, get } from "svelte/store";
import { baseURL } from "../websocket";

export interface Console {
    hist: string;
}

export const consoleState = writable<Console>({ hist: '' });

export async function postCommand(command: string) {
    const response = await fetch(`${get(baseURL)}/console`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ command })
    });
    if (!response.ok) {
        throw new Error('Failed to post command');
    }
}
