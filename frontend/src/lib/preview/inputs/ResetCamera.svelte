<script lang="ts">
	import { previewState as p } from '$api/incoming/preview';
	import { Button } from 'flowbite-svelte';
	import { resetCamera } from '../preview3D';
	let isDisabled: boolean;

	$: isDisabled = $p.nComp != 3;
</script>

<Button
	outline
	class="h-11 w-full justify-between {isDisabled ? 'cursor-not-allowed opacity-50' : ''}"
	on:click={resetCamera}
	disabled={isDisabled}>Reset Camera</Button
>
