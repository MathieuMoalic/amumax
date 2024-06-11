import { writable } from "svelte/store";

export interface Field {
    name: string;
    value: string;
    description: string;
}

export interface Parameters {
    regions: number[];
    fields: Field[];
}

export const parametersState = writable<Parameters>({
    regions: [],
    fields: []
});