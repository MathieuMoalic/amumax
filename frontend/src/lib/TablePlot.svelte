<script lang="ts">
	import { tablePlotState, refreshTablePlot, TablePlotUrl, tablePlotx, tablePloty } from './api';
</script>

<section>
	<h2>Table Plot</h2>
	<div>
		<br />
		<div>
			<p>
				Auto Save Interval: <input bind:value={$tablePlotState.autoSaveInterval} /> s
			</p>
			<b>
				x:
				<select bind:value={$tablePlotx} on:change={refreshTablePlot}>
					{#each $tablePlotState.columns as q}
						<option value={q}>{q}</option>
					{/each}
				</select>
				y:
				<select bind:value={$tablePloty} on:change={refreshTablePlot}>
					{#each $tablePlotState.columns as q}
						<option value={q}>{q}</option>
					{/each}
				</select>
			</b>
			<br />
			<p><span></span></p>
			<img src={$TablePlotUrl} alt="TablePlot" />
		</div>
	</div>
</section>

<style>
	section {
		grid-area: tableplot;
	}
</style>
