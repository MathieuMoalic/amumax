<script lang="ts">
	import { onMount } from 'svelte';
	import { meshState } from './api';
	onMount(() => {
		// console.log('meshState', $meshState);
	});
</script>

<span>▾ mesh</span> <br />
<div>
	<table>
		<tbody>
			<tr>
				<td>gridsize: </td>
				<td><input bind:value={$meshState.Nx} /></td>
				<td> × <input bind:value={$meshState.Ny} /></td>
				<td> × <input bind:value={$meshState.Nz} /></td>
				<td> cells </td>
			</tr>
			<tr>
				<td>cellsize: </td>
				<td><input bind:value={$meshState.dx} /> </td>
				<td> × <input bind:value={$meshState.dy} /></td>
				<td> × <input bind:value={$meshState.dz} /></td>
				<td> m<sup>3</sup> </td>
			</tr>
			<tr>
				<td>PBC: </td>
				<td><input bind:value={$meshState.PBCx} /> </td>
				<td> × <input bind:value={$meshState.PBCy} /></td>
				<td> × <input bind:value={$meshState.PBCz} /></td>
				<td> repetitions </td>
			</tr>
			<tr>
				<td>worldsize:</td>
				<td><span>{$meshState.Tx} </span> </td>
				<td> × <span> {$meshState.Ty} </span></td>
				<td> × <span> {$meshState.Tz} </span></td>
				<td> nm<sup>3</sup> </td>
			</tr>
		</tbody>
	</table>

	<button>update</button>
</div>
