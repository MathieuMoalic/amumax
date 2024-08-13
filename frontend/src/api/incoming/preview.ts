import { writable } from "svelte/store";

export type VectorField = Array<{ x: number; y: number; z: number }>;
export type ScalarField = Array<Array<number>>;

export interface Preview {
    quantity: string;
    unit: string;
    component: number;
    layer: number;
    maxPoints: number;
    dimensions: [number, number, number];
    type: string;
    vectorFieldValues: VectorField;
    vectorFieldPositions: VectorField;
    scalarField: ScalarField;
    min: number;
    max: number;
    refresh: boolean;
}


export const previewState = writable<Preview>(
    {
        quantity: "",
        unit: "",
        component: 0,
        layer: 0,
        maxPoints: 0,
        dimensions: [0, 0, 0],
        type: "",
        vectorFieldValues: [],
        vectorFieldPositions: [],
        scalarField: [],
        min: 0,
        max: 0,
        refresh: false
    }
);