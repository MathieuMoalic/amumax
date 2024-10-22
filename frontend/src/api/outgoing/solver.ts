import { get } from "svelte/store";
import { post } from "$api/post";
import { solverState } from "$api/incoming/solver";

export function postSolverType(type: string) {
    post('solver/type', { type });
}

export function postRun(runtime: string) {
    post('solver/run', { runtime });
}

export function postSteps(steps: string) {
    post('solver/steps', { steps });
}

export function postRelax() {
    post('solver/relax', {});
}

export function postMinimize() {
    post('solver/minimize', {});
}

export function postBreak() {
    post('solver/break', {});
}

export function postFixdt() {
    post('solver/fixdt', { fixdt: get(solverState).fixdt });
}

export function postMindt() {
    post('solver/mindt', { mindt: get(solverState).mindt });
}

export function postMaxdt() {
    post('solver/maxdt', { maxdt: get(solverState).maxdt });
}

export function postMaxerr() {
    post('solver/maxerr', { maxerr: get(solverState).maxerr });
}


