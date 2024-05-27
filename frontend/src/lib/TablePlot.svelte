<script lang="ts">
	import { onMount } from 'svelte';
	import { tablePlotState, postTableColumns } from './api';

	onMount(async () => {
		const Plotly = await import('plotly.js-dist-min');
		Plotly.newPlot('gd', {
			data: [{ x: $tablePlotState.xdata, y: $tablePlotState.ydata }],
			layout: {
				width: 600,
				height: 400,
				title: 'Simple Plot',
				paper_bgcolor: 'lightgray',
				plot_bgcolor: 'white'
			}
		});
	});
	async function redraw() {
		const Plotly = await import('plotly.js-dist-min');
		Plotly.react('gd', [{ x: $tablePlotState.xdata, y: $tablePlotState.ydata }]);
		// Plotly.extendTraces('gd', { x: [$tablePlotState.xdata], y: [$tablePlotState.ydata] }, [0]);
	}
</script>

<section>
	<h2>Table Plot</h2>
	<p>
		Auto Save Interval: <input bind:value={$tablePlotState.autoSaveInterval} /> s
	</p>
	<b>
		x:
		<select bind:value={$tablePlotState.xColumn} on:change={postTableColumns}>
			{#each $tablePlotState.columns as q}
				<option value={q}>{q}</option>
			{/each}
		</select>
		y:
		<select bind:value={$tablePlotState.yColumn} on:change={postTableColumns}>
			{#each $tablePlotState.columns as q}
				<option value={q}>{q}</option>
			{/each}
		</select>
	</b>
	<div id="gd"></div>
	<button on:click={redraw}>Redraw</button>
</section>

<style>
	section {
		grid-area: tableplot;
	}
</style>
