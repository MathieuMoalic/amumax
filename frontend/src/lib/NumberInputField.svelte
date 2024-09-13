<script lang="ts">
	import { Input } from 'flowbite-svelte';

	let value = '';
	export let func = (_: any) => {};
	export let placeholder: any;
	function inputChanged() {
		if (value === '') {
			return;
		}
		func(value);
	}
</script>

<Input
	class="w-32"
	type="number"
	bind:value
	on:change={inputChanged}
	placeholder=" {placeholder}"
/>
