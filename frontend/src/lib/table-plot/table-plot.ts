import * as echarts from 'echarts';
import { get } from 'svelte/store';
import { tablePlotState } from '$api/incoming/table-plot';

export function plotTable() {
    if (chartInstance === undefined) {
        init();
    } else {
        update();
    }
}

let chartInstance: echarts.ECharts;

function update() {
    let t = get(tablePlotState);
    chartInstance.setOption({
        series: [
            {
                name: 'y',
                data: t.data,
            }
        ],
        xAxis: {
            name: `${t.xColumn} (${t.xColumnUnit})`,
            min: t.xmin,
            max: t.xmax,
        },
        yAxis: {
            name: `${t.yColumn} (${t.yColumnUnit})`,
            min: t.ymin,
            max: t.ymax,
        }
    });
}

export function init() {
    var chartDom = document.getElementById('table-plot')!;
    if (chartDom === null) {
        return;
    }
    // https://apache.github.io/echarts-handbook/en/best-practices/canvas-vs-svg
    chartInstance = echarts.init(chartDom, undefined, { renderer: 'svg' });
    let t = get(tablePlotState);

    // @ts-ignore
    chartInstance.setOption({
        axisPointer: {
            show: true,
            type: 'line',
            lineStyle: {
                color: '#6e9bcb',
                width: 2,
                type: 'dashed'
            },

            label: {
                backgroundColor: '#282a36',
                color: '#fff',
                formatter: function (params: any) {
                    return parseFloat(params.value).toPrecision(2);
                },
                padding: [8, 5, 8, 5],
                borderColor: '#6e9bcb',
                borderWidth: 1,
            }
        },
        animation: true,
        grid: {
            containLabel: false,
            left: '10%',
            right: '10%',
        },
        xAxis: {
            name: `${t.xColumn} (${t.xColumnUnit})`,
            min: t.xmin,
            max: t.xmax,
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
                    return parseFloat(value).toPrecision(2);
                },
                color: '#fff',
                // showMinLabel: true,
            }
        },
        yAxis: {
            name: `${t.yColumn} (${t.yColumnUnit})`,
            min: t.ymin,
            max: t.ymax,
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
                    return parseFloat(value).toPrecision(2);
                },
                // formatter: "{value}",
                color: '#fff',
                showMinLabel: true,
            }
        },
        series: [
            {
                type: 'line',
                name: 'y',
                showSymbol: false,
                data: t.data,
            }
        ]
    });
}

export function resizeECharts() {
    window.addEventListener('resize', function () {
        chartInstance.resize();
    });
}