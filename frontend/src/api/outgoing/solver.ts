import { writable } from "svelte/store";
import { get } from "svelte/store";
import { baseURL } from "$api/websocket";

export interface Solver {
    solver: string;
    runTime: number;
    runSteps: number;
    fixdt: number;
    mindt: number;
    maxdt: number;
    maxerr: number;
}

export const solver = writable<Solver>({
    solver: "",
    runTime: 0,
    runSteps: 0,
    fixdt: 0,
    mindt: 0,
    maxdt: 0,
    maxerr: 0
});

export async function postSolver() {
    let t = get(solver);
    let command = t.solver;
    const response = await fetch(`${get(baseURL)}/solver`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ command })
    });
    if (!response.ok) {
        throw new Error('Failed to post solver');
    }
    let parsedData: Solver = await response.json();
    solver.set(parsedData);
}