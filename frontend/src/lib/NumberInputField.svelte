<script lang="ts">
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

<input type="number" bind:value on:change={inputChanged} placeholder=" {placeholder}" />

<style>
	input {
		width: 70px;
	}
</style>
