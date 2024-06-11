import { writable, get } from "svelte/store";
import { baseURL } from "../websocket";

export interface Solver {
    type: string;
    steps: number;
    time: number;
    dt: number;
    errPerStep: number;
    maxTorque: number;
    fixdt: number;
    mindt: number;
    maxdt: number;
    maxerr: number;
}

export const solverState = writable<Solver>({
    type: "",
    steps: 0,
    time: 0,
    dt: 0,
    errPerStep: 0,
    maxTorque: 0,
    fixdt: 0,
    mindt: 0,
    maxdt: 0,
    maxerr: 0
});

export async function postSolverType(type: string) {
    const response = await fetch(`${get(baseURL)}/solver`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ type })
    });
    if (!response.ok) {
        throw new Error('Failed to post solver type');
    }
}

export async function postRun(duration: number) {
    const response = await fetch(`${get(baseURL)}/run`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ duration })
    });
    if (!response.ok) {
        throw new Error(response.statusText);
    }
}
export async function postSteps(steps: number) {
    const response = await fetch(`${get(baseURL)}/steps`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ steps })
    });
    if (!response.ok) {
        throw new Error('Failed to post steps');
    }
}
export async function postRelax() {
    const response = await fetch(`${get(baseURL)}/relax`, {
        method: 'POST'
    });
    if (!response.ok) {
        throw new Error('Failed to post relax');
    }
}
export async function postBreak() {
    const response = await fetch(`${get(baseURL)}/break`, {
        method: 'POST'
    });
    if (!response.ok) {
        throw new Error('Failed to post break');
    }
}