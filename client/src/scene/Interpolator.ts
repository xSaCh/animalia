import type { WorldState, Entity } from "../models/world.js";

function lerp(a: number, b: number, t: number): number {
  return a + (b - a) * t;
}

function lerpVec(
  ax: number,
  ay: number,
  bx: number,
  by: number,
  t: number
): [number, number] {
  return [lerp(ax, bx, t), lerp(ay, by, t)];
}

export interface InterpolatedEntity {
  id: number;
  type: string;
  position: [number, number];
  direction: [number, number];
  state: string;
}

export class Interpolator {
  private prev: WorldState | null = null;
  private next: WorldState | null = null;
  private nextReceivedAt = 0;
  private snapshotIntervalMs = 200;

  update(state: WorldState): void {
    this.prev = this.next;
    this.next = state;
    this.nextReceivedAt = performance.now();
    if (state.config?.tps) {
      this.snapshotIntervalMs = 1000 / state.config.tps;
    }
  }

  getInterpolatedEntities(now: number): InterpolatedEntity[] {
    const next = this.next;
    if (!next) return [];

    const prev = this.prev;
    if (!prev) {
      return next.entities.map((e) => ({
        id: e.id,
        type: e.type,
        position: [e.position.x, e.position.y] as [number, number],
        direction: [e.direction.x, e.direction.y] as [number, number],
        state: e.state,
      }));
    }

    const elapsed = now - this.nextReceivedAt;
    const t = Math.min(1, Math.max(0, elapsed / this.snapshotIntervalMs));

    const nextIds = new Set(next.entities.map((e) => e.id));
    const result: InterpolatedEntity[] = [];

    for (const nextEnt of next.entities) {
      const prevEnt = prev.entities.find((e) => e.id === nextEnt.id);
      if (!prevEnt) {
        result.push({
          id: nextEnt.id,
          type: nextEnt.type,
          position: [nextEnt.position.x, nextEnt.position.y],
          direction: [nextEnt.direction.x, nextEnt.direction.y],
          state: nextEnt.state,
        });
        continue;
      }
      const [px, py] = lerpVec(
        prevEnt.position.x,
        prevEnt.position.y,
        nextEnt.position.x,
        nextEnt.position.y,
        t
      );
      const [dx, dy] = lerpVec(
        prevEnt.direction.x,
        prevEnt.direction.y,
        nextEnt.direction.x,
        nextEnt.direction.y,
        t
      );
      result.push({
        id: nextEnt.id,
        type: nextEnt.type,
        position: [px, py],
        direction: [dx, dy],
        state: nextEnt.state,
      });
    }
    return result;
  }

  getLatestWorld(): WorldState | null {
    return this.next;
  }
}
