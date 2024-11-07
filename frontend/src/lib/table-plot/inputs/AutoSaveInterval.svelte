<script lang="ts">
	import { tablePlotState } from '$api/incoming/table-plot';
	import { postAutoSaveInterval } from '$api/outgoing/table-plot';
	import { ButtonGroup, Input, InputAddon } from 'flowbite-svelte';
	let value = '';
	function inputChanged() {
		if (value === '') {
			return;
		}
		postAutoSaveInterval(value);
		value = '';
	}
</script>

<ButtonGroup class="flex h-11">
	<InputAddon class="w-fit whitespace-nowrap !bg-transparent">AutoSave Interval</InputAddon>
	<Input
		class="w-full truncate"
		bind:value
		on:change={inputChanged}
		placeholder=" {$tablePlotState.autoSaveInterval}"
	/>
	<InputAddon class="w-16 !bg-transparent">s</InputAddon>
</ButtonGroup>
