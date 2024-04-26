<script lang="ts">
	import quantities from './quantities';
	import { refreshImage, imageUrl, renderLayer, renderComp, renderQuant } from './api';
</script>

<section>
	<span> ▾ display </span>
	<br />
	<div>
		<p>
			Quantity:
			<select bind:value={$renderQuant} on:change={refreshImage}>
				{#each quantities as q}
					<option value={q}>{q}</option>
				{/each}
			</select>
			<select bind:value={$renderComp}>
				{#each ['All', 'x', 'y', 'z'] as q}
					<option value={q}>{q}</option>
				{/each}
			</select>
			<br />
		</p>
		Slice:<input type="range" min="0" max="0" bind:value={$renderLayer} />
		<span>{$renderLayer}</span>
		<img src={$imageUrl} alt="Display" />
	</div>
</section>
