<script lang="ts">
	import { onMount } from 'svelte';
	import { tablePlotState, postTableColumns } from './api';
	import { Chart, registerables } from 'chart.js';
	import zoomPlugin from 'chartjs-plugin-zoom';

	let chart: Chart | null = null;

	Chart.register(...registerables, zoomPlugin);

	async function onNewColumns() {
		await postTableColumns();
		redraw();
	}

	onMount(() => {
		createChart();
	});

	function createChart() {
		chart = new Chart(document.getElementById('gd') as HTMLCanvasElement, {
			type: 'scatter',
			data: {
				datasets: [
					{
						data: $tablePlotState.data,
						backgroundColor: 'rgba(240, 192, 192, 1)',
						borderColor: 'rgba(75, 192, 192, 1)',
						borderWidth: 2,
						pointRadius: 0,
						showLine: true
					}
				]
			},
			options: {
				responsive: true,
				scales: {
					x: {
						beginAtZero: false,
						grid: {
							color: '#f8f8f2'
						},
						ticks: {
							color: '#f8f8f2'
						}
					},
					y: {
						beginAtZero: false,
						grid: {
							color: '#f8f8f2'
						},
						ticks: {
							color: '#f8f8f2'
						}
					}
				},
				plugins: {
					legend: {
						display: false
					},
					zoom: {
						pan: {
							enabled: true
						},
						zoom: {
							wheel: {
								enabled: true,
								speed: 0.1
							}
						}
					}
				}
			}
		});
	}

	function redraw() {
		if (chart) {
			chart.data.datasets[0].data = $tablePlotState.data;
			chart.update();
		}
	}
	function resetZoom() {
		if (chart) {
			chart.resetZoom('active');
		}
	}
</script>

<section>
	<h2>Table Plot</h2>
	<p>
		Auto Save Interval: <input bind:value={$tablePlotState.autoSaveInterval} /> s
	</p>
	<b>
		x:
		<select bind:value={$tablePlotState.xColumn} on:change={onNewColumns}>
			{#each $tablePlotState.columns as q}
				<option value={q}>{q}</option>
			{/each}
		</select>
		y:
		<select bind:value={$tablePlotState.yColumn} on:change={onNewColumns}>
			{#each $tablePlotState.columns as q}
				<option value={q}>{q}</option>
			{/each}
		</select>
	</b>
	<button on:click={resetZoom}>Reset Zoom</button>
	<button on:click={redraw}>Redraw</button>
	<div class="plot">
		<canvas id="gd"></canvas>
	</div>
</section>

<style>
	section {
		grid-area: tableplot;
	}
	div.plot {
		padding: 20px;
	}
</style>
