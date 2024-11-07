<script lang="ts">
	import { tablePlotState } from '$api/incoming/table-plot';
	import { postMaxPoints } from '$api/outgoing/table-plot';
	import { ButtonGroup, Input, InputAddon } from 'flowbite-svelte';
	let value = '';
	function inputChanged() {
		if (value === '') {
			return;
		}
		postMaxPoints(value);
	}
</script>

<ButtonGroup class="flex h-11">
	<InputAddon class="w-fit whitespace-nowrap !bg-transparent ">Max Points</InputAddon>
	<Input
		class="w-full truncate"
		bind:value
		on:change={inputChanged}
		placeholder=" {$tablePlotState.maxPoints}"
	/>
</ButtonGroup>
