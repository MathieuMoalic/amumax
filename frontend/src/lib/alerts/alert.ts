// svelte store that manages the alerts
import { writable } from 'svelte/store';

export let alert = writable("");

export function setAlert(message: string) {
    alert.set(message);
    setTimeout(() => {
        alert.set("");
    }, 10000);
}