import type { WorldState } from "../models/world.js";

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

interface EntityTrack {
  id: number;
  type: string;
  state: string;
  direction: [number, number];
  velocity: [number, number]; // world units / second
  authPos: [number, number]; // latest authoritative server position
  displayPos: [number, number]; // current rendered position
  correctionStartPos: [number, number]; // render position when snapshot arrived
  correctionStartTime: number;
  correctionDurationMs: number;
}

export class Interpolator {
  private latestWorld: WorldState | null = null;
  private tracks = new Map<number, EntityTrack>();
  private lastSnapshotAt = 0;
  private snapshotIntervalMs = 1000; // server event interval (network snapshot)
  private readonly minCorrectionMs = 120;
  private readonly maxCorrectionMs = 300;
  private readonly maxPredictSeconds = 1.5; // guard rail against runaway extrapolation

  update(state: WorldState): void {
    const now = performance.now();

    if (this.lastSnapshotAt > 0) {
      // Track real server event cadence (your case: ~1 second events).
      const observedMs = now - this.lastSnapshotAt;
      // EMA to reduce jitter while still adapting.
      this.snapshotIntervalMs = lerp(
        this.snapshotIntervalMs,
        Math.max(100, Math.min(2000, observedMs)),
        0.25
      );
    }

    const liveIds = new Set(state.entities.map((e) => e.id));
    // Remove tracks for entities that no longer exist.
    for (const id of this.tracks.keys()) {
      if (!liveIds.has(id)) this.tracks.delete(id);
    }

    for (const entity of state.entities) {
      const track = this.tracks.get(entity.id);

      if (!track) {
        // First time seeing this entity: seed with authoritative state.
        this.tracks.set(entity.id, {
          id: entity.id,
          type: entity.type,
          state: entity.state,
          direction: [entity.direction.x, entity.direction.y],
          velocity: [0, 0],
          authPos: [entity.position.x, entity.position.y],
          displayPos: [entity.position.x, entity.position.y],
          correctionStartPos: [entity.position.x, entity.position.y],
          correctionStartTime: now,
          correctionDurationMs: this.minCorrectionMs,
        });
        continue;
      }

      // Capture the entity's current displayed position before reconciliation.
      const currentDisplay = this.computeDisplayPos(track, now);
      track.displayPos = currentDisplay;

      const dtSec = Math.max((now - this.lastSnapshotAt) / 1000, 0.001);
      const measuredVx = (entity.position.x - track.authPos[0]) / dtSec;
      const measuredVy = (entity.position.y - track.authPos[1]) / dtSec;

      // Prefer measured velocity between authoritative snapshots.
      // If near-zero, fallback to server direction (keeps orientation/intent).
      let vx = measuredVx;
      let vy = measuredVy;
      const measuredSpeed = Math.hypot(vx, vy);
      if (measuredSpeed < 0.0001) {
        const dirLen = Math.hypot(entity.direction.x, entity.direction.y);
        if (dirLen > 0) {
          // Small default pace while still allowing "intent" movement.
          const fallbackSpeed = 1.0;
          vx = (entity.direction.x / dirLen) * fallbackSpeed;
          vy = (entity.direction.y / dirLen) * fallbackSpeed;
        } else {
          vx = 0;
          vy = 0;
        }
      }

      track.type = entity.type;
      track.state = entity.state;
      track.direction = [entity.direction.x, entity.direction.y];
      track.velocity = [vx, vy];
      track.authPos = [entity.position.x, entity.position.y];
      track.correctionStartPos = currentDisplay;
      track.correctionStartTime = now;
      track.correctionDurationMs = Math.max(
        this.minCorrectionMs,
        Math.min(this.maxCorrectionMs, this.snapshotIntervalMs * 0.35)
      );
    }

    this.latestWorld = state;
    this.lastSnapshotAt = now;
  }

  getInterpolatedEntities(now: number): InterpolatedEntity[] {
    const result: InterpolatedEntity[] = [];
    for (const track of this.tracks.values()) {
      const pos = this.computeDisplayPos(track, now);
      track.displayPos = pos;

      const speed = Math.hypot(track.velocity[0], track.velocity[1]);
      const dir: [number, number] =
        speed > 0.0001
          ? [track.velocity[0] / speed, track.velocity[1] / speed]
          : track.direction;

      result.push({
        id: track.id,
        type: track.type,
        position: pos,
        direction: dir,
        state: track.state,
      });
    }
    return result;
  }

  private computeDisplayPos(track: EntityTrack, now: number): [number, number] {
    const elapsedSinceSnapshotSec = Math.max(
      0,
      Math.min((now - track.correctionStartTime) / 1000, this.maxPredictSeconds)
    );

    // Dead-reckoned target position from latest authoritative snapshot.
    const targetX = track.authPos[0] + track.velocity[0] * elapsedSinceSnapshotSec;
    const targetY = track.authPos[1] + track.velocity[1] * elapsedSinceSnapshotSec;

    // Smoothly reconcile from previously rendered position to predicted path.
    const t = Math.max(
      0,
      Math.min(1, (now - track.correctionStartTime) / track.correctionDurationMs)
    );
    const [x, y] = lerpVec(
      track.correctionStartPos[0],
      track.correctionStartPos[1],
      targetX,
      targetY,
      t
    );
    return [x, y];
  }

  getLatestWorld(): WorldState | null {
    return this.latestWorld;
  }
}
