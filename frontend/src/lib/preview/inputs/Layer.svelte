<script lang="ts">
	import Slider from '$components/Slider.svelte';
	import { meshState } from '$api/incoming/mesh';
	import { previewState as p } from '$api/incoming/preview';
	import { postLayer } from '$api/outgoing/preview';
	let isDisabled: boolean;

	$: isDisabled = $meshState.Nz < 2;
	$: values = Array.from({ length: $meshState.Nz }, (_, i) => i);
	$: layer = $p.layer;
</script>

<Slider label="Z Layer" bind:value={layer} {values} onChangeFunction={postLayer} {isDisabled} />
