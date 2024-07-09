import { previewState } from '$api/incoming/preview';
import * as echarts from 'echarts';
import { get } from 'svelte/store';

export function plotScalarField() {
    if (chart === undefined || chart.isDisposed()) {
        init();
    }
    update();
}

let chart: echarts.ECharts;

function getScalarField(): number[][] {
    // Convert Uint8Array to ArrayBuffer
    const buffer = get(previewState).buffer.buffer;
    if (buffer.byteLength % 4 !== 0) {
        throw new Error("Invalid buffer length. Buffer length must be a multiple of 4 bytes.");
    }
    const float32Array = new Float32Array(buffer);
    let dims = get(previewState).dimensions;
    const data: number[][] = [];
    for (let i = 0; i < dims[0]; i++) {
        for (let j = 0; j < dims[1]; j++) {
            let v = float32Array[j * dims[0] + i];
            data.push([i, j, v]);
        }
    }
    return data;
}

function update() {
    let data = getScalarField();
    chart.setOption({
        dataset: {
            source: data
        }
    });
}

function init() {
    var chartDom = document.getElementById('container')!;
    chart = echarts.init(chartDom);

    let dims = get(previewState).dimensions;
    let yData: number[] = Array.from({ length: dims[1] }, (_, i) => i);
    let xData: number[] = Array.from({ length: dims[1] }, (_, i) => i);

    let data = getScalarField();

    chart.setOption({
        tooltip: {},
        xAxis: {
            type: 'category',
            data: xData
        },
        yAxis: {
            type: 'category',
            data: yData
        },
        visualMap: {
            min: -1,
            max: 1,
            calculable: true,
            realtime: false,
            inRange: {
                color: [
                    '#313695',
                    '#4575b4',
                    '#74add1',
                    '#abd9e9',
                    '#e0f3f8',
                    '#ffffbf',
                    '#fee090',
                    '#fdae61',
                    '#f46d43',
                    '#d73027',
                    '#a50026'
                ]
            }
        },
        series: [
            {
                name: get(previewState).quantity,
                type: 'heatmap',
                emphasis: {
                    itemStyle: {
                        borderColor: '#333',
                        borderWidth: 1
                    }
                },
                progressive: 1000,
                animation: false
            }
        ],
        dataset: {
            source: data
        }
    });
}

export function disposeECharts() {
    const container = document.getElementById('container')!;
    if (container) {
        const echartsInstance = echarts.getInstanceByDom(container);
        if (echartsInstance) {
            echartsInstance.dispose();
        }
    }
}
