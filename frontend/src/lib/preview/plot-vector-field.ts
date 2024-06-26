import { previewState } from '$api/incoming/preview';
import * as THREE from 'three';
import * as BufferGeometryUtils from 'three/examples/jsm/utils/BufferGeometryUtils.js';
import { TrackballControls } from 'three/examples/jsm/controls/TrackballControls.js';
import { get } from 'svelte/store';
import { writable } from 'svelte/store';

type Dimensions = [number, number, number];
interface Display {
    mesh: THREE.InstancedMesh;
    scene: THREE.Scene;
    camera: THREE.PerspectiveCamera;
    renderer: THREE.WebGLRenderer;
    controls: TrackballControls;
    isInitialized: boolean;
    parsingTime: number;
    dimensions: Dimensions;
}

export const display = writable<Display | null>(null);

export function plotVectorField() {
    if (get(display) === null) {
        init();
    } else {
        update();
    }
}

function createMesh(dimensions: Dimensions): THREE.InstancedMesh {
    const shaftGeometry = new THREE.CylinderGeometry(0.05, 0.05, 0.55, 8);
    shaftGeometry.translate(0, -0.1, 0); // Center the shaft
    const headGeometry = new THREE.ConeGeometry(0.2, 0.4, 8);
    headGeometry.translate(0, 0.4, 0); // Position the head at the top of the centered shaft
    let arrowGeometry = BufferGeometryUtils.mergeGeometries([shaftGeometry, headGeometry]);

    const arrowMaterial = new THREE.MeshBasicMaterial({ color: 0x0000ff });

    const totalVectors = dimensions[0] * dimensions[1] * dimensions[2];
    return new THREE.InstancedMesh(arrowGeometry, arrowMaterial, totalVectors);
}

function createCamera(dimensions: Dimensions): THREE.PerspectiveCamera {
    const camera = new THREE.PerspectiveCamera(50, window.innerWidth / window.innerHeight, 1, 100);
    camera.position.set(dimensions[0] / 2, dimensions[1] / 2, 15);
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

function createControls(dimensions: Dimensions, camera: THREE.PerspectiveCamera, renderer: THREE.WebGLRenderer) {
    const controls = new TrackballControls(camera, renderer.domElement);
    controls.dynamicDampingFactor = 1;
    controls.panSpeed = 0.8;
    controls.rotateSpeed = 1;
    controls.target.set(dimensions[0] / 2, dimensions[1] / 2, 0); // Set the center of the grid as the target
    controls.update();
    return controls;
}

function createScene() {
    const scene = new THREE.Scene();
    scene.background = new THREE.Color(0x282a36);
    return scene;
}

function addArrowsToMesh(mesh: THREE.InstancedMesh, dimensions: Dimensions) {
    const vectorField = getVectorField();
    const dummy = new THREE.Object3D();
    vectorField.forEach((_, i) => {
        const posx = i % dimensions[0];
        const posy = Math.floor(i / dimensions[0]) % dimensions[1];
        const posz = Math.floor(i / (dimensions[0] * dimensions[1]));
        dummy.position.set(posx, posy, posz);
        dummy.updateMatrix();
        mesh.setMatrixAt(i, dummy.matrix);
    });

    mesh.instanceMatrix.needsUpdate = true; // Notify Three.js that instance matrices need to be updated
}

export function init() {
    let dimensions: Dimensions = [17, 17, 1];
    let scene = createScene();
    let camera = createCamera(dimensions);
    let renderer = createRenderer();
    let controls = createControls(dimensions, camera, renderer);
    let mesh = createMesh(dimensions);
    addArrowsToMesh(mesh, dimensions);
    scene.add(mesh);
    update();
    display.set({
        mesh,
        scene,
        camera,
        renderer,
        controls,
        isInitialized: true,
        parsingTime: 0,
        dimensions
    });

    function animate() {
        requestAnimationFrame(animate);
        controls.update();
        renderer.render(scene, camera);
    }
    animate();
}

function getVectorField(): Array<{ x: number; y: number; z: number }> {
    const buffer = get(previewState).buffer;
    console.log("getVectorField:" + buffer.byteLength);
    const float32Array = new Float32Array(buffer);
    const vectors: Array<{ x: number; y: number; z: number }> = [];

    for (let i = 0; i < float32Array.length; i += 3) {
        vectors.push({
            x: float32Array[i],
            y: float32Array[i + 1],
            z: float32Array[i + 2]
        });
    }
    console.log(vectors)
    return vectors;
}
export function update() {
    let d = get(display);
    if (d) {
        const dummy = new THREE.Object3D();
        const defaultVector = new THREE.Vector3(0, 1, 0); // Default orientation of the arrow
        const mesh = d.mesh;
        let vectorField = getVectorField();

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
        display.update(currentDisplay => {
            if (currentDisplay) {
                return { ...currentDisplay, mesh: mesh };
            }
            return currentDisplay;
        }
        );
    };
}

export function resetCamera(camera: THREE.PerspectiveCamera, controls: TrackballControls) {
    camera.position.set(8, 8, 17);
    controls.target.set(8, 8, 0);
}

