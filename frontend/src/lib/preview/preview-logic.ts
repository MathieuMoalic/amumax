import { previewState } from "$api/incoming/preview";
import { get } from "svelte/store";
import { plotVectorField } from "./plot-vector-field";
import { plotScalarField } from "./plot-scalar-field";

export function plotPreview() {
    console.log(get(previewState));
    let type = get(previewState).type;
    if (type === 'vector') {
        console.log('Plotting vector field');
        plotVectorField();
    } else if (type === 'scalar') {
        console.log('Plotting scalar field');
        plotScalarField();
    } else {
        console.log('Unknown field type');
    }
}