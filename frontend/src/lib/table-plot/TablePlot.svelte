<script lang="ts">
	import { onMount } from 'svelte';
	import { tablePlotState } from '$api/incoming/table-plot';
	import { postTableColumns } from '$api/outgoing/table-plot';
	import { Chart, registerables } from 'chart.js';
	import zoomPlugin from 'chartjs-plugin-zoom';
	import { plotChart, createChart, resetZoom, chart } from './tablePlot';

	Chart.register(...registerables, zoomPlugin);

	async function onNewColumns() {
		await postTableColumns();
		plotChart();
	}
	onMount(() => {
		createChart();
	});
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
	<button on:click={plotChart}>Redraw</button>
	<div class="plot-container">
		<canvas id="plot"></canvas>
	</div>
</section>

<style>
	section {
		grid-area: tableplot;
	}
	.plot-container {
		position: relative;
		margin: auto;
		height: 70%;
		padding: 20px;
	}
</style>
