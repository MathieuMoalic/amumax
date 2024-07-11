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
    let disp = get(display);
    if (disp === null || !disp.isInitialized) {
        init();
    }
    update();
}

function createMesh(): THREE.InstancedMesh {
    let dims = get(previewState).dimensions;
    const shaftGeometry = new THREE.CylinderGeometry(0.05, 0.05, 0.55, 8);
    shaftGeometry.translate(0, -0.06, 0); // Center the shaft
    const headGeometry = new THREE.ConeGeometry(0.2, 0.4, 8);
    headGeometry.translate(0, 0.4, 0); // Position the head at the top of the centered shaft
    let arrowGeometry = BufferGeometryUtils.mergeGeometries([shaftGeometry, headGeometry]);

    const arrowMaterial = new THREE.MeshBasicMaterial({ vertexColors: false });

    const totalVectors = dims[0] * dims[1] * dims[2];
    const mesh = new THREE.InstancedMesh(arrowGeometry, arrowMaterial, totalVectors);
    return mesh;
}

function createCamera(): THREE.PerspectiveCamera {
    let fov = 50;
    let div = document.getElementById('container');
    let width = div?.offsetWidth || 1;
    let height = div?.offsetHeight || 1;
    let aspect = width / height;
    let near = 0.1;
    let far = 1000;
    const camera = new THREE.PerspectiveCamera(fov, aspect, near, far);

    let dims = get(previewState).dimensions;
    let posx = dims[0] / 2;
    let posy = dims[1] / 2;
    let posz = Math.max(dims[0], dims[1]) * 1.5;
    camera.position.set(posx, posy, posz);
    return camera;
}

function createRenderer() {
    const renderer = new THREE.WebGLRenderer({ antialias: true });
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
    let dims = get(previewState).dimensions;
    const controls = new TrackballControls(camera, renderer.domElement);
    controls.dynamicDampingFactor = 1;
    controls.panSpeed = 0.8;
    controls.rotateSpeed = 1;
    controls.target.set(dims[0] / 2, dims[1] / 2, 0); // Set the center of the grid as the target
    controls.update();
    return controls;
}

function createScene() {
    const scene = new THREE.Scene();
    scene.background = new THREE.Color(0x282a36);
    return scene;
}

function addArrowsToMesh(mesh: THREE.InstancedMesh, scene: THREE.Scene) {
    const vectorField = getVectorField();
    const dummy = new THREE.Object3D();
    let dims = get(previewState).dimensions;

    // if (!mesh.instanceColor) {
    //     mesh.instanceColor = new THREE.InstancedBufferAttribute(colors, 3);
    // }
    const totalVectors = dims[0] * dims[1] * dims[2];
    let instanceColor = new THREE.InstancedBufferAttribute(new Float32Array(totalVectors * 3), 3)
    const colors = instanceColor.array;
    const color = new THREE.Color();


    vectorField.forEach((vector, i) => {
        if (vector.x === 0 && vector.y === 0 && vector.z === 0) {
            return;
        }
        const posx = i % dims[0];
        const posy = Math.floor(i / dims[0]) % dims[1];
        const posz = Math.floor(i / (dims[0] * dims[1]));
        dummy.position.set(posx, posy, posz);

        let h = Math.atan2(vector.y, vector.x) / Math.PI / 2;
        let s = Math.sqrt(vector.x ** 2 + vector.y ** 2 + vector.z ** 2);
        let l = (vector.z + 1) / 2;
        color.setHSL(h, s, l); // Map z to HSL for a smooth gradient
        colors[i * 3 + 0] = color.r;
        colors[i * 3 + 1] = color.g;
        colors[i * 3 + 2] = color.b;

        dummy.updateMatrix();
        mesh.matrix.makeScale(1, 1, 1)
        mesh.setMatrixAt(i, dummy.matrix);
        // TODO: REMOVE DUMMY
    });

    mesh.instanceMatrix.needsUpdate = true;
    instanceColor.needsUpdate = true; // Notify Three.js that colors need to be updated
    mesh.instanceColor = instanceColor;
}

function init() {
    let scene = createScene();
    let camera = createCamera();
    let renderer = createRenderer();
    let controls = createControls(camera, renderer);
    let mesh = createMesh();
    addArrowsToMesh(mesh, scene);
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
        dimensions: get(previewState).dimensions
    });

    function animate() {
        requestAnimationFrame(animate);
        controls.update();
        renderer.render(scene, camera);
    }
    animate();
}

function getVectorField(): Array<{ x: number; y: number; z: number }> {
    // Convert Uint8Array to ArrayBuffer
    const buffer = get(previewState).buffer.buffer;

    // Ensure the buffer is a multiple of 4 bytes since each float32 is 4 bytes.
    if (buffer.byteLength % 4 !== 0) {
        throw new Error("Invalid buffer length. Buffer length must be a multiple of 4 bytes.");
    }

    // Create Float32Array from the buffer
    const float32Array = new Float32Array(buffer);

    // Ensure the array length is a multiple of 3 (since each vector has 3 components)
    if (float32Array.length % 3 !== 0) {
        throw new Error("Invalid Float32Array length. Length must be a multiple of 3.");
    }

    const vectors: Array<{ x: number; y: number; z: number }> = [];

    for (let i = 0; i < float32Array.length; i += 3) {
        vectors.push({
            x: float32Array[i],
            y: float32Array[i + 1],
            z: float32Array[i + 2]
        });
    }

    return vectors;
}

function update() {
    let d = get(display);
    if (d) {
        const dummy = new THREE.Object3D();
        const defaultVector = new THREE.Vector3(0, 1, 0); // Default orientation of the arrow
        const mesh = d.mesh;
        let vectorField = getVectorField();


        let instanceColor = mesh.instanceColor;
        if (!instanceColor) {
            throw new Error('instanceColor not found');
        }
        const colors = instanceColor.array;
        const color = new THREE.Color();

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

            let h = Math.atan2(vector.y, vector.x) / Math.PI / 2;
            let s = Math.sqrt(vector.x ** 2 + vector.y ** 2 + vector.z ** 2);
            let l = (vector.z + 1) / 2;
            color.setHSL(h, s, l); // Map z to HSL for a smooth gradient
            colors[i * 3 + 0] = color.r;
            colors[i * 3 + 1] = color.g;
            colors[i * 3 + 2] = color.b;

            dummy.updateMatrix();
            mesh.setMatrixAt(i, dummy.matrix);
        });

        mesh.instanceMatrix.needsUpdate = true;
        instanceColor.needsUpdate = true; // Notify Three.js that colors need to be updated

        display.update(currentDisplay => {
            if (currentDisplay) {
                return { ...currentDisplay, mesh: mesh };
            }
            return currentDisplay;
        });
    }
}

export function resetCamera() {
    let dims = get(previewState).dimensions;
    let posx = dims[0] / 2;
    let posy = dims[1] / 2;
    let posz = Math.max(dims[0], dims[1]) * 1.5;

    let displayInstance = get(display);
    if (displayInstance) {
        let camera = displayInstance.camera;
        let controls = displayInstance.controls;

        // Reset camera position
        camera.position.set(posx, posy, posz);

        // Reset camera rotation
        camera.up.set(0, 1, 0);  // Ensure the camera's up direction is the Y-axis
        camera.lookAt(dims[0] / 2, dims[1] / 2, 0);

        // Reset controls
        controls.target.set(dims[0] / 2, dims[1] / 2, 0);
        controls.update();
    }
}

export function disposeThreeJS() {
    const container = document.getElementById('container');
    const displayInstance = get(display);
    if (displayInstance && displayInstance.renderer) {
        displayInstance.renderer.dispose();
        displayInstance.scene.children.forEach(child => {
            if (child instanceof THREE.Mesh) {
                child.geometry.dispose();
                if (child.material instanceof THREE.Material) {
                    child.material.dispose();
                }
            }
        });
        display.set(null); // Reset the display store
        if (container) {
            container.innerHTML = '';
        }
    }
}