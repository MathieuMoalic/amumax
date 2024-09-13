<script lang="ts">
	import { previewState as p } from '$api/incoming/preview';
	import { postComponent } from '$api/outgoing/preview';
	import { Radio } from 'flowbite-svelte';
</script>

{#if $p.nComp == 3}
	<ul
		class="w-full items-center divide-x divide-gray-200 rounded-lg border border-gray-200 disabled:bg-slate-50 sm:flex rtl:divide-x-reverse dark:divide-gray-600 dark:border-gray-600 dark:bg-gray-800"
	>
		{#each ['3D', 'x', 'y', 'z'] as choice}
			<li class="w-full">
				<Radio
					name="hor-list"
					class="p-3"
					on:change={() => {
						postComponent(choice);
					}}
					bind:group={$p.component}
					value={choice}
				>
					{choice}
				</Radio>
			</li>
		{/each}
	</ul>
{:else}
	<div></div>
{/if}
