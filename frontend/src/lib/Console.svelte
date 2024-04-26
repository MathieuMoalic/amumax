<script lang="ts">
	import { consoleState, postCommand } from './api';
	let command = '';
	function handleKeydown(event: KeyboardEvent) {
		if (event.key === 'Enter') {
			postCommand(command);
			event.preventDefault(); // Prevent the default action to avoid form submission or newline in input
		}
	}
</script>

<span>▾ console</span> <br />
<div>
	<textarea rows="16" cols="84">{$consoleState.hist}</textarea> <br />
	<input
		placeholder="type commands here, or up/down"
		size="86"
		bind:value={command}
		on:keydown={handleKeydown}
	/>
</div>
