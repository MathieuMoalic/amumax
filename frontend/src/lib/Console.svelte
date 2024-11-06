<script lang="ts">
	import { consoleState } from '$api/incoming/console';
	import { postCommand } from '$api/outgoing/console';

	import Prism from 'prismjs';
	import 'prismjs/components/prism-go'; // Ensure the import path is correct
	import { onMount, tick } from 'svelte';

	let command = '';
	let codeDiv: HTMLDivElement | null = null;

	// Handle Enter key to submit commands
	function handleKeydown(event: KeyboardEvent) {
		if (event.key === 'Enter') {
			postCommand(command);
			event.preventDefault(); // Prevent the default action to avoid form submission or newline in input
		}
	}
	async function scrollDown() {
		if (codeDiv) {
			// Wait for the DOM to update (e.g. after a new line in hist is added)
			await tick();
			const isAtBottom = codeDiv.scrollTop + codeDiv.clientHeight >= codeDiv.scrollHeight - 100;
			if (isAtBottom) {
				codeDiv.scrollBy(0, codeDiv.scrollHeight);
			}
		}
	}

	$: {
		$consoleState.hist; // Trigger the reactive statement when hist changes
		scrollDown();
	}
	onMount(() => {
		scrollDown();
	});
</script>

<section>
	<h2 class="mb-4 text-2xl font-semibold">Console</h2>
	<div class="flex flex-col gap-2">
		<div class="code" bind:this={codeDiv}>
			{@html Prism.highlight($consoleState.hist, Prism.languages['go'], 'go')}
		</div>
		<br />
		<input
			placeholder="type commands here, or up/down"
			size="86"
			bind:value={command}
			on:keydown={handleKeydown}
		/>
	</div>
</section>

<style>
	section {
		grid-area: console;
	}
	.code {
		white-space: pre-wrap;
		overflow-y: scroll;
		height: 30rem;
		border: 1px solid #f8f8f2;
		width: 100%;
		padding: 5px;
	}
	input {
		padding: 5px;
		border: 1px solid #f8f8f2;
		width: 100%;
	}
</style>
