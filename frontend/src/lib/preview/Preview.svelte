<script lang="ts">
	import quantities from './quantities';
	import { meshState } from '$api/incoming/mesh';
	import { previewState as p } from '$api/incoming/preview';
	import { display, resetCamera } from './plot-vector-field';
	import {
		postComponent,
		postLayer,
		postMaxPoints,
		postQuantity,
		postRefresh
	} from '$api/outgoing/preview';
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
		<input type="range" min="1" max={$meshState.Nz} bind:value={$p.layer} on:change={postLayer} />
		{$p.layer}
		{$meshState.Nz}
		<button on:click={postRefresh}>Refresh</button>
	</p>
	Max Points:<input type="number" bind:value={$p.maxPoints} on:change={postMaxPoints} />
	<div id="container"></div>
	Parsing time: {$display?.parsingTime.toFixed(0)} ms
	<button on:click={resetCamera}>Reset Camera</button>
</section>

<style>
	section {
		grid-area: display;
	}
	#container {
		width: 100%;
		height: 500px;
		border: 1px solid black;
	}
</style>
