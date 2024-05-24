<script lang="ts">
	import quantities from './quantities';
	import { meshState } from './api';
	import { refreshImage, imageUrl, renderLayer, imageComponent, imageQuantity } from './api';
</script>

<section>
	<h2>Preview</h2>
	<p>
		Quantity:
		<select bind:value={$imageQuantity} on:change={refreshImage}>
			{#each quantities as quantity}
				<option value={quantity}>{quantity}</option>
			{/each}
		</select>
		<select bind:value={$imageComponent}>
			{#each ['All', 'x', 'y', 'z'] as choice}
				<option value={choice}>{choice}</option>
			{/each}
		</select>
		Z-layer: 0
		<input type="range" min="0" max={$meshState.Nz - 1} bind:value={$renderLayer} />
		{$renderLayer}
	</p>
	<img src={$imageUrl} alt="Display" width="100%" />
</section>

<style>
	section {
		grid-area: display;
	}
</style>
