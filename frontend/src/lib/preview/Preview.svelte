<script lang="ts">
	import quantities from './quantities';
	import { meshState } from '$api/incoming/mesh';
	import { previewState as p } from '$api/incoming/preview';
	import { display } from './plot-vector-field';
	import {
		postComponent,
		postLayer,
		postMaxPoints,
		postQuantity,
		postRefresh
	} from '$api/outgoing/preview';
	import { onMount } from 'svelte';

	let maxPointsField = 'loading...';
	onMount(() => {
		maxPointsField = $p.maxPoints.toString();
	});
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
		Z-layer: 1
		<input
			type="range"
			min="0"
			max={$meshState.Nz - 1}
			bind:value={$p.layer}
			on:change={postLayer}
		/>
		{$p.layer}
		<button on:click={postRefresh}>Refresh</button>
	</p>
	Max Points:<input bind:value={maxPointsField} on:change={postMaxPoints} />
	<div id="container"></div>
	Parsing time: {$display?.parsingTime.toFixed(0)} ms
</section>

<style>
	section {
		grid-area: display;
	}
	#container {
		width: 500px;
		height: 500px;
		border: 1px solid black;
	}
</style>
