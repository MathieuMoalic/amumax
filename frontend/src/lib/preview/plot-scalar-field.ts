import { previewState } from '$api/incoming/preview';
import * as echarts from 'echarts';
import { get } from 'svelte/store';
import { meshState } from '$api/incoming/mesh';

export function plotScalarField(newDimensions: boolean) {
    if (chart === undefined) {
        init();
    } else if (chart.isDisposed()) {
        init();
    } else if (newDimensions) {
        disposeECharts();
        init();
    } else {
        update();
    }
}

let chart: echarts.ECharts;
let previousDimension: number[];

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
            if (v === 0) { v = NaN; }
            data.push([i, j, v]);
        }
    }
    return data;
}

function update() {
    let data = getScalarField();
    let ps = get(previewState);
    chart.setOption({
        tooltip: {
            formatter: function (params: any) {
                if (params.value === undefined) {
                    return "NaN";
                }
                return `${params.value[2]} ${ps.unit}`;
            },
        },
        dataset: {
            source: data
        },
        visualMap: [{
            max: ps.max,
            min: ps.min,
        }]
    });
}

function init() {
    var chartDom = document.getElementById('container')!;
    // https://apache.github.io/echarts-handbook/en/best-practices/canvas-vs-svg
    chart = echarts.init(chartDom, undefined, { renderer: 'svg' });
    let ps = get(previewState);
    let dims = ps.dimensions;
    previousDimension = dims;
    let mesh = get(meshState);
    let xData = Array.from({ length: dims[0] }, (_, i) => i * mesh.dx * 1e9);
    let yData = Array.from({ length: dims[1] }, (_, i) => i * mesh.dy * 1e9);

    let aspectRatio = dims[0] / dims[1]; // Calculate the aspect ratio

    let gridWidth, gridHeight;

    // Ensure the grid dimensions maintain the aspect ratio
    if (aspectRatio > 1) {
        gridWidth = '80%';
        gridHeight = `${80 / aspectRatio}%`;
    } else {
        gridWidth = `${80 * aspectRatio}%`;
        gridHeight = '80%';
    }

    let data = getScalarField();

    // @ts-ignore
    chart.setOption({
        tooltip: {
            position: "top",
            formatter: function (params: any) {
                if (params.value === undefined) {
                    return "NaN";
                }
                return `${params.value[2]} ${ps.unit}`;
            },
            backgroundColor: '#282a36',
            borderColor: '#6e9bcb',
            textStyle: {
                color: '#fff'
            }

        },
        axisPointer: {
            show: true,
            type: 'line',
            triggerEmphasis: false,
            lineStyle: {
                color: '#6e9bcb',
                width: 2,
                type: 'dashed'
            },
            label: {
                backgroundColor: '#282a36',
                color: '#fff',
                formatter: function (params: any) {
                    if (params.value === undefined) {
                        return "NaN";
                    }
                    return ` ${parseFloat(params.value).toFixed(0)} nm`;
                },
                padding: [8, 5, 8, 5],
                borderColor: '#6e9bcb',
                borderWidth: 1,
            }
        },
        xAxis: {
            type: 'category',
            data: xData,
            name: "x (nm)",
            nameLocation: 'middle',
            nameGap: 25,
            nameTextStyle: {
                color: '#fff',
            },
            axisTick: {
                alignWithLabel: true,
                length: 6,
                lineStyle: {
                    type: 'solid',
                    color: '#fff',
                },
            },
            axisLabel: {
                show: true,
                formatter: function (value: string, _index: string) {
                    return parseFloat(value).toFixed(0);
                },
                color: '#fff',
                showMinLabel: true,
            }
        },
        yAxis: {
            type: 'category',
            data: yData,
            name: "y (nm)",
            nameLocation: 'middle',
            nameGap: 45,
            nameTextStyle: {
                color: '#fff',
            },
            axisTick: {
                alignWithLabel: true,
                length: 6,
                lineStyle: {
                    type: 'solid',
                    color: '#fff',
                },
            },
            axisLabel: {
                show: true,
                formatter: function (value: string, _index: string) {
                    return parseFloat(value).toFixed(0);
                },
                color: '#fff',
                showMinLabel: true,
            }
        },
        visualMap: [{
            type: 'continuous',
            min: ps.min,
            max: ps.max,
            calculable: true,
            realtime: true,
            precision: 2,
            itemWidth: 9,
            itemHeight: 140,
            text: [ps.unit, ''],
            textStyle: {
                color: '#fff'
            },
            inRange: {
                color: [
                    '#313695',
                    '#ffffff',
                    '#a50026'
                ]
            },
            left: 'right',
        }],
        series: [
            {
                name: ps.quantity,
                type: 'heatmap',
                emphasis: {
                    itemStyle: {
                        borderColor: '#333',
                        borderWidth: 1
                    },
                },
                selectedMode: true,
                progressive: 0,
                animation: true
            }
        ],
        dataset: {
            source: data
        },
        grid: {
            containLabel: false,
            left: '10%',
            right: '10%',
        },
        toolbox: {
            show: true,
            itemSize: 20,
            iconStyle: {
                borderColor: '#fff'
            },
            feature: {
                dataZoom: {
                    xAxisIndex: 0,
                    yAxisIndex: 0,
                    brushStyle: {
                        color: '#282a3655',
                        borderColor: '#6e9bcb',
                        borderWidth: 2,
                    },
                },
                dataView: { show: false },
                restore: {
                    show: false
                },
                saveAsImage: {
                    type: 'png',
                    name: 'preview',
                }

            }
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

export function resizeECharts() {
    window.addEventListener('resize', function () {
        chart.resize();
    });
}
