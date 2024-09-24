<script lang="ts">
	import Header from '$lib/Header.svelte';
	import Display from '$lib/preview/Preview.svelte';
	import TablePlot from '$lib/table-plot/TablePlot.svelte';
	import Solver from '$lib/Solver.svelte';
	import Console from '$lib/Console.svelte';
	import Mesh from '$lib/Mesh.svelte';
	import Parameters from '$lib/Parameters.svelte';
</script>

<Header />
<div class="grid-container">
	<Display />
	<TablePlot />
	<Solver />
	<Console />
	<Mesh />
	<Parameters />
</div>

<style>
	.grid-container {
		display: grid;
		grid-template-areas:
			'display tableplot'
			'console solver'
			'parameters mesh';
		margin: 0;
		grid-gap: 10px;
		padding: 10px;
		max-width: 100%;
		box-sizing: border-box;
	}
	@media (max-width: 1200px) {
		.grid-container {
			grid-template-areas:
				'header'
				'display'
				'tableplot'
				'console'
				'solver'
				'mesh'
				'geometry'
				'initialmag'
				'parameters';
		}
	}
	:global(section) {
		border: 2px solid #f8f8f2;
		padding: 10px;
		padding-top: rem;
		border-radius: 5px;
		background-color: #282a36;
		/* max-width: 30%; */
		overflow: hidden;
	}
	:global(body) {
		font-family: 'Arial', sans-serif;
	}
	:global(section > h2) {
		background-color: #6272a4;
		color: #f8f8f2;
		padding: 1px;
		padding-left: 1rem;
		border-radius: 3px 3px 0 0;
		margin-top: 0;
	}
</style>
