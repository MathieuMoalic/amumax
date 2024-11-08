<script lang="ts">
	import { previewState as p } from '$api/incoming/preview';
	import { postComponent } from '$api/outgoing/preview';
	import { Radio } from 'flowbite-svelte';
	let isDisabled: boolean;

	$: isDisabled = $p.nComp != 3;
</script>

<ul class="divide-x divide-gray-600 rounded-lg border border-gray-600 bg-gray-800 sm:flex">
	{#each ['3D', 'x', 'y', 'z'] as choice}
		<li class="w-full">
			<Radio
				name="hor-list"
				class="p-3 {isDisabled ? 'cursor-not-allowed opacity-50' : ''}"
				on:change={() => {
					postComponent(choice);
				}}
				bind:group={$p.component}
				value={choice}
				disabled={isDisabled}
			>
				{choice}
			</Radio>
		</li>
	{/each}
</ul>
