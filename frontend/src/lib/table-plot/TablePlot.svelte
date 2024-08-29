<script lang="ts">
	import { tablePlotState } from '$api/incoming/table-plot';
	import {
		postAutoSaveInterval,
		postMaxPoints,
		postStep,
		postXColumn,
		postYColumn
	} from '$api/outgoing/table-plot';
	import { onMount } from 'svelte';
	import { resizeECharts } from './table-plot';
	import NumberInputField from '$lib/NumberInputField.svelte';

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
		<div>
			Auto Save Interval: <NumberInputField
				func={postAutoSaveInterval}
				placeholder={$tablePlotState.autoSaveInterval}
			/> s
		</div>
		<div>
			Max Points: <NumberInputField func={postMaxPoints} placeholder={$tablePlotState.maxPoints} />
		</div>
		<div>
			Step: <NumberInputField func={postStep} placeholder={$tablePlotState.step} />
		</div>
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
		({$tablePlotState.data.length} points)
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
