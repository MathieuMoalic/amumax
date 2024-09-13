<script lang="ts">
	import { meshState } from '$api/incoming/mesh';
	import { previewState as p } from '$api/incoming/preview';
	import { postLayer } from '$api/outgoing/preview';
	import { Range, Label } from 'flowbite-svelte';
</script>

<div class="m-3 grid grid-cols-5 gap-2">
	<div class="col-span-1 flex items-center justify-center">
		<Label>Layer: {$p.layer}</Label>
	</div>
	<div class="col-span-4 flex items-center justify-center">
		<Range
			id="range-steps"
			min="0"
			max={$meshState.Nz - 1}
			bind:value={$p.layer}
			step="1"
			on:change={postLayer}
		/>
	</div>
</div>
