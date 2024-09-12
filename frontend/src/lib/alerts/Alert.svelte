<script>
	import { Toast } from 'flowbite-svelte';
	import { alert } from './alert';
	import { ExclamationCircleSolid } from 'flowbite-svelte-icons';
</script>

{#if $alert !== ''}
	<div class="fixed left-3/4 top-10 z-50 gap-10">
		<Toast color="red">
			<svelte:fragment slot="icon">
				<ExclamationCircleSolid class="h-5 w-5" />
				<span class="sr-only">Warning icon</span>
			</svelte:fragment>
			{$alert}
		</Toast>
	</div>
{/if}
