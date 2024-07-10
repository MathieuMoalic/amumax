<script lang="ts">
	import { solverState } from '$api/incoming/solver';
	import {
		postSolverType,
		postRun,
		postSteps,
		postRelax,
		postBreak,
		postFixdt,
		postMindt,
		postMaxdt,
		postMaxerr
	} from '$api/outgoing/solver';

	let solvertypes = ['bw_euler', 'euler', 'heun', 'rk4', 'rk23', 'rk45', 'rkf56'];
	let runSeconds = 1e-9;
	let runSteps = 100;
</script>

<section>
	<h2>Solver</h2>
	<div class="flex-container">
		<div class="grid-container">
			<div>
				<div class="header">Solver</div>
				<div class="value">
					<select bind:value={$solverState.type} on:change={postSolverType}>
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
				<div class="header">
					<button
						on:click={() => {
							postRun(runSeconds);
						}}>Run</button
					>
				</div>
				<div class="value">
					<input type="number" bind:value={runSeconds} />
				</div>
				<div class="unit">s</div>
			</div>
			<div>
				<div class="header">
					<button
						on:click={() => {
							postSteps(runSteps);
						}}>Run Steps</button
					>
				</div>
				<div class="value"><input type="number" bind:value={runSteps} /></div>
				<div class="unit"></div>
			</div>
			<div>
				<div class="header"><button on:click={postRelax}>Relax</button></div>
				<div class="value"></div>
				<div class="unit"></div>
			</div>
			<div>
				<div class="header"><button on:click={postBreak}>Break</button></div>
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
				<div class="value">
					<input type="number" bind:value={$solverState.fixdt} on:change={postFixdt} />
				</div>
				<div class="unit">s</div>
			</div>
			<div>
				<div class="header">mindt:</div>
				<input type="number" bind:value={$solverState.mindt} on:change={postMindt} />
				<div class="unit">s</div>
			</div>
			<div>
				<div class="header">maxdt:</div>
				<input type="number" bind:value={$solverState.maxdt} on:change={postMaxdt} />
				<div class="unit">s</div>
			</div>
			<div>
				<div class="header">maxerr:</div>
				<input type="number" bind:value={$solverState.maxerr} on:change={postMaxerr} />
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
