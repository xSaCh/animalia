import * as THREE from "three";
import type { StaticObstacles } from "../models/world.js";

const CELL = 1;

function makeBox(w: number, h: number, d: number, color: number): THREE.Group {
  const group = new THREE.Group();
  
  const geometry = new THREE.BoxGeometry(w, h, d);
  const material = new THREE.MeshBasicMaterial({ color });
  const mesh = new THREE.Mesh(geometry, material);
  group.add(mesh);

  // Add edges for "blueprint" look
  const edges = new THREE.EdgesGeometry(geometry);
  const line = new THREE.LineSegments(
    edges,
    new THREE.LineBasicMaterial({ color: 0x000000, linewidth: 2 })
  );
  group.add(line);

  return group;
}

/** Walls: gray cuboids. Water: blue low-poly. Food: green/yellow low-poly. */
export function createObstacleMeshes(obstacles: StaticObstacles): THREE.Group {
  const group = new THREE.Group();

  for (const o of obstacles.walls) {
    const mesh = makeBox(CELL * 0.9, CELL * 0.5, CELL * 0.9, 0x555555);
    mesh.position.set(o.position.x + 0.5, 0.25, o.position.y + 0.5);
    group.add(mesh);
  }

  for (const o of obstacles.water_sources) {
    const mesh = makeBox(CELL * 0.85, CELL * 0.3, CELL * 0.85, 0x2288cc);
    mesh.position.set(o.position.x + 0.5, 0.15, o.position.y + 0.5);
    group.add(mesh);
  }

  for (const o of obstacles.food_sources) {
    const mesh = makeBox(CELL * 0.7, CELL * 0.4, CELL * 0.7, 0x88aa22);
    mesh.position.set(o.position.x + 0.5, 0.2, o.position.y + 0.5);
    group.add(mesh);
  }

  for (const o of obstacles.rest_areas) {
    const mesh = makeBox(CELL * 0.8, CELL * 0.2, CELL * 0.8, 0x666666);
    mesh.position.set(o.position.x + 0.5, 0.1, o.position.y + 0.5);
    group.add(mesh);
  }

  return group;
}
