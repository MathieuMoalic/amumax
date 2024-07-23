import { previewState } from "$api/incoming/preview";
import { get } from "svelte/store";
import { plotVectorField, disposeThreeJS } from "./plot-vector-field";
import { plotScalarField, disposeECharts } from "./plot-scalar-field";

let previousType = "";
let previousDimensions: number[];

function arraysEqual(a: number[], b: number[]) {
    if (a === b) return true;
    if (a == null || b == null) return false;
    if (a.length !== b.length) return false;

    for (let i = 0; i < a.length; i++) {
        if (a[i] !== b[i]) return false;
    }
    return true;
}

function newDimensions() {
    let dims = get(previewState).dimensions;
    if (previousDimensions === undefined) {
        previousDimensions = dims;
        return true;
    }
    if (!arraysEqual(dims, previousDimensions)) {
        previousDimensions = dims;
        return true;
    }
    return false;
}

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
        plotVectorField(newDimensions());
    } else if (type === 'scalar') {
        plotScalarField(newDimensions());
    } else {
        console.log('Unknown field type');
    }
}