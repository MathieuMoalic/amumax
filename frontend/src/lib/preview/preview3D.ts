import { previewState } from '$api/incoming/preview';
import * as THREE from 'three';
import * as BufferGeometryUtils from 'three/examples/jsm/utils/BufferGeometryUtils.js';
import { TrackballControls } from 'three/examples/jsm/controls/TrackballControls.js';
import { get } from 'svelte/store';
import { writable } from 'svelte/store';
import { disposePreview2D } from './preview2D';


export function preview3D() {
    if (get(previewState).refresh) {
        disposePreview2D();
        disposePreview3D();
        init();
        update();
    } else {
        update();
    }
}

export function disposePreview3D() {
    const container = document.getElementById('container');
    const displayInstance = get(threeDPreview);
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
        threeDPreview.set(null); // Reset the display store
        if (container) {
            container.innerHTML = '';
        }
    }
}

type Dimensions = [number, number, number];
interface ThreeDPreview {
    mesh: THREE.InstancedMesh;
    scene: THREE.Scene;
    camera: THREE.PerspectiveCamera;
    renderer: THREE.WebGLRenderer;
    controls: TrackballControls;
    isInitialized: boolean;
    parsingTime: number;
    dimensions: Dimensions;
}
export const threeDPreview = writable<ThreeDPreview | null>(null);

export function plotVectorField(newDimensions: boolean) {
    let disp = get(threeDPreview);
    if (disp === null || !disp.isInitialized) {
        init();
    } else if (newDimensions) {
        disposePreview3D();
        init();
    }
    update();
}

function createMesh(): THREE.InstancedMesh {
    const shaftGeometry = new THREE.CylinderGeometry(0.05, 0.05, 0.55, 8);
    shaftGeometry.translate(0, -0.06, 0); // Center the shaft
    const headGeometry = new THREE.ConeGeometry(0.2, 0.4, 8);
    headGeometry.translate(0, 0.4, 0); // Position the head at the top of the centered shaft
    let arrowGeometry = BufferGeometryUtils.mergeGeometries([shaftGeometry, headGeometry]);

    const arrowMaterial = new THREE.MeshBasicMaterial({ vertexColors: false });

    const mesh = new THREE.InstancedMesh(arrowGeometry, arrowMaterial, get(previewState).vectorFieldValues.length);
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
    const vectorFieldValues = get(previewState).vectorFieldValues;
    const vectorFieldPositions = get(previewState).vectorFieldPositions;

    const dummy = new THREE.Object3D();

    // Times 3 because each color has 3 components (RGB)
    const instanceColorLength = vectorFieldPositions.length * 3;
    let instanceColor = new THREE.InstancedBufferAttribute(new Float32Array(instanceColorLength), 3)
    const colors = instanceColor.array;
    const color = new THREE.Color();


    for (let i = 0; i < vectorFieldValues.length; i++) {
        const val = vectorFieldValues[i];
        const pos = vectorFieldPositions[i];
        dummy.position.set(pos.x, pos.y, pos.z);

        let h = Math.atan2(val.y, val.x) / Math.PI / 2;
        let s = Math.sqrt(val.x ** 2 + val.y ** 2 + val.z ** 2);
        let l = (val.z + 1) / 2;
        color.setHSL(h, s, l); // Map z to HSL for a smooth gradient
        colors[i * 3 + 0] = color.r;
        colors[i * 3 + 1] = color.g;
        colors[i * 3 + 2] = color.b;

        dummy.updateMatrix();
        mesh.setMatrixAt(i, dummy.matrix);
        // TODO: REMOVE DUMMY
    };
    mesh.matrix.makeScale(1, 1, 1)

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
    threeDPreview.set({
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


function update() {
    let d = get(threeDPreview);
    if (d) {
        const dummy = new THREE.Object3D();
        const defaultVector = new THREE.Vector3(0, 1, 0); // Default orientation of the arrow
        const mesh = d.mesh;
        let vectorField = get(previewState).vectorFieldValues;


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

        threeDPreview.update(currentDisplay => {
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

    let displayInstance = get(threeDPreview);
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

