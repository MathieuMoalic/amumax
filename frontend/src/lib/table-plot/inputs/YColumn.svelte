<script lang="ts">
	import { tablePlotState } from '$api/incoming/table-plot';
	import { postYColumn } from '$api/outgoing/table-plot';
	import { Button, Dropdown, DropdownItem } from 'flowbite-svelte';
	import { ChevronDownOutline } from 'flowbite-svelte-icons';
	let dropdownOpen = false;
</script>

<Button outline color="primary" class="h-11 w-full justify-between">
	<span>Y Axis: </span>
	<span class="truncate text-white">{$tablePlotState.yColumn}</span>
	<ChevronDownOutline class="h-5 w-5 text-gray-500" />
</Button>
<Dropdown bind:open={dropdownOpen} class="w-3/4">
	{#each $tablePlotState.columns as q}
		<DropdownItem
			on:click={(_) => {
				postYColumn(q);
				dropdownOpen = false;
			}}
		>
			{q}
		</DropdownItem>
	{/each}
</Dropdown>
