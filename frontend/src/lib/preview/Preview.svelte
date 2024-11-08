<script lang="ts">
	import { meshState } from '$api/incoming/mesh';
	import { previewState as p } from '$api/incoming/preview';
	import { threeDPreview, resetCamera } from './preview3D';
	import { resizeECharts } from './preview2D';
	import { onMount } from 'svelte';
	import { Button } from 'flowbite-svelte';
	import QuantityDropdown from './inputs/QuantityDropdown.svelte';
	import Component from './inputs/Component.svelte';
	import Layer from './inputs/Layer.svelte';
	import XDataPoints from './inputs/XDataPoints.svelte';
	import YDataPoints from './inputs/YDataPoints.svelte';
	import ResetCamera from './inputs/ResetCamera.svelte';

	onMount(() => {
		resizeECharts();
	});
</script>

<section>
	<h2 class="mb-4 text-2xl font-semibold">Preview</h2>

	<div class="m-1 flex flex-wrap" id="parent-fields">
		<div class="basis-1/2">
			<QuantityDropdown />
		</div>
		<div class="basis-1/2">
			<Component />
		</div>

		<div class="basis-1/2">
			{#if $p.xPossibleSizes.length > 0}
				<XDataPoints />
			{/if}
		</div>
		<div class="basis-1/2">
			{#if $p.yPossibleSizes.length > 0}
				<YDataPoints />
			{/if}
		</div>
		<div class="basis-1/2">
			<ResetCamera />
		</div>

		<div class="basis-1/2">
			<Layer />
		</div>
	</div>
	<hr />
	<div class="relative h-[500px] w-full">
		{#if $p.scalarField == null && $p.vectorFieldPositions == null}
			<div class="absolute inset-0 flex items-center justify-center text-6xl text-gray-600">
				NO DATA
			</div>
		{/if}
		<div id="container"></div>
	</div>
	<hr />
</section>

<style>
	section {
		grid-area: display;
	}
	#parent-fields > div {
		@apply p-1;
	}
	#container {
		width: 100%;
		height: 500px;
	}
</style>
