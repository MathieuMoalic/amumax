<script lang="ts">
	import { parametersState } from '$api/incoming/parameters';
	let selectedRegion = 0;
</script>

<section>
	<h2>Parameters</h2>
	<div>
		Region:
		<select bind:value={selectedRegion} on:change={() => {}}>
			{#each $parametersState.regions as region}
				<option value={region}>{region}</option>
			{/each}
		</select>
		<div class="grid-container">
			{#each $parametersState.fields as field}
				<div class="header">{field.name}:</div>
				<input placeholder=" {field.value}" />
				<div class="description">{field.description}</div>
			{/each}
		</div>
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
