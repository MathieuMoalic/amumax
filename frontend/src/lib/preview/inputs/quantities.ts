export const quantities: { [category: string]: string[] } = {
	Common: ['m', 'torque', 'regions', 'Msat', 'Aex', 'alpha'],

	'Magnetic Fields': [
		'B_anis',
		'B_custom',
		'B_demag',
		'B_eff',
		'B_exch',
		'B_ext',
		'B_mel',
		'B_therm'
	],

	Energy: [
		'Edens_anis',
		'Edens_custom',
		'Edens_demag',
		'Edens_exch',
		'Edens_mel',
		'Edens_therm',
		'Edens_total',
		'Edens_Zeeman'
	],

	Anisotropy: ['anisC1', 'anisC2', 'anisU', 'Kc1', 'Kc2', 'Kc3', 'Ku1', 'Ku2'],

	DMI: ['Dbulk', 'Dind', 'DindCoupling'],

	External: [
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
		'ext_topologicalchargelattice'
	],

	'Spin-transfer Torque': ['xi', 'STTorque'],

	Strain: ['exx', 'exy', 'exz', 'eyy', 'eyz', 'ezz'],

	Current: ['J', 'Pol'],

	Slonczewski: ['EpsilonPrime', 'FixedLayer', 'FreeLayerThickness', 'Lambda'],

	'Magneto-elastic': ['B1', 'B2', 'F_mel', 'B_mel'],

	Miscellaneous: [
		'frozenspins',
		'NoDemagSpins',
		'MFM',
		'spinAngle',
		'LLtorque',
		'm_full',
		'Temp',
		'geom'
	]
};
