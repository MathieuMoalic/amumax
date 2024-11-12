<script lang="ts">
	import Slider from '$components/Slider.svelte';
	import { previewState as p } from '$api/incoming/preview';
	import { postXChosenSize } from '$api/outgoing/preview';
</script>

<Slider
	label="X Data Points"
	bind:value={$p.xChosenSize}
	values={$p.xPossibleSizes}
	onChangeFunction={postXChosenSize}
/>
