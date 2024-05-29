<script lang="ts">
	import quantities from './quantities';
	import { meshState } from './api';
	import { renderLayer, imageComponent, imageQuantity } from './api';
	import { onMount } from 'svelte';
	import { init, update, parsingTime } from './display';

	onMount(() => {});
</script>

<section>
	<h2>Preview</h2>
	<p>
		Quantity:
		<select bind:value={$imageQuantity} on:change={update}>
			{#each quantities as quantity}
				<option value={quantity}>{quantity}</option>
			{/each}
		</select>
		<select bind:value={$imageComponent}>
			{#each ['All', 'x', 'y', 'z'] as choice}
				<option value={choice}>{choice}</option>
			{/each}
		</select>
		Z-layer: 1
		<input type="range" min="0" max={$meshState.Nz - 1} bind:value={$renderLayer} />
		{$renderLayer}
		<button on:click={update}>Refresh</button>
	</p>
	<div id="container"></div>
	Parsing time: {$parsingTime} ms
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
