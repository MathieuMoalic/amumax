<script lang="ts">
	import { parametersState as p } from '$api/incoming/parameters';
	import { postSelectedRegion } from '$api/outgoing/parameters';
	import { Toggle, Button, Dropdown, DropdownItem } from 'flowbite-svelte';
	import { ChevronDownOutline } from 'flowbite-svelte-icons';
	let dropdownOpen = false;
	let showZeroValues = false;
</script>

<section>
	<h2 class="mb-4 text-2xl font-semibold">Parameters</h2>
	<div class="m-3 flex flex-col gap-2">
		<div class="m-3 grid grid-cols-2 gap-2">
			<div class="flex items-center justify-center">
				<Button outline>
					Region: {$p.selectedRegion}
					<ChevronDownOutline class="ms-2 h-6 w-6 text-white dark:text-white" />
				</Button>
				<Dropdown bind:open={dropdownOpen} class="h-48 w-40 overflow-y-auto py-1">
					{#each $p.regions as region}
						<DropdownItem
							on:click={(_) => {
								postSelectedRegion(region);
								dropdownOpen = false;
							}}
						>
							{region}
						</DropdownItem>
					{/each}
				</Dropdown>
			</div>
			<div class="flex items-center justify-center">
				<Toggle bind:checked={showZeroValues}>Show Unchanged Parameters</Toggle>
			</div>
		</div>
	</div>
	<div class="grid-container">
		{#each $p.fields as field}
			{#if field.changed || showZeroValues}
				<div class="header">{field.name}:</div>
				<input placeholder=" {field.value}" />
				<div class="description">{field.description}</div>
			{/if}
		{/each}
	</div>
</section>

<style>
	section {
		grid-area: parameters;
	}
	.grid-container {
		display: grid;
		grid-template-columns: auto auto auto;
		gap: 10px;
		align-items: center;
	}
	input {
		width: 100%;
		box-sizing: border-box;
	}
	.header {
		text-align: right;
		padding-right: 10px;
	}
	.description {
		font-style: italic;
		color: bisque;
		overflow: hidden;
	}
</style>
