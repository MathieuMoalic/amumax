<script lang="ts">
	import { meshState } from '$api/incoming/mesh';
</script>

<section>
	<h2>Mesh</h2>
	<div class="grid-container">
		<div class="header">Nx:</div>
		<div><input placeholder=" {$meshState.Nx}" /></div>
		<div>×</div>
		<div class="header">Ny:</div>
		<div><input placeholder=" {$meshState.Ny}" /></div>
		<div>×</div>
		<div class="header">Nz:</div>
		<div><input placeholder=" {$meshState.Nz}" /></div>
		<div>cells</div>

		<div class="header">dx:</div>
		<div><input placeholder=" {$meshState.dx.toExponential(3)}" /></div>
		<div>×</div>
		<div class="header">dy:</div>
		<div><input placeholder=" {$meshState.dy.toExponential(3)}" /></div>
		<div>×</div>
		<div class="header">dz:</div>
		<div><input placeholder=" {$meshState.dz.toExponential(3)}" /></div>
		<div>m<sup>3</sup></div>

		<div class="header">PBCx:</div>
		<div><input placeholder=" {$meshState.PBCx}" /></div>
		<div>×</div>
		<div class="header">PBCy:</div>
		<div><input placeholder=" {$meshState.PBCy}" /></div>
		<div>×</div>
		<div class="header">PBCz:</div>
		<div><input placeholder=" {$meshState.PBCz}" /></div>
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
	input {
		width: 100%;
		box-sizing: border-box;
	}
	.header {
		text-align: right;
		padding-right: 10px;
	}
	section {
		grid-area: mesh;
	}
</style>
