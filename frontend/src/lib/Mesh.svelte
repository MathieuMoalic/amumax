<script lang="ts">
	import { meshState as m } from '$api/incoming/mesh';
	import { Label, Input, Button, InputAddon, ButtonGroup, Checkbox } from 'flowbite-svelte';
</script>

<section>
	<h2 class="mb-4 text-2xl font-semibold">Mesh</h2>

	<div class="mb-6 grid gap-6 md:grid-cols-3">
		<ButtonGroup class="w-full" size="sm">
			<InputAddon class="w-20">Nx</InputAddon>
			<Input id="input-addon-sm" type="text" placeholder={$m.dx} />
		</ButtonGroup>
		<ButtonGroup class="w-full" size="sm">
			<InputAddon class="w-20">Ny</InputAddon>
			<Input id="input-addon-sm" type="text" placeholder={$m.dy} />
		</ButtonGroup>
		<ButtonGroup class="w-full" size="sm">
			<InputAddon class="w-20">Nz</InputAddon>
			<Input id="input-addon-sm" type="text" placeholder={$m.dz} />
		</ButtonGroup>
	</div>
	<div class="mb-6 grid gap-6 md:grid-cols-3">
		<ButtonGroup class="w-full" size="sm">
			<InputAddon class="w-20">dx</InputAddon>
			<Input id="input-addon-sm" type="text" placeholder={$m.Nx} />
			<InputAddon>m</InputAddon>
		</ButtonGroup>
		<ButtonGroup class="w-full" size="sm">
			<InputAddon class="w-20">dy</InputAddon>
			<Input id="input-addon-sm" type="text" placeholder={$m.Ny} />
			<InputAddon>m</InputAddon>
		</ButtonGroup>
		<ButtonGroup class="w-full" size="sm">
			<InputAddon class="w-20">dz</InputAddon>
			<Input id="input-addon-sm" type="text" placeholder={$m.Nz} />
			<InputAddon>m</InputAddon>
		</ButtonGroup>
	</div>
	<div class="mb-6 grid gap-6 md:grid-cols-3">
		<ButtonGroup class="w-full" size="sm">
			<InputAddon class="w-20">Tx</InputAddon>
			<Input id="input-addon-sm" type="text" placeholder={$m.Tx} />
			<InputAddon>m</InputAddon>
		</ButtonGroup>
		<ButtonGroup class="w-full" size="sm">
			<InputAddon class="w-20">Ty</InputAddon>
			<Input id="input-addon-sm" type="text" placeholder={$m.Ty} />
			<InputAddon>m</InputAddon>
		</ButtonGroup>
		<ButtonGroup class="w-full" size="sm">
			<InputAddon class="w-20">Tz</InputAddon>
			<Input id="input-addon-sm" type="text" placeholder={$m.Tz} />
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
