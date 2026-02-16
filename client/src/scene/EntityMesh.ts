import * as THREE from "three";
import type { InterpolatedEntity } from "./Interpolator.js";

const ENTITY_WIDTH = 0.5;
const ENTITY_HEIGHT = 0.6;
const ENTITY_DEPTH = 0.5;

function entityColor(type: string): number {
  switch (type) {
    case "goat":
      return 0xddccaa;
    case "wolf":
      return 0x888888;
    default:
      return 0xaaaaaa;
  }
}

/** One cuboid per entity; position and rotation updated each frame from interpolated state. */
export function createEntityMesh(entity: InterpolatedEntity): THREE.Group {
  const group = new THREE.Group();
  const baseColor = entityColor(entity.type);

  const geometry = new THREE.BoxGeometry(
    ENTITY_WIDTH,
    ENTITY_HEIGHT,
    ENTITY_DEPTH
  );
  const material = new THREE.MeshBasicMaterial({
    color: baseColor,
  });
  const mesh = new THREE.Mesh(geometry, material);
  mesh.name = "entity-body";
  group.add(mesh);

  // Add edges
  const edges = new THREE.EdgesGeometry(geometry);
  const line = new THREE.LineSegments(
    edges,
    new THREE.LineBasicMaterial({ color: 0x000000 })
  );
  line.name = "entity-outline";
  group.add(line);

  group.userData = { entityId: entity.id, baseColor };
  updateEntityMesh(group, entity);
  return group;
}

export function updateEntityMesh(
  mesh: THREE.Object3D, // Changed to Object3D because it's a Group now
  entity: InterpolatedEntity
): void {
  mesh.position.set(entity.position[0] + 0.5, ENTITY_HEIGHT / 2, entity.position[1] + 0.5);
  const [dx, dy] = entity.direction;
  const angle = Math.atan2(dx, dy);
  mesh.rotation.y = angle;
}

export function setEntityHighlight(mesh: THREE.Object3D, selected: boolean): void {
  const group = mesh as THREE.Group;
  const baseColor = group.userData.baseColor as number;
  const body = group.getObjectByName("entity-body");
  const outline = group.getObjectByName("entity-outline");

  if (body instanceof THREE.Mesh && body.material instanceof THREE.MeshBasicMaterial) {
    body.material.color.set(selected ? 0xffd54f : baseColor);
  }

  if (outline instanceof THREE.LineSegments && outline.material instanceof THREE.LineBasicMaterial) {
    outline.material.color.set(selected ? 0xff8f00 : 0x000000);
  }

  group.scale.setScalar(selected ? 1.25 : 1);
}
