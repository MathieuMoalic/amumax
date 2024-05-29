import {
    displayData
} from './api';
import * as THREE from 'three';
import * as BufferGeometryUtils from 'three/examples/jsm/utils/BufferGeometryUtils.js';
import { TrackballControls } from 'three/examples/jsm/controls/TrackballControls.js';
import { get } from 'svelte/store';
import { writable } from 'svelte/store';

export const displayMesh = writable<Mesh>(new THREE.InstancedMesh(new THREE.BoxGeometry(), new THREE.MeshBasicMaterial(), 0));
export const isDisplayInitialized = writable<boolean>(false);
export const parsingTime = writable<number>(0);


let dims: [number, number, number] = [17, 17, 1];

function createMesh() {
    // Adjust shaft geometry so its center aligns with the desired position
    const shaftGeometry = new THREE.CylinderGeometry(0.05, 0.05, 0.55, 8);
    shaftGeometry.translate(0, -0.1, 0); // Center the shaft


    // Adjust head geometry to align with the shaft, then adjust both to be centered
    const headGeometry = new THREE.ConeGeometry(0.2, 0.4, 8);
    headGeometry.translate(0, 0.4, 0); // Position the head at the top of the centered shaft

    // Merge geometries
    let arrowGeometry = BufferGeometryUtils.mergeGeometries([shaftGeometry, headGeometry]);

    const arrowMaterial = new THREE.MeshBasicMaterial({ color: 0x0000ff });

    const totalVectors = dims[0] * dims[1] * dims[2];
    displayMesh.set(new THREE.InstancedMesh(arrowGeometry, arrowMaterial, totalVectors))
}

type Mesh = THREE.InstancedMesh<THREE.BufferGeometry<THREE.NormalBufferAttributes>, THREE.Material | THREE.Material[], THREE.InstancedMeshEventMap>



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
function createScene() {
    const scene = new THREE.Scene();
    scene.background = new THREE.Color(0x282a36);
    return scene;
}
export function init() {
    const scene = createScene();
    const camera = createCamera();
    const renderer = createRenderer();
    const controls = createControls(camera, renderer);
    createMesh();
    addArrowsToMesh();
    update();
    scene.add(get(displayMesh));

    function animate() {
        requestAnimationFrame(animate);
        controls.update();
        renderer.render(scene, camera);
    }
    animate();
}

function addArrowsToMesh() {
    const vectorField = get(displayData);
    const dummy = new THREE.Object3D();

    const mesh = get(displayMesh);

    vectorField.forEach((_, i) => {
        const posx = i % dims[0];
        const posy = Math.floor(i / dims[0]) % dims[1];
        const posz = Math.floor(i / (dims[0] * dims[1]));

        // Set the position
        dummy.position.set(posx, posy, posz);
        dummy.updateMatrix();
        mesh.setMatrixAt(i, dummy.matrix);
    });

    mesh.instanceMatrix.needsUpdate = true; // Notify Three.js that instance matrices need to be updated
}
export function update() {
    const vectorField = get(displayData);

    const dummy = new THREE.Object3D();
    const defaultVector = new THREE.Vector3(0, 1, 0); // Default orientation of the arrow
    const mesh = get(displayMesh);

    vectorField.forEach((vector, i) => {
        if (vector.x === 0 && vector.y === 0 && vector.z === 0) {
            return;
        }

        mesh.getMatrixAt(i, dummy.matrix);
        dummy.matrix.decompose(dummy.position, dummy.quaternion, dummy.scale);
        dummy.quaternion.setFromUnitVectors(
            defaultVector,
            new THREE.Vector3(vector.x, vector.y, vector.z).normalize()
        );
        dummy.updateMatrix();
        mesh.setMatrixAt(i, dummy.matrix);
    });
    mesh.instanceMatrix.needsUpdate = true;
}

