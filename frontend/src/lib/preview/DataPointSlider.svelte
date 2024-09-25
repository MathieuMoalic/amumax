<script lang="ts">
	import { Range } from 'flowbite-svelte';
	import { post } from '$api/post';
	import { previewState as p } from '$api/incoming/preview';
	import { meshState as m } from '$api/incoming/mesh';

	function inputChanged() {
		updateStep(step);
		post('preview/maxpoints', { maxPoints: pointsCount });
	}

	function largestPowerOfTwoExponent(maxPoint: number): number {
		let i = 0;
		while (1 << i <= maxPoint) {
			i++;
		}
		return i - 1; // Return the exponent of the largest power of two
	}

	function getStep(): number {
		let stepOut: number;
		const value = localStorage.getItem('previewStep');
		if (value !== null) {
			stepOut = parseInt(value);
		} else {
			stepOut = largestPowerOfTwoExponent($p.maxPoints);
			console.debug('step', stepOut);
		}
		return stepOut;
	}

	function updateStep(step: number) {
		localStorage.setItem('previewStep', step.toString());
	}

	function getSliderSteps(): number[] {
		let steps: number[] = [];
		let maxSteps = largestPowerOfTwoExponent(maxPoints);
		for (let i = 3; i <= maxSteps + 1; i += 1) {
			steps.push(i);
		}
		return steps;
	}

	let maxPoints: number = $m.Nx * $m.Ny;
	let step: number = getStep();
	let steps: number[] = getSliderSteps();

	$: pointsCount = step === steps[steps.length - 1] ? maxPoints : 1 << step;
</script>

<div class="flex w-64 items-center rounded-md border border-gray-600 p-2 pr-0">
	<Range
		id="range-steps"
		min={steps[0]}
		max={steps[steps.length - 1]}
		bind:value={step}
		on:change={inputChanged}
		step={1}
		class="w-32"
	/>
	<div class="w-20 overflow-clip">{pointsCount}/{maxPoints}</div>
</div>
