<script lang="ts">
	import quantities from './quantities';
	import { meshState } from '$api/incoming/mesh';
	import { preview } from '$api/outgoing/preview';
	import { display } from './plot-vector-field';
	import { postPreview } from '$api/outgoing/preview';

	let maxPointsField = '1000';

	function onMaxPointsFieldChange() {
		$preview.maxPoints = parseInt(maxPointsField);
		postPreview();
	}
</script>

<section>
	<h2>Preview</h2>
	<p>
		Quantity:
		<select bind:value={$preview.quantity} on:change={postPreview}>
			{#each quantities as quantity}
				<option value={quantity}>{quantity}</option>
			{/each}
		</select>
		<select bind:value={$preview.component} on:change={postPreview}>
			{#each ['All', 'x', 'y', 'z'] as choice}
				<option value={choice}>{choice}</option>
			{/each}
		</select>
		Z-layer: 1
		<input type="range" min="0" max={$meshState.Nz - 1} bind:value={$preview.layer} />
		{$preview.layer}
		<button on:click={postPreview}>Refresh</button>
	</p>
	Max Points:<input bind:value={maxPointsField} on:change={onMaxPointsFieldChange} />
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
