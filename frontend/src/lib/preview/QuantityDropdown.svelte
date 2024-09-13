<script>
	import { Button, Dropdown, DropdownItem, DropdownDivider } from 'flowbite-svelte';
	import { ChevronDownOutline, ChevronRightOutline } from 'flowbite-svelte-icons';

	import { previewState as p } from '$api/incoming/preview';
	import { quantities } from './quantities';
	import { postQuantity } from '$api/outgoing/preview';
	let dropdownOpen = false;
</script>

<Button outline class="flex items-center overflow-hidden text-ellipsis whitespace-nowrap">
	<span class="truncate">
		{$p.quantity}
	</span>
	<ChevronDownOutline class="ms-2 h-6 w-6 text-white dark:text-white" />
</Button>

<Dropdown bind:open={dropdownOpen} placement="right-start">
	{#each quantities['Common'] as quantity}
		<DropdownItem
			on:click={(_) => {
				postQuantity(quantity);
				dropdownOpen = false;
			}}
		>
			{quantity}
		</DropdownItem>
	{/each}
	<DropdownDivider />
	{#each Object.entries(quantities) as [category, items]}
		{#if category != 'Common'}
			<DropdownItem class="flex items-center justify-between">
				{category}<ChevronRightOutline class="text-primary-700 ms-2 h-6 w-6 dark:text-white" />
			</DropdownItem>
			<Dropdown placement="right-start" trigger="hover">
				{#each items as quantity}
					<DropdownItem
						on:click={(_) => {
							postQuantity(quantity);
							dropdownOpen = false;
						}}
					>
						{quantity}
					</DropdownItem>
				{/each}
			</Dropdown>
		{/if}
	{/each}
</Dropdown>
