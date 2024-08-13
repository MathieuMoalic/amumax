<script lang="ts">
	import { tablePlotState } from '$api/incoming/table-plot';
	import { postAutoSaveInterval, postXColumn, postYColumn } from '$api/outgoing/table-plot';
	import { onMount } from 'svelte';
	import { resizeECharts } from './table-plot';
	onMount(() => {
		resizeECharts();
	});
</script>

<section>
	<h2>Table Plot</h2>
	{#if $tablePlotState.data.length === 0}
		<div class="msg">
			<p>No table data, use <code>TableSave()</code> to save data.</p>
		</div>
	{:else}
		<p>
			Auto Save Interval: <input
				type="number"
				bind:value={$tablePlotState.autoSaveInterval}
				on:change={postAutoSaveInterval}
			/> s
		</p>
		<b>
			x:
			<select bind:value={$tablePlotState.xColumn} on:change={postXColumn}>
				{#each $tablePlotState.columns as q}
					<option value={q}>{q}</option>
				{/each}
			</select>
			y:
			<select bind:value={$tablePlotState.yColumn} on:change={postYColumn}>
				{#each $tablePlotState.columns as q}
					<option value={q}>{q}</option>
				{/each}
			</select>
		</b>
		<div id="table-plot"></div>
	{/if}
</section>

<style>
	section {
		grid-area: tableplot;
	}
	#table-plot {
		width: 100%;
		height: 500px;
	}
	.msg {
		color: #888;
	}
</style>
