<script>
	import { Button, Dropdown, DropdownItem, DropdownDivider } from 'flowbite-svelte';
	import { ChevronDownOutline, ChevronRightOutline } from 'flowbite-svelte-icons';
	import { previewState as p } from '$api/incoming/preview';
	import quantities from './quantities';
</script>

<Button
	>Dropdown button<ChevronDownOutline class="ms-2 h-6 w-6 text-white dark:text-white" /></Button
>
<Dropdown>
	<DropdownItem>Dashboard</DropdownItem>
	<DropdownItem class="flex items-center justify-between">
		Dropdown<ChevronRightOutline class="text-primary-700 ms-2 h-6 w-6 dark:text-white" />
	</DropdownItem>
	<Dropdown placement="right-start">
		<DropdownItem>Overview</DropdownItem>
		<DropdownItem>My downloads</DropdownItem>
		<DropdownItem>Billing</DropdownItem>
	</Dropdown>
	<DropdownItem>Earnings</DropdownItem>
	<DropdownItem slot="footer">Sign out</DropdownItem>
</Dropdown>

<!-- <Select bind:value={$p.quantity} on:change={postQuantity}>
    {#each quantities as quantity}
        <option value={quantity}>{quantity}</option>
    {/each}
</Select> -->
