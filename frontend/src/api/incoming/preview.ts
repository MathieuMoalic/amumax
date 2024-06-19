import { writable } from 'svelte/store';

export type VectorField = Array<{ x: number; y: number; z: number }>;

export const displayData = writable<VectorField>([]);
