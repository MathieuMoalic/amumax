import { previewState } from "$api/incoming/preview";
import { get } from "svelte/store";
import { plotVectorField, disposeThreeJS } from "./plot-vector-field";
import { plotScalarField, disposeECharts } from "./plot-scalar-field";

let previousType = "";

export function plotPreview() {
    let type = get(previewState).type;
    if (type !== previousType) {
        if (previousType === 'vector') {
            disposeThreeJS();
        } else if (previousType === 'scalar') {
            disposeECharts();
        }
    }
    previousType = type;
    if (type === 'vector') {
        plotVectorField();
    } else if (type === 'scalar') {
        plotScalarField();
    } else {
        console.log('Unknown field type');
    }
}