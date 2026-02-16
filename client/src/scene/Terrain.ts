import * as THREE from "three";

/** Voxel-style terrain. */
export function createTerrain(width: number, height: number): THREE.Group {
  const group = new THREE.Group();
  
  // Use a single geometry for all ground blocks
  const geometry = new THREE.BoxGeometry(1, 1, 1);
  geometry.translate(0, -0.5, 0); // Pivot at top face

  // Use MeshBasicMaterial for no lighting
  // A nice bright green for grass
  const material = new THREE.MeshBasicMaterial({ color: 0x7cfc00 });

  const mesh = new THREE.InstancedMesh(geometry, material, width * height);
  
  const dummy = new THREE.Object3D();
  let i = 0;
  for (let x = 0; x < width; x++) {
    for (let z = 0; z < height; z++) {
      dummy.position.set(x + 0.5, 0, z + 0.5);
      dummy.updateMatrix();
      mesh.setMatrixAt(i++, dummy.matrix);
    }
  }
  
  group.add(mesh);

  // Custom grid using LineSegments
  const vertices: number[] = [];
  
  // Lines along X (horizontal)
  for (let z = 0; z <= height; z++) {
    vertices.push(0, 0.02, z);
    vertices.push(width, 0.02, z);
  }
  // Lines along Z (vertical)
  for (let x = 0; x <= width; x++) {
    vertices.push(x, 0.02, 0);
    vertices.push(x, 0.02, height);
  }
  
  const gridGeo = new THREE.BufferGeometry();
  gridGeo.setAttribute('position', new THREE.Float32BufferAttribute(vertices, 3));
  const gridMat = new THREE.LineBasicMaterial({ color: 0x333333, opacity: 0.3, transparent: true });
  const gridLines = new THREE.LineSegments(gridGeo, gridMat);
  group.add(gridLines);

  return group;
}
