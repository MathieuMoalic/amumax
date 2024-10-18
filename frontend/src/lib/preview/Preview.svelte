<script lang="ts">
	import { meshState } from '$api/incoming/mesh';
	import { previewState as p } from '$api/incoming/preview';
	import { threeDPreview, resetCamera } from './preview3D';
	import { resizeECharts } from './preview2D';
	import { onMount } from 'svelte';
	import { Button } from 'flowbite-svelte';
	import QuantityDropdown from './QuantityDropdown.svelte';
	import Component from './Component.svelte';
	import Layer from './Layer.svelte';
	import DataPointSlider from './DataPointSlider.svelte';

	onMount(() => {
		resizeECharts();
	});
</script>

<section>
	<h2 class="mb-4 text-2xl font-semibold">Preview</h2>

	<div class="m-3 grid grid-cols-10 gap-2">
		<div class="field col-span-2">
			<QuantityDropdown />
		</div>
		<div class="field col-span-4">
			<Component />
		</div>

		{#if $p.scalarField == null && $p.vectorFieldPositions == null}
			<div class="field col-span-4">
				No data to display for {$p.quantity}
			</div>
		{:else}
			<div class="field col-span-4">
				{#if $meshState.Nx > 1}
					<DataPointSlider />
				{/if}
			</div>
			{#if $threeDPreview !== null}
				<div class="field col-span-3">
					<Button on:click={resetCamera} outline>Reset Camera</Button>
				</div>
			{:else}
				<div class="field col-span-3"></div>
			{/if}

			{#if $meshState.Nz > 1}
				<div class="field col-span-7">
					<Layer />
				</div>
			{/if}
		{/if}
	</div>
	<hr />
	<div id="container"></div>
	<hr />
</section>

<style>
	section {
		grid-area: display;
	}
	#container {
		width: 100%;
		height: 500px;
	}
	.field {
		@apply flex items-center justify-center border-0 border-green-500;
	}
</style>
