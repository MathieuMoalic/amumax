<script lang="ts">
	import quantities from './quantities';
	import { meshState } from '$api/incoming/mesh';
	import { previewState as p } from '$api/incoming/preview';
	import { threeDPreview, disposePreview3D, resetCamera, preview3D } from './preview3D';
	import { disposePreview2D, init, resizeECharts } from './preview2D';
	import { postComponent, postLayer, postMaxPoints, postQuantity } from '$api/outgoing/preview';
	import { onMount } from 'svelte';
	onMount(() => {
		resizeECharts();
	});
	let maxPointsInput = '';
	function maxPointsInputChanged() {
		if (maxPointsInput === '') {
			return;
		}
		postMaxPoints(parseInt(maxPointsInput));
	}
</script>

<section>
	<h2>Preview</h2>
	<p>
		Quantity:
		<select bind:value={$p.quantity} on:change={postQuantity}>
			{#each quantities as quantity}
				<option value={quantity}>{quantity}</option>
			{/each}
		</select>
		<select bind:value={$p.component} on:change={postComponent}>
			{#each ['All', 'x', 'y', 'z'] as choice}
				<option value={choice}>{choice}</option>
			{/each}
		</select>
	</p>
	<p>
		Z-layer: 0
		<input
			type="range"
			min="0"
			max={$meshState.Nz - 1}
			bind:value={$p.layer}
			on:change={postLayer}
		/>
		{$p.layer}
		{$meshState.Nz}
	</p>
	<p>
		Max Points:<input
			bind:value={maxPointsInput}
			on:change={maxPointsInputChanged}
			placeholder=" {$p.maxPoints}"
		/>
		Number of points in the preview:
		{#if $p.scalarField !== null}
			{$p.scalarField?.length}
		{:else if $p.vectorFieldPositions !== null}
			{$p.vectorFieldPositions?.length}
		{/if}
	</p>
	<div id="container"></div>
	<p>
		Parsing time: {$threeDPreview?.parsingTime.toFixed(0)} ms
	</p>
	<p>
		<button on:click={resetCamera}>Reset Camera</button>
		<button on:click={disposePreview2D}>disposePreview2D</button>
		<button on:click={disposePreview3D}>disposePreview3D</button>
		<button on:click={init}>init 2D</button>
		<button on:click={preview3D}>init 3D</button>
	</p>
</section>

<style>
	section {
		grid-area: display;
	}
	#container {
		width: 100%;
		height: 500px;
	}
</style>
