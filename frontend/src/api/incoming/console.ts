import { writable } from "svelte/store";

export interface Console {
    hist: string;
}

export const consoleState = writable<Console>({ hist: '' });