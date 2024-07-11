<script lang="ts">
	import Header from '$lib/Header.svelte';
	import Display from '$lib/preview/Preview.svelte';
	import TablePlot from '$lib/table-plot/TablePlot.svelte';
	import Solver from '$lib/Solver.svelte';
	import Console from '$lib/Console.svelte';
	import Mesh from '$lib/Mesh.svelte';
	import Geometry from '$lib/Geometry.svelte';
	import InitialMag from '$lib/InitialMag.svelte';
	import Parameters from '$lib/Parameters.svelte';
	import Footer from '$lib/Footer.svelte';
	import { onMount } from 'svelte';
	import { initializeWebSocket, ws } from '$api/websocket';

	onMount(() => {
		try {
			initializeWebSocket(ws);
		} catch (error) {
			console.error('Error connecting to websocket:', error);
		}
	});
</script>

<div class="grid-container">
	<Header />
	<Display />
	<TablePlot />
	<Solver />
	<Console />
	<Mesh />
	<Geometry />
	<InitialMag />
	<Parameters />
	<Footer />
</div>

<style>
	.grid-container {
		display: grid;
		grid-template-areas:
			'header header header header'
			'display display tableplot tableplot'
			'console console solver solver'
			'parameters parameters mesh mesh'
			'parameters parameters geometry geometry'
			'parameters parameters initialmag initialmag'
			'footer footer footer footer';
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
				'parameters'
				'footer';
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
