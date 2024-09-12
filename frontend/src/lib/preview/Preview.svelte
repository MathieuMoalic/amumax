<script lang="ts">
	import quantities from './quantities';
	import { meshState } from '$api/incoming/mesh';
	import { previewState as p } from '$api/incoming/preview';
	import { threeDPreview, resetCamera } from './preview3D';
	import { resizeECharts } from './preview2D';
	import { postComponent, postLayer, postMaxPoints, postQuantity } from '$api/outgoing/preview';
	import { onMount } from 'svelte';
	import NumberInputField from '$lib/NumberInputField.svelte';
	import { Button, Select } from 'flowbite-svelte';
	import QuantityDropdown from './QuantityDropdown.svelte';

	onMount(() => {
		resizeECharts();
	});
</script>

<section>
	<h2>Preview</h2>
	<div class="grid grid-cols-2 gap-6">
		<QuantityDropdown />
		<div>
			Quantity:
			<Select bind:value={$p.quantity} on:change={postQuantity}>
				{#each quantities as quantity}
					<option value={quantity}>{quantity}</option>
				{/each}
			</Select>
		</div>
		<div>
			{#if $p.nComp == 3}
				{#each ['3D', 'x', 'y', 'z'] as choice}
					<label>
						<input
							type="radio"
							bind:group={$p.component}
							value={choice}
							on:change={postComponent}
						/>
						{choice}
					</label>
				{/each}
			{/if}
		</div>
		{#if $p.scalarField == null && $p.vectorFieldPositions == null}
			<p>No data to display for {$p.quantity}</p>
		{:else}
			{#if $meshState.Nz > 1}
				<div>
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
				</div>
			{/if}
			<div>
				Max Points: <NumberInputField func={postMaxPoints} placeholder={$p.maxPoints} />
				({#if $p.scalarField !== null}
					{$p.scalarField?.length}
				{:else if $p.vectorFieldPositions !== null}
					{$p.vectorFieldPositions?.length}
				{/if}
				points)
			</div>
			{#if $threeDPreview !== null}
				<div>
					<Button on:click={resetCamera} outline>Reset Camera</Button>
				</div>
			{/if}
		{/if}
	</div>

	<div id="container"></div>
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
