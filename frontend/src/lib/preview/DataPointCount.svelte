<script lang="ts">
	import { previewState as p } from '$api/incoming/preview';
	import { postMaxPoints } from '$api/outgoing/preview';
	import { Input } from 'flowbite-svelte';

	let value = '';
	function inputChanged() {
		if (value === '') {
			return;
		}
		postMaxPoints(value);
	}
</script>

<div class="m-3 grid grid-cols-2 gap-2">
	<Input
		class="w-24"
		type="number"
		bind:value
		on:change={inputChanged}
		placeholder=" {$p.dataPointsCount}"
	/>
	<div class="flex items-center justify-center">data points</div>
</div>
