import { writable } from "svelte/store";
import { get } from "svelte/store";
import { baseURL } from "$api/post";

export interface Console {
    command: string;
}

export const console = writable<Console>({
    command: "",
});

export async function postCommand() {
    let t = get(console);
    let command = t.command;
    const response = await fetch(`${get(baseURL)}/console`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ command })
    });
    if (!response.ok) {
        throw new Error('Failed to post console command');
    }
    let parsedData: Console = await response.json();
    console.set(parsedData);
}