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
	import { Input, Label, Button, Select } from 'flowbite-svelte';
	let solvertypes = ['bw_euler', 'euler', 'heun', 'rk4', 'rk23', 'rk45', 'rkf56'];
	let runSeconds = 1e-9;
	let runSteps = 100;
</script>

<section>
	<h2 class="mb-4 text-2xl font-semibold">Solver</h2>

	<div class="grid grid-cols-2 gap-6">
		<!-- Left Column -->
		<div class="grid gap-1">
			<!-- Solver Type -->
			<div class="field">
				<Button for="solver_type" class="h-6 w-28" outline>Solver</Button>
				<Select
					id="solver_type"
					bind:value={$solverState.type}
					on:change={postSolverType}
					class=" block h-6 w-32 rounded-md border border-gray-300 bg-white  text-gray-700 shadow-sm focus:ring focus:ring-blue-500"
				>
					{#each solvertypes as solvertype}
						<option value={solvertype}>{solvertype}</option>
					{/each}
				</Select>
			</div>

			<!-- Run -->
			<div class="field">
				<Button class="h-6 w-28" outline on:click={() => postRun(runSeconds)}>Run</Button>
				<Input class="h-6 w-32" type="number" id="run_seconds" bind:value={runSeconds} />
				<span class="ml-6">s</span>
			</div>

			<!-- Run Steps -->
			<div class="field">
				<Button class="h-6 w-28" outline on:click={() => postSteps(runSteps)}>Run Steps</Button>
				<Input class="h-6 w-32" type="number" id="run_steps" bind:value={runSteps} />
			</div>

			<!-- Relax -->
			<div class="field">
				<Button class="h-6 w-28" outline on:click={postRelax}>Relax</Button>
			</div>

			<!-- Break -->
			<div class="field">
				<Button class="h-6 w-28" outline on:click={postBreak}>Break</Button>
			</div>
		</div>

		<!-- Right Column -->
		<div class="grid gap-1">
			<!-- Steps -->
			<div class="field">
				<Label for="steps" class="h-6 w-32">Steps</Label>
				<Input class="h-6 w-32" type="text" id="steps" placeholder={$solverState.steps} readonly />
			</div>

			<!-- Time -->
			<div class="field">
				<Label for="time" class="h-6 w-32">Time</Label>
				<Input
					class="h-6 w-32"
					type="text"
					id="time"
					placeholder={$solverState.time.toExponential(3)}
					readonly
				/>
				<span class="ml-6">s</span>
			</div>

			<!-- dt -->
			<div class="field">
				<Label for="dt" class="h-6 w-32">dt</Label>
				<Input
					class="h-6 w-32"
					type="text"
					id="dt"
					placeholder={$solverState.dt.toExponential(3)}
					readonly
				/>
				<span class="ml-6">s</span>
			</div>

			<!-- err/step -->
			<div class="field">
				<Label for="err_step" class="h-6 w-32">err/step</Label>
				<Input
					class="h-6 w-32"
					type="text"
					id="err_step"
					placeholder={$solverState.errPerStep.toExponential(3)}
					readonly
				/>
			</div>

			<!-- maxTorque -->
			<div class="field">
				<Label for="maxTorque" class="h-6 w-32">maxTorque</Label>
				<Input
					class="h-6 w-32"
					type="text"
					id="maxTorque"
					placeholder={$solverState.maxTorque.toExponential(3)}
					readonly
				/>
				<span class="ml-6">T</span>
			</div>

			<!-- fixdt -->
			<div class="field">
				<Label for="fixdt" class="h-6 w-32">fixdt</Label>
				<Input
					class="h-6 w-32"
					type="number"
					id="fixdt"
					bind:value={$solverState.fixdt}
					on:change={postFixdt}
				/>
				<span class="ml-6">s</span>
			</div>

			<!-- mindt -->
			<div class="field">
				<Label for="mindt" class="h-6 w-32">mindt</Label>
				<Input
					class="h-6 w-32"
					type="number"
					id="mindt"
					bind:value={$solverState.mindt}
					on:change={postMindt}
				/>
				<span class="ml-6">s</span>
			</div>

			<!-- maxdt -->
			<div class="field">
				<Label for="maxdt" class="h-6 w-32">maxdt</Label>
				<Input
					class="h-6 w-32"
					type="number"
					id="maxdt"
					bind:value={$solverState.maxdt}
					on:change={postMaxdt}
				/>
				<span class="ml-6">s</span>
			</div>

			<!-- maxerr -->
			<div class="field">
				<Label for="maxerr" class="h-6 w-32">maxerr</Label>
				<Input
					class="h-6 w-32"
					type="number"
					id="maxerr"
					bind:value={$solverState.maxerr}
					on:change={postMaxerr}
				/>
			</div>
		</div>
	</div>
</section>

<style>
	section {
		grid-area: solver;
	}
	.field {
		@apply grid grid-cols-3 gap-7;
	}
</style>
