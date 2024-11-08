<script lang="ts">
	export let label: string = 'Label';
	export let values: number[] = [];
	export let value: number;
	export let onChangeFunction: (value: number) => void;
	export let isDisabled: boolean = false;

	// Derived variables
	$: sliderMin = 0;
	$: sliderMax = values.length - 1;
	let index = values.indexOf(value);
	// Ensure value is one of the values in the array
	$: {
		const idx = values.indexOf(value);
		if (idx !== -1) {
			index = idx;
		} else {
			index = 0;
			value = values[0];
		}
	}

	// Update value when index changes
	$: if (values[index] !== value) {
		value = values[index];
	}
</script>

<div
	class="relative h-11 overflow-hidden rounded-md border border-gray-600 bg-gray-800 {isDisabled
		? 'cursor-not-allowed opacity-50'
		: ''}"
	style="--index:{index}; --min:{sliderMin}; --max:{sliderMax};"
>
	<!-- Filled portion of the slider -->
	<div class="filled-portion absolute left-0 top-0 h-full rounded-md bg-gray-700"></div>

	<!-- Invisible range input to capture user interactions -->
	<input
		type="range"
		min={sliderMin}
		max={sliderMax}
		step="1"
		bind:value={index}
		on:input={() => onChangeFunction(value)}
		class="absolute left-0 top-0 h-full w-full cursor-pointer opacity-0"
		disabled={isDisabled}
	/>

	<!-- Min and Max labels -->
	<div class="pointer-events-none absolute inset-0 flex items-center justify-between px-4">
		<span class=" text-gray-400">{values[0]}</span>
		<span class=" text-gray-400">{values[values.length - 1]}</span>
	</div>

	<!-- Current value label -->
	<div class="pointer-events-none absolute inset-0 flex items-center justify-center">
		<span class=" text-gray-200">{label} = {value}</span>
	</div>
</div>

<style>
	.filled-portion {
		width: calc((var(--index) - var(--min)) / (var(--max) - var(--min)) * 100%);
		transition: width 0.2s ease;
	}
</style>
