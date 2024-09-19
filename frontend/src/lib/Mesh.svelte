<script lang="ts">
	import { meshState as m } from '$api/incoming/mesh';
	import { Input, InputAddon, ButtonGroup } from 'flowbite-svelte';
</script>

<section>
	<h2 class="mb-4 text-2xl font-semibold">Mesh</h2>

	<div class="mb-6 grid gap-6 md:grid-cols-3">
		<ButtonGroup class="w-full" size="sm">
			<InputAddon class="w-20">dx</InputAddon>
			<Input id="input-addon-sm" type="text" placeholder={$m.dx.toPrecision(8)} />
			<InputAddon>m</InputAddon>
		</ButtonGroup>
		<ButtonGroup class="w-full" size="sm">
			<InputAddon class="w-20">dy</InputAddon>
			<Input id="input-addon-sm" type="text" placeholder={$m.dy.toPrecision(8)} />
			<InputAddon>m</InputAddon>
		</ButtonGroup>
		<ButtonGroup class="w-full" size="sm">
			<InputAddon class="w-20">dz</InputAddon>
			<Input id="input-addon-sm" type="text" placeholder={$m.dz.toPrecision(8)} />
			<InputAddon>m</InputAddon>
		</ButtonGroup>
	</div>
	<div class="mb-6 grid gap-6 md:grid-cols-3">
		<ButtonGroup class="w-full" size="sm">
			<InputAddon class="w-20">Nx</InputAddon>
			<Input id="input-addon-sm" type="text" placeholder={$m.Nx} />
		</ButtonGroup>
		<ButtonGroup class="w-full" size="sm">
			<InputAddon class="w-20">Ny</InputAddon>
			<Input id="input-addon-sm" type="text" placeholder={$m.Ny} />
		</ButtonGroup>
		<ButtonGroup class="w-full" size="sm">
			<InputAddon class="w-20">Nz</InputAddon>
			<Input id="input-addon-sm" type="text" placeholder={$m.Nz} />
		</ButtonGroup>
	</div>
	<div class="mb-6 grid gap-6 md:grid-cols-3">
		<ButtonGroup class="w-full" size="sm">
			<InputAddon class="w-20">Tx</InputAddon>
			<Input id="input-addon-sm" type="text" placeholder={$m.Tx.toExponential(6)} />
			<InputAddon>m</InputAddon>
		</ButtonGroup>
		<ButtonGroup class="w-full" size="sm">
			<InputAddon class="w-20">Ty</InputAddon>
			<Input id="input-addon-sm" type="text" placeholder={$m.Ty.toExponential(6)} />
			<InputAddon>m</InputAddon>
		</ButtonGroup>
		<ButtonGroup class="w-full" size="sm">
			<InputAddon class="w-20">Tz</InputAddon>
			<Input id="input-addon-sm" type="text" placeholder={$m.Tz.toExponential(6)} />
			<InputAddon>m</InputAddon>
		</ButtonGroup>
	</div>
	<div class="mb-6 grid gap-6 md:grid-cols-3">
		<ButtonGroup class="w-full" size="sm">
			<InputAddon class="w-20">PBCx</InputAddon>
			<Input id="input-addon-sm" type="text" placeholder={$m.PBCx} />
		</ButtonGroup>
		<ButtonGroup class="w-full" size="sm">
			<InputAddon class="w-20">PBCy</InputAddon>
			<Input id="input-addon-sm" type="text" placeholder={$m.PBCy} />
		</ButtonGroup>
		<ButtonGroup class="w-full" size="sm">
			<InputAddon class="w-20">PBCz</InputAddon>
			<Input id="input-addon-sm" type="text" placeholder={$m.PBCz} />
		</ButtonGroup>
	</div>
</section>

<style>
	section {
		grid-area: mesh;
	}
</style>
