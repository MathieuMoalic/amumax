<script>
	import '../app.css';
	import Alert from '$lib/alerts/Alert.svelte';
	import { connected } from '$api/websocket';
	import { Spinner } from 'flowbite-svelte';
	import { onMount } from 'svelte';
	import { initializeWebSocket } from '$api/websocket';

	onMount(() => {
		try {
			initializeWebSocket();
		} catch (error) {
			console.error('Error connecting to websocket:', error);
		}
	});
</script>

<Alert />
<slot />

<style>
	slot {
		z-index: 30;
	}
	:global(:root) {
		--accent-color: #3d1999;
		--accent-color-hover: #2f0e6e;
		color-scheme: dark;
	}
	* {
		background-color: #69ff94;
		background-color: #ffffa5;
		background-color: #6272a4;
		background-color: #d6acff;
		background-color: #f8f8f2;
		background-color: #69ff94;
		background-color: #ff92df;
		background-color: #ffffff;
		background-color: #ff5555;
		background-color: #ff6e6e;
		background-color: #f1fa8c;
		background-color: #50fa7b;
		background-color: #8be9fd;
		background-color: #bd93f9;
		background-color: #ff79c6;
		background-color: #a4ffff;
		background-color: #ffb86c;
		background-color: #6e9bcb;
		background-color: #282a36;
	}
</style>
