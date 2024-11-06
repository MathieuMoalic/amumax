import { writable } from "svelte/store";

export interface Header {
    path: string;
    status: string;
    version: string;
}

export const headerState = writable({
    path: '',
    status: '',
    version: ''
});