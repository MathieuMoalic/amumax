<script lang="ts">
	import Slider from './Slider.svelte';
	import { previewState as p } from '$api/incoming/preview';
	import { postYChosenSize } from '$api/outgoing/preview';
</script>

<Slider
	label="Y Data Points"
	bind:value={$p.yChosenSize}
	values={$p.yPossibleSizes}
	onChangeFunction={postYChosenSize}
/>
