import { previewState } from '$api/incoming/preview';
import * as echarts from 'echarts';
import { get } from 'svelte/store';
import { meshState } from '$api/incoming/mesh';
import { disposePreview3D } from './preview3D';

export function preview2D() {
	if (get(previewState).scalarField === null) {
		disposePreview2D();
		disposePreview3D();
		return;
	}

	const container = document.getElementById('container')!;
	if (get(previewState).refresh || !container) {
		disposePreview2D();
		disposePreview3D();
		init();
	} else {
		update();
	}
}

function getColorMap(min: number, max: number) {
	if (min < 0 && max > 0) {
		return ['#313695', '#ffffff', '#a50026'];
	} else {
		return ['#ffffff', '#a50026'];
	}
}

let chartInstance: echarts.ECharts;

function update() {
	if (chartInstance === undefined || chartInstance.isDisposed()) {
		init();
		return;
	}
	let data = get(previewState).scalarField;
	let ps = get(previewState);
	chartInstance.setOption({
		tooltip: {
			formatter: function (params: any) {
				if (params.value === undefined) {
					return 'NaN';
				}
				return `${params.value[2]} ${ps.unit}`;
			}
		},
		dataset: {
			source: data
		},
		visualMap: [
			{
				max: ps.max,
				min: ps.min
			}
		]
	});
}

function init() {
	var chartDom = document.getElementById('container')!;
	// https://apache.github.io/echarts-handbook/en/best-practices/canvas-vs-svg
	chartInstance = echarts.init(chartDom, undefined, { renderer: 'svg' });
	let ps = get(previewState);
	let dims = [ps.xChosenSize, ps.yChosenSize];
	let mesh = get(meshState);
	let xData = Array.from({ length: dims[0] }, (_, i) => i * mesh.dx * 1e9 * mesh.Nx/ps.xChosenSize);
	let yData = Array.from({ length: dims[1] }, (_, i) => i * mesh.dy * 1e9* mesh.Ny/ps.yChosenSize);

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

	let data = get(previewState).scalarField;

	// @ts-ignore
	chartInstance.setOption({
		tooltip: {
			position: 'top',
			formatter: function (params: any) {
				if (params.value === undefined) {
					return 'NaN';
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
						return 'NaN';
					}
					return ` ${parseFloat(params.value).toFixed(0)} nm`;
				},
				padding: [8, 5, 8, 5],
				borderColor: '#6e9bcb',
				borderWidth: 1
			}
		},
		xAxis: {
			type: 'category',
			data: xData,
			name: 'x (nm)',
			nameLocation: 'middle',
			nameGap: 25,
			nameTextStyle: {
				color: '#fff'
			},
			axisTick: {
				alignWithLabel: true,
				length: 6,
				lineStyle: {
					type: 'solid',
					color: '#fff'
				}
			},
			axisLabel: {
				show: true,
				formatter: function (value: string, _index: string) {
					return parseFloat(value).toFixed(0);
				},
				color: '#fff',
				showMinLabel: true
			}
		},
		yAxis: {
			type: 'category',
			data: yData,
			name: 'y (nm)',
			nameLocation: 'middle',
			nameGap: 45,
			nameTextStyle: {
				color: '#fff'
			},
			axisTick: {
				alignWithLabel: true,
				length: 6,
				lineStyle: {
					type: 'solid',
					color: '#fff'
				}
			},
			axisLabel: {
				show: true,
				formatter: function (value: string, _index: string) {
					return parseFloat(value).toFixed(0);
				},
				color: '#fff',
				showMinLabel: true
			}
		},
		visualMap: [
			{
				type: 'continuous',
				min: ps.min,
				max: ps.max,
				calculable: true,
				realtime: true,
				// precision: 2,
				formatter: function (value: any) {
					return parseFloat(value).toPrecision(2);
				},
				itemWidth: 9,
				itemHeight: 140,
				text: [ps.unit, ''],
				textStyle: {
					color: '#fff'
				},
				inRange: {
					color: getColorMap(ps.min, ps.max)
				},
				left: 'right'
			}
		],
		series: [
			{
				name: ps.quantity,
				type: 'heatmap',
				emphasis: {
					itemStyle: {
						borderColor: '#333',
						borderWidth: 1
					}
				},
				selectedMode: false,
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
			right: '10%'
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
						borderWidth: 2
					}
				},
				dataView: { show: false },
				restore: {
					show: false
				},
				saveAsImage: {
					type: 'png',
					name: 'preview'
				}
			}
		}
	});
}

export function disposePreview2D() {
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
		if (chartInstance === undefined || chartInstance.isDisposed()) {
			return;
		}
	});
}
