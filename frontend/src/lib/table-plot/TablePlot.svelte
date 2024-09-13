<script lang="ts">
	import { tablePlotState } from '$api/incoming/table-plot';
	import { onMount } from 'svelte';
	import { resizeECharts } from './table-plot';
	import TimeStep from './fields/TimeStep.svelte';
	import Step from './fields/Step.svelte';
	import XColumn from './fields/XColumn.svelte';
	import YColumn from './fields/YColumn.svelte';
	import DataPointsCount from './fields/DataPointsCount.svelte';

	onMount(() => {
		resizeECharts();
	});
</script>

<section>
	<h2 class="mb-4 text-2xl font-semibold">Table Plot</h2>
	{#if $tablePlotState.data.length === 0}
		<div class="msg">
			<p>No table data, use <code>TableSave()</code> to save data.</p>
		</div>
	{:else}
		<div class="m-3 grid grid-cols-10 gap-2">
			<div class="field col-span-2">
				<TimeStep />
			</div>
			<div class="field col-span-2">
				<DataPointsCount />
			</div>
			<div class="field col-span-2">
				<Step />
			</div>
			<div class="field col-span-2">
				<XColumn />
			</div>
			<div class="field col-span-2">
				<YColumn />
			</div>
		</div>
		<hr />
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
	.field {
		@apply flex items-center justify-center border-0 border-green-500;
	}
</style>
