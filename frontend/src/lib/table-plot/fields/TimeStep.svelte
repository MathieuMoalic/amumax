<script lang="ts">
	import { tablePlotState } from '$api/incoming/table-plot';
	import { postAutoSaveInterval } from '$api/outgoing/table-plot';
	import { Input } from 'flowbite-svelte';
	let value = '';
	function inputChanged() {
		if (value === '') {
			return;
		}
		postAutoSaveInterval(value);
	}
</script>

<div class="m-3 flex flex-col gap-2">
	<div class="flex items-center justify-center">Time step (s)</div>
	<Input
		class="w-24"
		type="number"
		bind:value
		on:change={inputChanged}
		placeholder=" {$tablePlotState.autoSaveInterval}"
	/>
</div>
