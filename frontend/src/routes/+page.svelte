<script lang="ts">
	import Header from '$lib/Header.svelte';
	import Display from '$lib/Display.svelte';
	import TablePlot from '$lib/TablePlot.svelte';
	import Solver from '$lib/Solver.svelte';
	import Console from '$lib/Console.svelte';
	import Mesh from '$lib/Mesh.svelte';
	import Geometry from '$lib/Geometry.svelte';
	import InitialMag from '$lib/InitialMag.svelte';
	import Parameters from '$lib/Parameters.svelte';
	import Footer from '$lib/Footer.svelte';
	import { onDestroy, onMount } from 'svelte';
	import { baseURL, paused } from '$lib/api';
	import { wsConnect } from '$lib/websocket';

	// let interval = 0;
	onMount(() => {
		try {
			wsConnect();
		} catch (error) {
			console.error('Error connecting to websocket:', error);
		}
		// interval = setInterval(() => {
		// 	if (!$paused) {
		// 		refresh();
		// 	}
		// }, 1000);

		// baseURL.set(window.location.origin);
		baseURL.set('http://localhost:5001');
		// console.log('baseURL:', get(baseURL));
		// fetchEngineState().catch((error) => {
		// 	console.error('Error fetching engine state:', error);
		// });
	});
	onDestroy(() => {
		// clearInterval(interval);
	});
</script>

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
