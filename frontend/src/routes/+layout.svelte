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
</style>
