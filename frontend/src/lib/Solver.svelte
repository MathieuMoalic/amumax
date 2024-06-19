<script lang="ts">
	import { solverState } from '$api/incoming/solver';
	import { postSolverType, postRun, postSteps, postRelax, postBreak } from '$api/outgoing/solver';
	import { headerState } from '$api/incoming/header';
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
	<h2>Solver</h2>
	<div class="flex-container">
		<div class="grid-container">
			<div>
				<div class="header">Solver</div>
				<div class="value">
					<select bind:value={selectedSolver} on:change={changeSolver}>
						{#each solvertypes as solvertype}
							<option value={solvertype}>
								{solvertype}
							</option>
						{/each}
					</select>
				</div>
				<div class="unit"></div>
			</div>
			<div>
				<div class="header"><button on:click={run}>Run</button></div>
				<div class="value">
					<input bind:value={runSeconds} />
				</div>
				<div class="unit">s</div>
			</div>
			<div>
				<div class="header"><button on:click={steps}>Run Steps</button></div>
				<div class="value"><input bind:value={runSteps} /></div>
				<div class="unit"></div>
			</div>
			<div>
				<div class="header"><button on:click={relax}>Relax</button></div>
				<div class="value"></div>
				<div class="unit"></div>
			</div>
			<div>
				<div class="header"><button on:click={pbreak}>Break</button></div>
				<div class="value"></div>
				<div class="unit"></div>
			</div>
		</div>
		<div class="grid-container">
			<div>
				<div class="header">steps:</div>
				<div class="value">{$solverState.steps}</div>
				<div class="unit"></div>
			</div>
			<div>
				<div class="header">time:</div>
				<div class="value">{$solverState.time.toExponential(3)}</div>
				<div class="unit">s</div>
			</div>
			<div>
				<div class="header">dt:</div>
				<div class="value">{$solverState.dt.toExponential(3)}</div>
				<div class="unit">s</div>
			</div>
			<div>
				<div class="header">err/step:</div>
				<div class="value">{$solverState.errPerStep.toExponential(3)}</div>
				<div class="unit"></div>
			</div>
			<div>
				<div class="header">maxTorque:</div>
				<div class="value">{$solverState.maxTorque.toExponential(3)}</div>
				<div class="unit">T</div>
			</div>
			<div>
				<div class="header">fixdt:</div>
				<div class="value"><input placeholder=" {$solverState.fixdt.toExponential(3)}" /></div>
				<div class="unit">s</div>
			</div>
			<div>
				<div class="header">mindt:</div>
				<div class="value"><input placeholder=" {$solverState.mindt.toExponential(3)}" /></div>
				<div class="unit">s</div>
			</div>
			<div>
				<div class="header">maxdt:</div>
				<div class="value"><input placeholder=" {$solverState.maxdt.toExponential(3)}" /></div>
				<div class="unit">s</div>
			</div>
			<div>
				<div class="header">maxerr:</div>
				<div class="value"><input placeholder=" {$solverState.maxerr.toExponential(3)}" /></div>
				<div class="unit">/step</div>
			</div>
		</div>
	</div>
</section>

<style>
	section {
		grid-area: solver;
	}
	.flex-container {
		display: flex;
		gap: 20px; /* Adjust the gap between the grid containers as needed */
	}
	.grid-container {
		display: grid;
		grid-template-columns: 100px 100px auto;
		gap: 10px 10px; /* Adjust the gap as needed */
		align-items: center;
	}

	.grid-container > div {
		display: contents; /* Allows the divs to participate in the grid layout */
	}

	.header {
		text-align: right;
		padding-right: 10px;
	}

	.value {
		text-align: left;
	}

	.unit {
		text-align: left;
		padding-left: 5px;
	}

	input {
		width: 100px;
		box-sizing: border-box;
	}
</style>
