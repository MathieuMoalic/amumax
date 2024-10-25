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
		postMaxerr,
		postMinimize
	} from '$api/outgoing/solver';
	import { Input } from 'flowbite-svelte';
	import { Button, Dropdown, DropdownItem, ButtonGroup, InputAddon } from 'flowbite-svelte';
	import { ChevronDownOutline } from 'flowbite-svelte-icons';

	let solvertypes = ['bw_euler', 'euler', 'heun', 'rk4', 'rk23', 'rk45', 'rkf56'];
	let runSeconds = '1e-9';
	let runSteps = '100';
	let dropdownOpen = false;
</script>

<section>
	<h2 class="mb-6 text-2xl font-semibold">Solver</h2>

	<div class="grid gap-6 md:grid-cols-2">
		<!-- Left Column -->
		<div class="space-y-4">
			<!-- Solver Type -->
			<div>
				<Button outline color="primary" class="flex w-full items-center justify-between">
					<span>Solver: {$solverState.type}</span>

					<ChevronDownOutline class="h-5 w-5 text-gray-500" />
				</Button>
				<Dropdown
					bind:open={dropdownOpen}
					class="w-56"
					placement="bottom-start"
					on:click={() => (dropdownOpen = false)}
				>
					{#each solvertypes as solvertype}
						<DropdownItem
							class="flex items-center gap-2 text-base font-semibold"
							on:click={() => {
								postSolverType(solvertype);
								dropdownOpen = false;
							}}
						>
							{solvertype}
						</DropdownItem>
					{/each}
				</Dropdown>
			</div>

			<!-- Run -->
			<div class="flex items-center space-x-2">
				<Button outline on:click={() => postRun(runSeconds)}>Run</Button>
				<ButtonGroup class="h-11 w-full">
					<Input
						class="w-full"
						id="run_seconds"
						bind:value={runSeconds}
						placeholder="Time in seconds"
						on:change={() => postRun(runSeconds)}
					/>
					<InputAddon class="w-16 !bg-transparent">s</InputAddon>
				</ButtonGroup>
			</div>

			<!-- Run Steps -->
			<div class="flex items-center space-x-2">
				<Button
					class="inline-flex h-11 items-center whitespace-nowrap px-4 py-2"
					outline
					on:click={() => postSteps(runSteps)}>Run Steps</Button
				>
				<Input
					class="w-full"
					id="run_steps"
					bind:value={runSteps}
					placeholder="Number of steps"
					on:change={() => postRun(runSteps)}
				/>
			</div>

			<!-- Relax -->
			<div>
				<Button outline class="w-full" on:click={postRelax}>Relax</Button>
			</div>

			<!-- Relax -->
			<div>
				<Button outline class="w-full" on:click={postMinimize}>Minimize</Button>
			</div>

			<!-- Break -->
			<div>
				<Button outline class="w-full" on:click={postBreak}>Break</Button>
			</div>
		</div>

		<!-- Right Column -->
		<div class="space-y-4">
			<!-- Steps -->
			<div class="flex items-center space-x-2">
				<ButtonGroup class="h-11 w-full">
					<InputAddon class="w-44 !bg-transparent">Steps</InputAddon>
					<Input class="w-full" value={$solverState.steps} readonly />
					<InputAddon class="w-24 !bg-transparent"></InputAddon>
				</ButtonGroup>
			</div>

			<!-- Time -->
			<div class="flex items-center space-x-2">
				<ButtonGroup class="h-11 w-full">
					<InputAddon class="w-44 !bg-transparent">Time</InputAddon>
					<Input class="w-full" type="text" value={$solverState.time.toExponential(3)} readonly />
					<InputAddon class="w-24 !bg-transparent">s</InputAddon>
				</ButtonGroup>
			</div>

			<!-- dt -->
			<div class="flex items-center space-x-2">
				<ButtonGroup class="h-11 w-full">
					<InputAddon class="w-44 !bg-transparent">dt</InputAddon>
					<Input class="w-full" type="text" value={$solverState.dt.toExponential(3)} readonly />
					<InputAddon class="w-24 !bg-transparent">s</InputAddon>
				</ButtonGroup>
			</div>

			<!-- err/step -->
			<div class="flex items-center space-x-2">
				<ButtonGroup class="h-11 w-full">
					<InputAddon class="w-44 !bg-transparent">Err/step</InputAddon>
					<Input
						class="w-full"
						type="text"
						value={$solverState.errPerStep.toExponential(3)}
						readonly
					/>
					<InputAddon class="w-24 !bg-transparent"></InputAddon>
				</ButtonGroup>
			</div>

			<!-- maxTorque -->
			<div class="flex items-center space-x-2">
				<ButtonGroup class="h-11 w-full">
					<InputAddon class="w-44 !bg-transparent">Max Torque</InputAddon>
					<Input
						class="w-full"
						type="text"
						value={$solverState.maxTorque.toExponential(3)}
						readonly
					/>
					<InputAddon class="w-24 !bg-transparent">T</InputAddon>
				</ButtonGroup>
			</div>

			<!-- fixdt -->
			<div class="flex items-center space-x-2">
				<ButtonGroup class="h-11 w-full">
					<InputAddon class="w-44 !bg-transparent">Fixdt</InputAddon>
					<Input
						class="w-full"
						type="number"
						bind:value={$solverState.fixdt}
						on:change={postFixdt}
					/>
					<InputAddon class="w-24 !bg-transparent">s</InputAddon>
				</ButtonGroup>
			</div>

			<!-- mindt -->
			<div class="flex items-center space-x-2">
				<ButtonGroup class="h-11 w-full">
					<InputAddon class="w-44 !bg-transparent">Mindt</InputAddon>
					<Input
						class="w-full"
						type="number"
						bind:value={$solverState.mindt}
						on:change={postMindt}
					/>
					<InputAddon class="w-24 !bg-transparent">s</InputAddon>
				</ButtonGroup>
			</div>

			<!-- maxdt -->
			<div class="flex items-center space-x-2">
				<ButtonGroup class="h-11 w-full">
					<InputAddon class="w-44 !bg-transparent">Maxdt</InputAddon>
					<Input
						class="w-full"
						type="number"
						bind:value={$solverState.maxdt}
						on:change={postMaxdt}
					/>
					<InputAddon class="w-24 !bg-transparent">s</InputAddon>
				</ButtonGroup>
			</div>

			<!-- maxerr -->
			<div class="flex items-center space-x-2">
				<ButtonGroup class="h-11 w-full">
					<InputAddon class="w-44 !bg-transparent">MaxErr</InputAddon>
					<Input
						class="w-full"
						type="number"
						bind:value={$solverState.maxerr}
						on:change={postMaxerr}
					/>
					<InputAddon class="w-24 !bg-transparent"></InputAddon>
				</ButtonGroup>
			</div>
		</div>
	</div>
</section>

<style>
	section {
		grid-area: solver;
	}
</style>
