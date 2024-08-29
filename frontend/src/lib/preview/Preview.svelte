<script lang="ts">
	import quantities from './quantities';
	import { meshState } from '$api/incoming/mesh';
	import { previewState as p } from '$api/incoming/preview';
	import { threeDPreview, resetCamera } from './preview3D';
	import { resizeECharts } from './preview2D';
	import { postComponent, postLayer, postMaxPoints, postQuantity } from '$api/outgoing/preview';
	import { onMount } from 'svelte';
	import NumberInputField from '$lib/NumberInputField.svelte';
	onMount(() => {
		resizeECharts();
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
		Max Points: <NumberInputField func={postMaxPoints} placeholder={$p.maxPoints} />
		({#if $p.scalarField !== null}
			{$p.scalarField?.length}
		{:else if $p.vectorFieldPositions !== null}
			{$p.vectorFieldPositions?.length}
		{/if}
		points)
	</p>
	<p>
		{#if $threeDPreview?.parsingTime !== null}
			Parsing time: {$threeDPreview?.parsingTime.toFixed(0)} ms
		{/if}
	</p>
	{#if $threeDPreview !== null}
		<p>
			<button on:click={resetCamera}>Reset Camera</button>
		</p>
	{/if}
	<div id="container"></div>
</section>
/mnt/storage_2/scratch/pl0095-01/zelent/mannga/FMR/coupled_resonator/v11/Material_2/Tx_6e-06/xsize_600/ysize_350/sq_parm_0/rotation_45/anetnna_0/B0_0.03.mx3

<style>
	section {
		grid-area: display;
	}
	#container {
		width: 100%;
		height: 500px;
	}
</style>
