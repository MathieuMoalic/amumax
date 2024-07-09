import { writable } from "svelte/store";

export type VectorField = Array<{ x: number; y: number; z: number }>;
export type ScalarField = Array<number>;

export interface Preview {
    quantity: string;
    component: string;
    layer: number;
    maxPoints: number;
    dimensions: [number, number, number];
    type: string;
    buffer: Uint8Array;
}

export const previewState = writable<Preview>({
    quantity: 'm',
    component: 'All',
    layer: 0,
    maxPoints: 1000,
    dimensions: [0, 0, 0],
    type: 'vector',
    buffer: new Uint8Array(0),
});