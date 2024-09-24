<script lang="ts">
	import { Range } from 'flowbite-svelte';
	import { post } from '$api/post';
	import { meshState as m } from '$api/incoming/mesh';

	function largestPowerOfTwoExponent(maxPoint: number): number {
		let i = 0;
		while (1 << i <= maxPoint) {
			i++;
		}
		return i - 1; // Return the exponent of the largest power of two
	}

	function inputChanged() {
		post('preview/maxpoints', { maxPoints: value });
	}

	let maxPoints: number = 0;
	let steps: number[] = [];

	let exp = 8; // Initial value
	$: value = exp === steps[steps.length - 1] ? maxPoints : 1 << exp;

	maxPoints = $m.Nx * $m.Ny;
	let maxSteps = largestPowerOfTwoExponent(maxPoints);
	for (let step = 3; step <= maxSteps + 1; step += 1) {
		steps.push(step);
	}
</script>

<div class="flex w-64 items-center rounded-md border border-gray-600 p-2 pr-0">
	<Range
		id="range-steps"
		min={steps[0]}
		max={steps[steps.length - 1]}
		bind:value={exp}
		on:change={inputChanged}
		step={1}
		class="w-32"
	/>
	<div class="w-20 overflow-clip">{value}/{maxPoints}</div>
</div>
