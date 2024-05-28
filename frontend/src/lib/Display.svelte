<script lang="ts">
	import quantities from './quantities';
	import { meshState } from './api';
	import {
		refreshImage,
		renderLayer,
		imageComponent,
		imageQuantity,
		fetchVectorField
	} from './api';
	import * as THREE from 'three';
	import * as BufferGeometryUtils from 'three/examples/jsm/utils/BufferGeometryUtils.js';
	import { TrackballControls } from 'three/examples/jsm/controls/TrackballControls.js';
	import { onMount } from 'svelte';

	let dims: [number, number, number] = [17, 17, 1];
	function createArrowMesh() {
		// Adjust shaft geometry so its center aligns with the desired position
		const shaftGeometry = new THREE.CylinderGeometry(0.05, 0.05, 0.8, 8);
		shaftGeometry.translate(0, -0.1, 0); // Center the shaft

		// Adjust head geometry to align with the shaft, then adjust both to be centered
		const headGeometry = new THREE.ConeGeometry(0.2, 0.4, 8);
		headGeometry.translate(0, 0.4, 0); // Position the head at the top of the centered shaft

		// Merge geometries
		let arrowGeometry = BufferGeometryUtils.mergeGeometries([shaftGeometry, headGeometry]);

		const arrowMaterial = new THREE.MeshBasicMaterial({ color: 0x0000ff });

		const totalVectors = dims[0] * dims[1] * dims[2];
		const arrowMesh = new THREE.InstancedMesh(arrowGeometry, arrowMaterial, totalVectors);
		return arrowMesh;
	}

	async function addVectorsToMesh(arrowMesh: THREE.InstancedMesh) {
		const vectorField = await fetchVectorField();
		const dummy = new THREE.Object3D();
		const defaultVector = new THREE.Vector3(0, 1, 0); // This is the default orientation of the arrow
		vectorField.forEach((vector, i) => {
			if (vector.x === 0 && vector.y === 0 && vector.z === 0) {
				return;
			}
			const posx = i % dims[0];
			const posy = Math.floor(i / dims[0]) % dims[1];
			const posz = Math.floor(i / (dims[0] * dims[1]));
			dummy.position.set(posx, posy, posz);
			dummy.quaternion.setFromUnitVectors(
				defaultVector,
				new THREE.Vector3(vector.x, vector.y, vector.z).normalize()
			);
			dummy.updateMatrix();
			arrowMesh.setMatrixAt(i, dummy.matrix);
		});
	}

	function createCamera() {
		const camera = new THREE.PerspectiveCamera(50, window.innerWidth / window.innerHeight, 1, 100);
		camera.position.set(8, 8, 17); // Set the camera position
		return camera;
	}
	function createRenderer() {
		const renderer = new THREE.WebGLRenderer();
		const container = document.getElementById('container');
		if (!container) {
			throw new Error('Container not found');
		}
		renderer.setSize(container.clientWidth, container.clientHeight);
		if (container) {
			container.appendChild(renderer.domElement);
		}
		return renderer;
	}
	function createControls(camera: THREE.PerspectiveCamera, renderer: THREE.WebGLRenderer) {
		// const controls = new OrbitControls(camera, renderer.domElement);
		const controls = new TrackballControls(camera, renderer.domElement);
		controls.dynamicDampingFactor = 1;
		controls.panSpeed = 0.8;
		controls.rotateSpeed = 1;
		controls.target.set(8, 8, 0); // Set the center of the grid as the target
		controls.update();
		return controls;
	}

	onMount(async () => {
		let arrowMesh = createArrowMesh();
		await addVectorsToMesh(arrowMesh);
		const scene = new THREE.Scene();
		scene.add(arrowMesh);

		let camera = createCamera();
		let renderer = createRenderer();
		let controls = createControls(camera, renderer);

		function animate() {
			requestAnimationFrame(animate);
			controls.update(); // Update controls in the animation loop
			renderer.render(scene, camera);
		}

		animate();
	});
</script>

<section>
	<h2>Preview</h2>
	<button on:click={fetchVectorField}>Fetch and print vector field</button>
	<p>
		Quantity:
		<select bind:value={$imageQuantity} on:change={refreshImage}>
			{#each quantities as quantity}
				<option value={quantity}>{quantity}</option>
			{/each}
		</select>
		<select bind:value={$imageComponent}>
			{#each ['All', 'x', 'y', 'z'] as choice}
				<option value={choice}>{choice}</option>
			{/each}
		</select>
		Z-layer: 0
		<input type="range" min="0" max={$meshState.Nz - 1} bind:value={$renderLayer} />
		{$renderLayer}
	</p>
	<div id="container"></div>
</section>

<style>
	section {
		grid-area: display;
	}
	#container {
		width: 500px;
		height: 500px;
	}
</style>
