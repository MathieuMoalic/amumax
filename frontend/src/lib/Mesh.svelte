<script lang="ts">
	import { meshState } from '$api/incoming/mesh';
	import { postMesh } from '$api/outgoing/mesh';
</script>

<section>
	<h2>Mesh</h2>
	<div class="grid-container">
		<div class="header">Nx:</div>
		<!-- <div><input bind:value={$meshState.Nx} on:change={postMesh} /></div> -->
		<div><span>{$meshState.Nx}</span></div>

		<div>×</div>

		<div class="header">Ny:</div>
		<!-- <div><input bind:value={$meshState.Ny} on:change={postMesh} /></div> -->
		<div><span>{$meshState.Ny}</span></div>

		<div>×</div>

		<div class="header">Nz:</div>
		<!-- <div><input bind:value={$meshState.Nz} on:change={postMesh} /></div> -->
		<div><span>{$meshState.Nz}</span></div>

		<div>cells</div>

		<div class="header">dx:</div>
		<!-- <div><input bind:value={$meshState.dx} on:change={postMesh} /></div> -->
		<div><span>{$meshState.dx}</span></div>

		<div>×</div>

		<div class="header">dy:</div>
		<!-- <div><input bind:value={$meshState.dy} on:change={postMesh} /></div> -->
		<div><span>{$meshState.dy}</span></div>

		<div>×</div>

		<div class="header">dz:</div>
		<!-- <div><input bind:value={$meshState.dz} on:change={postMesh} /></div> -->
		<div><span>{$meshState.dz}</span></div>

		<div>m<sup>3</sup></div>

		<div class="header">PBCx:</div>
		<!-- <div><input bind:value={$meshState.PBCx} on:change={postMesh} /></div> -->
		<div><span>{$meshState.PBCx}</span></div>

		<div>×</div>
		<div class="header">PBCy:</div>
		<!-- <div><input bind:value={$meshState.PBCy} on:change={postMesh} /></div> -->
		<div><span>{$meshState.PBCy}</span></div>

		<div>×</div>
		<div class="header">PBCz:</div>
		<!-- <div><input bind:value={$meshState.PBCz} on:change={postMesh} /></div> -->
		<div><span>{$meshState.PBCz}</span></div>

		<div>repetitions</div>

		<div class="header">Tx:</div>
		<div><span>{$meshState.Tx.toExponential(3)}</span></div>
		<div>×</div>
		<div class="header">Ty:</div>
		<div><span>{$meshState.Ty.toExponential(3)}</span></div>
		<div>×</div>
		<div class="header">Tz:</div>
		<div><span>{$meshState.Tz.toExponential(3)}</span></div>
		<div>m<sup>3</sup></div>
	</div>
</section>

<style>
	.grid-container {
		display: grid;
		grid-template-columns: auto auto auto auto auto auto auto auto auto;
		gap: 10px;
		align-items: center;
		padding: 10px;
	}
	.header {
		text-align: right;
		padding-right: 10px;
	}
	section {
		grid-area: mesh;
	}
</style>
