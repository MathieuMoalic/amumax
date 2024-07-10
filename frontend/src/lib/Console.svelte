<script lang="ts">
	import { consoleState } from '$api/incoming/console';
	import { postCommand } from '$api/outgoing/console';

	import Prism from 'prismjs';
	import 'prismjs/components/prism-go'; // Ensure the import path is correct

	let command = '';
	function handleKeydown(event: KeyboardEvent) {
		if (event.key === 'Enter') {
			postCommand(command);
			event.preventDefault(); // Prevent the default action to avoid form submission or newline in input
		}
	}
</script>

<section>
	<h2>Console</h2>
	<div class="parent">
		<div class="code">
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
	.parent {
		padding-right: 10px;
	}
	.code {
		white-space: pre-wrap;
		overflow-y: scroll;
		height: 20rem;
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
