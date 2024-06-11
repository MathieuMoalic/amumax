<script lang="ts">
	import quantities from './quantities';
	import { meshState } from '$api/incoming/mesh';
	import { preview } from '$api/outgoing/preview';
	import { onMount } from 'svelte';
	import { update, display } from '$api/incoming/preview';

	onMount(() => {});
</script>

<section>
	<h2>Preview</h2>
	<p>
		Quantity:
		<select bind:value={$preview.Quantity} on:change={update}>
			{#each quantities as quantity}
				<option value={quantity}>{quantity}</option>
			{/each}
		</select>
		<select bind:value={$preview.Component}>
			{#each ['All', 'x', 'y', 'z'] as choice}
				<option value={choice}>{choice}</option>
			{/each}
		</select>
		Z-layer: 1
		<input type="range" min="0" max={$meshState.Nz - 1} bind:value={$preview.Layer} />
		{$preview.Layer}
		<button on:click={update}>Refresh</button>
	</p>
	<div id="container"></div>
	Parsing time: {$display?.parsingTime} ms
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
