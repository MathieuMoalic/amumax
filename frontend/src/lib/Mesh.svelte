<script lang="ts">
	import { meshState as m } from '$api/incoming/mesh';
	import { Input, InputAddon, ButtonGroup } from 'flowbite-svelte';
</script>

<section>
	<h2 class="mb-4 text-2xl font-semibold">Mesh</h2>

	<div class="mb-6 grid gap-6 md:grid-cols-3">
		<ButtonGroup class="h-11 w-full">
			<InputAddon class="w-28 !bg-transparent">dx</InputAddon>
			<Input type="text" placeholder={$m.dx.toPrecision(8)} />
			<InputAddon class="w-14 !bg-transparent">m</InputAddon>
		</ButtonGroup>
		<ButtonGroup class="h-11 w-full">
			<InputAddon class="w-28 !bg-transparent">dy</InputAddon>
			<Input type="text" placeholder={$m.dy.toPrecision(8)} />
			<InputAddon class="w-14 !bg-transparent">m</InputAddon>
		</ButtonGroup>
		<ButtonGroup class="h-11 w-full">
			<InputAddon class="w-28 !bg-transparent">dz</InputAddon>
			<Input type="text" placeholder={$m.dz.toPrecision(8)} />
			<InputAddon class="w-14 !bg-transparent">m</InputAddon>
		</ButtonGroup>
	</div>

	<div class="mb-6 grid gap-6 md:grid-cols-3">
		<ButtonGroup class="h-11 w-full">
			<InputAddon class="w-28 !bg-transparent">Nx</InputAddon>
			<Input type="text" placeholder=" {$m.Nx}" />
			<InputAddon class="w-14 !bg-transparent"></InputAddon>
		</ButtonGroup>
		<ButtonGroup class="h-11 w-full">
			<InputAddon class="w-28 !bg-transparent">Ny</InputAddon>
			<Input type="text" placeholder=" {$m.Ny}" />
			<InputAddon class="w-14 !bg-transparent"></InputAddon>
		</ButtonGroup>
		<ButtonGroup class="h-11 w-full">
			<InputAddon class="w-28 !bg-transparent">Nz</InputAddon>
			<Input type="text" placeholder=" {$m.Nz}" />
			<InputAddon class="w-14 !bg-transparent"></InputAddon>
		</ButtonGroup>
	</div>

	<div class="mb-6 grid gap-6 md:grid-cols-3">
		<ButtonGroup class="h-11 w-full">
			<InputAddon class="w-28 !bg-transparent">Tx</InputAddon>
			<Input type="text" placeholder={$m.Tx.toExponential(6)} />
			<InputAddon class="w-14 !bg-transparent">m</InputAddon>
		</ButtonGroup>
		<ButtonGroup class="h-11 w-full">
			<InputAddon class="w-28 !bg-transparent">Ty</InputAddon>
			<Input type="text" placeholder={$m.Ty.toExponential(6)} />
			<InputAddon class="w-14 !bg-transparent">m</InputAddon>
		</ButtonGroup>
		<ButtonGroup class="h-11 w-full">
			<InputAddon class="w-28 !bg-transparent">Tz</InputAddon>
			<Input type="text" placeholder={$m.Tz.toExponential(6)} />
			<InputAddon class="w-14 !bg-transparent">m</InputAddon>
		</ButtonGroup>
	</div>

	<div class="mb-6 grid gap-6 md:grid-cols-3">
		<ButtonGroup class="h-11 w-full">
			<InputAddon class="w-28 !bg-transparent">PBCx</InputAddon>
			<Input type="text" placeholder=" {$m.PBCx}" />
			<InputAddon class="w-14 !bg-transparent"></InputAddon>
		</ButtonGroup>
		<ButtonGroup class="h-11 w-full">
			<InputAddon class="w-28 !bg-transparent">PBCy</InputAddon>
			<Input type="text" placeholder=" {$m.PBCy}" />
			<InputAddon class="w-14 !bg-transparent"></InputAddon>
		</ButtonGroup>
		<ButtonGroup class="h-11 w-full">
			<InputAddon class="w-28 !bg-transparent">PBCz</InputAddon>
			<Input type="text" placeholder=" {$m.PBCz}" />
			<InputAddon class="w-14 !bg-transparent"></InputAddon>
		</ButtonGroup>
	</div>
</section>

<style>
	section {
		grid-area: mesh;
	}
</style>
