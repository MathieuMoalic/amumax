<script lang="ts">
	import { parametersState } from './api';
	let selectedRegion = 0;
</script>

<section>
	<h2>Parameters</h2>
	<div>
		Region:
		<select bind:value={selectedRegion} on:change={() => {}}>
			{#each $parametersState.regions as region}
				<option value={region}>{region}</option>
			{/each}
		</select>
		<div class="grid-container">
			<div class="header">Aex:</div>
			<input placeholder=" {$parametersState.Aex}" />
			<div>J/m</div>
			<div class="description">Exchange stiffness</div>

			<div class="header">Msat:</div>
			<input placeholder=" {$parametersState.Msat}" />
			<div>A/m</div>
			<div class="description">Saturation magnetization</div>

			<div class="header">alpha:</div>
			<input placeholder=" {$parametersState.alpha}" />
			<div></div>
			<div class="description">Landau-Lifshitz damping constant</div>

			<div class="header">B_ext:</div>
			<input placeholder=" {$parametersState.B_ext}" />
			<div>T</div>
			<div class="description">Externally applied field</div>

			<div class="header">anisU:</div>
			<input placeholder=" {$parametersState.anisU}" />
			<div>Vector(x,y,z)</div>
			<div class="description">Uniaxial anisotropy direction</div>

			<div class="header">Ku1:</div>
			<input placeholder=" {$parametersState.Ku1}" />
			<div>J/m<sup>3</sup></div>
			<div class="description">1st order uniaxial anisotropy constant</div>

			<div class="header">Ku2:</div>
			<input placeholder=" {$parametersState.Ku2}" />
			<div>J/m<sup>3</sup></div>
			<div class="description">2nd order uniaxial anisotropy constant</div>

			<div class="header">anisC1:</div>
			<input placeholder=" {$parametersState.anisC1}" />
			<div>Vector(x,y,z)</div>
			<div class="description">Cubic anisotropy direction #1</div>

			<div class="header">anisC2:</div>
			<input placeholder=" {$parametersState.anisC2}" />
			<div>Vector(x,y,z)</div>
			<div class="description">Cubic anisotropy direction #2</div>

			<div class="header">Kc1:</div>
			<input placeholder=" {$parametersState.Kc1}" />
			<div>J/m<sup>3</sup></div>
			<div class="description">1st order cubic anisotropy constant</div>

			<div class="header">Kc2:</div>
			<input placeholder=" {$parametersState.Kc2}" />
			<div>J/m<sup>3</sup></div>
			<div class="description">2nd order cubic anisotropy constant</div>

			<div class="header">Kc3:</div>
			<input placeholder=" {$parametersState.Kc3}" />
			<div>J/m<sup>3</sup></div>
			<div class="description">3rd order cubic anisotropy constant</div>

			<div class="header">Dind:</div>
			<input placeholder=" {$parametersState.Dind}" />
			<div>J/m<sup>2</sup></div>
			<div class="description">Interfacial Dzyaloshinskii-Moriya strength</div>

			<div class="header">Dbulk:</div>
			<input placeholder=" {$parametersState.Dbulk}" />
			<div>J/m<sup>2</sup></div>
			<div class="description">Bulk Dzyaloshinskii-Moriya strength</div>

			<div class="header">B1:</div>
			<input placeholder=" {$parametersState.B1}" />
			<div>J/m<sup>3</sup></div>
			<div class="description">First magneto-elastic coupling constant</div>

			<div class="header">B2:</div>
			<input placeholder=" {$parametersState.B2}" />
			<div>J/m<sup>3</sup></div>
			<div class="description">Second magneto-elastic coupling constant</div>

			<div class="header">EpsilonPrime:</div>
			<input placeholder=" {$parametersState.EpsilonPrime}" />
			<div></div>
			<div class="description">Slonczewski secondary STT term ε'</div>

			<div class="header">exx:</div>
			<input placeholder=" {$parametersState.exx}" />
			<div></div>
			<div class="description">exx component of the strain tensor</div>

			<div class="header">exy:</div>
			<input placeholder=" {$parametersState.exy}" />
			<div></div>
			<div class="description">exy component of the strain tensor</div>

			<div class="header">exz:</div>
			<input placeholder=" {$parametersState.exz}" />
			<div></div>
			<div class="description">exz component of the strain tensor</div>

			<div class="header">eyy:</div>
			<input placeholder=" {$parametersState.eyy}" />
			<div></div>
			<div class="description">eyy component of the strain tensor</div>

			<div class="header">eyz:</div>
			<input placeholder=" {$parametersState.eyz}" />
			<div></div>
			<div class="description">eyz component of the strain tensor</div>

			<div class="header">ezz:</div>
			<input placeholder=" {$parametersState.ezz}" />
			<div></div>
			<div class="description">ezz component of the strain tensor</div>

			<div class="header">FixedLayer:</div>
			<input placeholder=" {$parametersState.FixedLayer}" />
			<div></div>
			<div class="description">Slonczewski fixed layer polarization</div>

			<div class="header">FreeLayerThickness:</div>
			<input placeholder=" {$parametersState.FreeLayerThickness}" />
			<div>m</div>
			<div class="description">Slonczewski free layer thickness</div>

			<div class="header">frozenspins:</div>
			<input placeholder=" {$parametersState.frozenspins}" />
			<div></div>
			<div class="description">Defines spins that should be fixed</div>

			<div class="header">J:</div>
			<input placeholder=" {$parametersState.J}" />
			<div>A/m<sup>2</sup></div>
			<div class="description">Electrical current density</div>

			<div class="header">Lambda:</div>
			<input placeholder=" {$parametersState.Lambda}" />
			<div></div>
			<div class="description">Slonczewski Λ parameter</div>

			<div class="header">MFMDipole:</div>
			<input placeholder=" {$parametersState.MFMDipole}" />
			<div>m</div>
			<div class="description">Height of vertically magnetized part of MFM tip</div>

			<div class="header">MFMLift:</div>
			<input placeholder=" {$parametersState.MFMLift}" />
			<div>m</div>
			<div class="description">MFM lift height</div>

			<div class="header">NoDemagSpins:</div>
			<input placeholder=" {$parametersState.NoDemagSpins}" />
			<div>Boolean</div>
			<div class="description">Disable magnetostatic interaction per region</div>

			<div class="header">Pol:</div>
			<input placeholder=" {$parametersState.Pol}" />
			<div></div>
			<div class="description">Electrical current polarization</div>

			<div class="header">Temp:</div>
			<input placeholder=" {$parametersState.Temp}" />
			<div>K</div>
			<div class="description">Temperature</div>

			<div class="header">xi:</div>
			<input placeholder=" {$parametersState.xi}" />
			<div></div>
			<div class="description">Non-adiabaticity of spin-transfer-torque</div>
		</div>
	</div>
</section>

<style>
	section {
		grid-area: parameters;
	}
	.grid-container {
		display: grid;
		grid-template-columns: auto auto auto auto;
		gap: 10px;
		align-items: center;
	}
	input {
		width: 100%;
		box-sizing: border-box;
	}
	.header {
		text-align: right;
		padding-right: 10px;
	}
	.description {
		font-style: italic;
		color: bisque;
		overflow: hidden;
	}
</style>
