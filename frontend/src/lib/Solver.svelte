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
		$headerState.status = 'running';
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

<section>
	<h2>Solvers</h2>
	<div>
		Type:{$solverState.type}
		<select bind:value={selectedSolver} on:change={changeSolver}>
			{#each solvertypes as solvertype}
				<option value={solvertype}>
					{solvertype}
				</option>
			{/each}
		</select>
	</div>
	<div>
		<button on:click={run}>Run</button>
		<input bind:value={runSeconds} />s
	</div>

	<div>
		<button on:click={steps}>Steps</button>
		<input bind:value={runSteps} />
	</div>
	<div>
		<button on:click={relax}>Relax</button>
	</div>
	<div>
		<button on:click={pbreak}>Break</button>
	</div>
	<div>
		step: {$solverState.steps}
	</div>
	<div>
		time: {$solverState.time.toExponential(3)} s
	</div>
	<div>
		dt: {$solverState.dt.toExponential(3)} s
	</div>
	<div>
		err/step: {$solverState.errPerStep.toExponential(3)}
	</div>
	<div>
		MaxTorque: {$solverState.maxTorque.toExponential(3)}
	</div>
	<div>
		fixdt: <input placeholder=" {$solverState.fixdt.toExponential(3)}" />s
	</div>
	<div>
		mindt: <input placeholder=" {$solverState.mindt.toExponential(3)}" />s
	</div>
	<div>
		maxdt: <input placeholder=" {$solverState.maxdt.toExponential(3)}" />s
	</div>
	<div>
		maxerr: <input placeholder=" {$solverState.maxerr.toExponential(3)}" />/step
	</div>
</section>

<style>
	section {
		grid-area: solver;
	}
	input {
		size: 8;
	}
</style>
