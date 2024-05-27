<script lang="ts">
	import { onMount } from 'svelte';
	import { tablePlotState, postTableColumns } from './api';
	import { Chart, registerables } from 'chart.js';
	import zoomPlugin from 'chartjs-plugin-zoom';
	import { redraw, createChart, resetZoom, chart } from './tablePlot';

	Chart.register(...registerables, zoomPlugin);

	async function onNewColumns() {
		await postTableColumns();
		redraw();
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
