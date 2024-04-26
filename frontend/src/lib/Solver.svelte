<script lang="ts">
	import {
		solverState,
		postSolverType,
		postRun,
		postSteps,
		postRelax,
		postBreak,
		headerState
	} from './api';
	let solvertypes = ['bw_euler', 'euler', 'heun', 'rk4', 'rk23', 'rk45', 'rkf56'];
	let selectedSolver = $solverState.type;
	function changeSolver() {
		postSolverType(selectedSolver).catch((error) => {
			console.error('Error posting solver type:', error);
		});
	}
	let runSeconds = 1e-9;
	let runSteps = 100;
	function run() {
		$headerState.Status = 'running';
		postRun(runSeconds).catch((error) => {
			console.error('Error posting run:', error);
		});
	}
	function steps() {
		postSteps(runSteps).catch((error) => {
			console.error('Error posting steps:', error);
		});
	}
	function relax() {
		postRelax().catch((error) => {
			console.error('Error posting relax:', error);
		});
	}
	function pbreak() {
		postBreak().catch((error) => {
			console.error('Error posting break:', error);
		});
	}
</script>

<div>
	<span>▾ solver</span> <br />
	<div>
		Type: <span>{$solverState.type}</span>
		<select bind:value={selectedSolver} on:change={changeSolver}>
			{#each solvertypes as solvertype}
				<option value={solvertype}>
					{solvertype}
				</option>
			{/each}
		</select>
		<table>
			<tbody>
				<tr>
					<td>
						<table>
							<tbody>
								<tr>
									<td> <button on:click={run}>Run</button></td>
									<td> <input bind:value={runSeconds} />s</td>
								</tr>
								<tr>
									<td> <button on:click={steps}>Steps</button></td>
									<td>
										<input bind:value={runSteps} />
									</td>
								</tr>
								<tr>
									<td> <button on:click={relax}>Relax</button></td>
									<td></td>
								</tr>
								<tr>
									<td> <button on:click={pbreak}>Break</button></td>
									<td></td>
								</tr>
							</tbody>
						</table>
					</td>
					<td> &nbsp; &nbsp; &nbsp; &nbsp; </td>
					<td>
						<table>
							<tbody>
								<tr>
									<td>step: </td>
									<td>{$solverState.steps} </td>
								</tr>
								<tr>
									<td>time: </td>
									<td><span>{$solverState.time.toExponential(3)}</span> s </td>
								</tr>
								<tr>
									<td>dt: </td>
									<td><span>{$solverState.dt.toExponential(3)}</span> s </td>
								</tr>
								<tr>
									<td>err/step: </td>
									<td><span>{$solverState.errPerStep.toExponential(3)}</span> </td>
								</tr>
								<tr>
									<td>MaxTorque:</td>
									<td><span>{$solverState.maxTorque.toExponential(3)}</span> </td>
								</tr>
							</tbody>
						</table>
					</td>
					<td> &nbsp; &nbsp; &nbsp; &nbsp; </td>
					<td>
						<table>
							<tbody>
								<tr>
									<td>fixdt: </td>
									<td
										><input placeholder=" {$solverState.fixdt.toExponential(3)}" />
										s
									</td>
								</tr>
								<tr>
									<td>mindt: </td>
									<td
										><input placeholder=" {$solverState.mindt.toExponential(3)}" />
										s
									</td>
								</tr>
								<tr>
									<td>maxdt: </td>
									<td
										><input placeholder=" {$solverState.maxdt.toExponential(3)}" />
										s
									</td>
								</tr>
								<tr>
									<td>maxerr: </td>
									<td><input placeholder=" {$solverState.maxerr.toExponential(3)}" />/step</td>
								</tr>
							</tbody>
						</table>
					</td>
				</tr>
			</tbody>
		</table>
	</div>
</div>

<style>
	input {
		size: 8;
	}
</style>
