import * as THREE from "three";
import { OrbitControls } from "three/examples/jsm/controls/OrbitControls.js";
import type { WorldState } from "../models/world.js";
import { Interpolator } from "./Interpolator.js";
import { createTerrain } from "./Terrain.js";
import { createObstacleMeshes } from "./ObstacleMesh.js";
import {
  createEntityMesh,
  updateEntityMesh,
} from "./EntityMesh.js";
import type { InterpolatedEntity } from "./Interpolator.js";

export class Scene {
  private scene: THREE.Scene;
  private camera: THREE.OrthographicCamera;
  private renderer: THREE.WebGLRenderer;
  private controls: OrbitControls;
  private interpolator = new Interpolator();
  private terrainGroup: THREE.Group | null = null;
  private obstacleGroup: THREE.Group | null = null;
  private entityMeshes = new Map<number, THREE.Object3D>();
  private entityGroup = new THREE.Group();

  // Frustum size for orthographic camera (view size in world units)
  // Increased to cover 30x30 world comfortably
  private frustumSize = 50;

  constructor(canvas: HTMLCanvasElement) {
    this.scene = new THREE.Scene();
    // Darker background for "blueprint" or "schematic" feel, or just clean white?
    // User asked for "impersive of low poly", often bright colors.
    // Let's use a soft gray-blue.
    this.scene.background = new THREE.Color(0xd0e0e3);

    const width = canvas.clientWidth;
    const height = canvas.clientHeight;
    const aspect = width / height;

    // Orthographic camera for isometric view
    this.camera = new THREE.OrthographicCamera(
      (this.frustumSize * aspect) / -2,
      (this.frustumSize * aspect) / 2,
      this.frustumSize / 2,
      this.frustumSize / -2,
      1,
      1000
    );

    // Isometric position
    this.camera.position.set(100, 100, 100);
    this.camera.lookAt(0, 0, 0);
    this.camera.up.set(0, 1, 0);

    this.renderer = new THREE.WebGLRenderer({ canvas, antialias: true });
    this.renderer.setSize(width, height);
    this.renderer.setPixelRatio(Math.min(window.devicePixelRatio, 2));
    // No shadows needed as per "no lighting calculation"

    // Controls
    this.controls = new OrbitControls(this.camera, canvas);
    this.controls.enableDamping = true;
    this.controls.dampingFactor = 0.05;
    this.controls.screenSpacePanning = true;
    this.controls.minZoom = 0.5;
    this.controls.maxZoom = 4;

    this.scene.add(this.entityGroup);
  }

  updateWorldState(state: WorldState): void {
    this.interpolator.update(state);

    const w = state.width;
    const h = state.height;

    // Dynamically adjust frustum to fit the world if needed
    // But fixed 50 is usually good for 30x30 with some padding.

    if (!this.terrainGroup) {
      this.terrainGroup = createTerrain(w, h);
      this.scene.add(this.terrainGroup);
      
      // Center controls target on the world center
      this.controls.target.set(w / 2, 0, h / 2);
      this.camera.position.set(w / 2 + 40, 40, h / 2 + 40);
      this.camera.lookAt(w / 2, 0, h / 2);
      this.controls.update();
    }

    if (!this.obstacleGroup) {
      this.obstacleGroup = createObstacleMeshes(state.static_obstacles);
      this.scene.add(this.obstacleGroup);
    }
  }

  private syncEntities(interpolated: InterpolatedEntity[]): void {
    const ids = new Set(interpolated.map((e) => e.id));
    for (const [id, mesh] of this.entityMeshes) {
      if (!ids.has(id)) {
        this.entityGroup.remove(mesh);
        // mesh is now a Group, so we need to dispose children
        mesh.traverse((child) => {
          if (child instanceof THREE.Mesh) {
            child.geometry.dispose();
            (child.material as THREE.Material).dispose();
          }
        });
        this.entityMeshes.delete(id);
      }
    }
    for (const ent of interpolated) {
      let mesh = this.entityMeshes.get(ent.id);
      if (!mesh) {
        mesh = createEntityMesh(ent);
        this.entityMeshes.set(ent.id, mesh);
        this.entityGroup.add(mesh);
      } else {
        updateEntityMesh(mesh, ent);
      }
    }
  }

  render(): void {
    this.controls.update();
    const now = performance.now();
    const entities = this.interpolator.getInterpolatedEntities(now);
    this.syncEntities(entities);

    this.renderer.render(this.scene, this.camera);
  }

  resize(width: number, height: number): void {
    const aspect = width / height;
    this.camera.left = (-this.frustumSize * aspect) / 2;
    this.camera.right = (this.frustumSize * aspect) / 2;
    this.camera.top = this.frustumSize / 2;
    this.camera.bottom = -this.frustumSize / 2;
    this.camera.updateProjectionMatrix();
    this.renderer.setSize(width, height);
  }
}
