import { get, writable } from "svelte/store";

export const baseURL = writable('http://localhost:5001');

export async function post(endpoint: string, data: any) {
    const response = await fetch(`${get(baseURL)}/${endpoint}`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ ...data })
    });
    if (!response.ok) {
        throw new Error(`Failed to post ${endpoint}`);
    }
}