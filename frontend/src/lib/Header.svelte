<script lang="ts">
	import { headerState } from '$api/incoming/header';
	import { connected } from '$api/websocket';
	import { PlayOutline, PauseOutline, CloseCircleOutline } from 'flowbite-svelte-icons';
</script>

<section class="sticky top-0 z-10 m-2.5 mb-0 p-2">
	<div class="flex w-screen items-center gap-2">
		<div class="text-3xl">
			{#if $connected}
				{#if $headerState.status === 'running'}
					<PlayOutline color="green" size="xl" />
				{:else if $headerState.status === 'paused'}
					<PauseOutline color="orange" size="xl" />
				{:else}
					<CloseCircleOutline color="red" size="xl" />
				{/if}
			{:else}
				<CloseCircleOutline color="red" size="xl" />
			{/if}
		</div>

		<div class="min-w-0 flex-grow truncate text-3xl text-green-500">
			{$headerState.path}
		</div>

		<div class="w-40 whitespace-nowrap text-xl text-gray-500">
			v.{$headerState.version}
		</div>
	</div>
</section>
