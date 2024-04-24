<script lang="ts">
	let renderQuant = 'm';
	let renderComp = 'All';
	let renderLayer = 0;
	let renderScale = 0;
	let quantities: string[] = [
		'Aex',
		'alpha',
		'anisC1',
		'anisC2',
		'anisU',
		'B1',
		'B2',
		'B_anis',
		'B_custom',
		'B_demag',
		'B_eff',
		'B_exch',
		'B_ext',
		'B_mel',
		'B_therm',
		'Dbulk',
		'Dind',
		'DindCoupling',
		'dt',
		'E_anis',
		'E_custom',
		'E_demag',
		'E_exch',
		'E_mel',
		'E_therm',
		'E_total',
		'E_Zeeman',
		'Edens_anis',
		'Edens_custom',
		'Edens_demag',
		'Edens_exch',
		'Edens_mel',
		'Edens_therm',
		'Edens_total',
		'Edens_Zeeman',
		'EpsilonPrime',
		'ExchCoupling',
		'ext_bubbledist',
		'ext_bubblepos',
		'ext_bubblespeed',
		'ext_corepos',
		'ext_dwpos',
		'ext_dwspeed',
		'ext_dwtilt',
		'ext_dwxpos',
		'ext_topologicalcharge',
		'ext_topologicalchargedensity',
		'ext_topologicalchargedensitylattice',
		'ext_topologicalchargelattice',
		'exx',
		'exy',
		'exz',
		'eyy',
		'eyz',
		'ezz',
		'F_mel',
		'FixedLayer',
		'FreeLayerThickness',
		'frozenspins',
		'geom',
		'J',
		'Kc1',
		'Kc2',
		'Kc3',
		'Ku1',
		'Ku2',
		'Lambda',
		'LastErr',
		'LLtorque',
		'm',
		'm_full',
		'MaxAngle',
		'maxTorque',
		'MFM',
		'Msat',
		'NEval',
		'NoDemagSpins',
		'PeakErr',
		'Pol',
		'regions',
		'spinAngle',
		'STTorque',
		'Temp',
		'torque',
		'xi'
	];
</script>

<section>
	<span> ▾ display </span>
	<br />
	<div>
		<p>
			Quantity:
			<select bind:value={renderQuant}>
				{#each quantities as q}
					<option value={q}>{q}</option>
				{/each}
			</select>
			<select bind:value={renderComp}>
				{#each ['All', 'x', 'y', 'z'] as q}
					<option value={q}>{q}</option>
				{/each}
			</select>
			<br />
		</p>
		Slice:<input type="range" min="0" max="0" bind:value={renderLayer} />
		<span>{renderLayer}</span>
		Scale: <input type="range" min="0" max="31" bind:value={renderScale} />
		<span>{renderScale}</span>
		<img alt="display" src="/render/m/?13149_0" />
	</div>
</section>
