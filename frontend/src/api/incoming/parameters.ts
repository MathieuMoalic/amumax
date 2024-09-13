import { writable } from 'svelte/store';

export interface Field {
	name: string;
	value: string;
	description: string;
	changed: boolean;
}

export interface Parameters {
	regions: number[];
	fields: Field[];
	selectedRegion: number;
}

export const parametersState = writable<Parameters>({
	regions: [],
	fields: [],
	selectedRegion: 0
});

export function sortFieldsByName() {
	parametersState.update((parameters) => {
		parameters.fields.sort((a, b) => a.name.localeCompare(b.name));
		return parameters;
	});
}
