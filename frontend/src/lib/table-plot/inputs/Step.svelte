<script lang="ts">
	import { tablePlotState } from '$api/incoming/table-plot';
	import { postStep } from '$api/outgoing/table-plot';
	import { ButtonGroup, Input, InputAddon } from 'flowbite-svelte';

	let value = '';
	function inputChanged() {
		if (value === '') {
			return;
		}
		postStep(value);
	}
</script>

<ButtonGroup class="flex h-11">
	<InputAddon class="w-16 whitespace-nowrap !bg-transparent ">Step</InputAddon>
	<Input
		class="w-full truncate"
		bind:value
		on:change={inputChanged}
		placeholder=" {$tablePlotState.step}"
	/>
</ButtonGroup>
