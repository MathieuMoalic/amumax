<script lang="ts">
	import { Range } from 'flowbite-svelte';

	export let possibleSizes: number[] = [];
	export let postChosenSize: (size: number) => void;
	let sliderValue = possibleSizes.length - 1; // Index in sizes array
	function onSliderChange(_: any) {
		postChosenSize(possibleSizes[sliderValue]);
	}
</script>

<Range bind:value={sliderValue} min="0" max={possibleSizes.length - 1} on:change={onSliderChange} />
{possibleSizes[sliderValue]}
