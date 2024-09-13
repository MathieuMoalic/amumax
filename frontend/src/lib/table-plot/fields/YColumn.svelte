<script lang="ts">
	import { tablePlotState } from '$api/incoming/table-plot';
	import { postYColumn } from '$api/outgoing/table-plot';
	import { Button, Dropdown, DropdownItem } from 'flowbite-svelte';
	import { ChevronDownOutline } from 'flowbite-svelte-icons';
	let dropdownOpen = false;
</script>

<div class="m-3 flex flex-col gap-2">
	<div class="flex items-center justify-center">Y axis</div>
	<Button outline class="flex items-center overflow-hidden text-ellipsis whitespace-nowrap">
		<span class="truncate">
			{$tablePlotState.yColumn}
		</span>
		<ChevronDownOutline class="ms-2 h-6 w-6 text-white dark:text-white" />
	</Button>
	<Dropdown bind:open={dropdownOpen}>
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
</div>
