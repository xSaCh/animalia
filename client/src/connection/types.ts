import type { WorldState } from "../models/world.js";

export type WorldStateCallback = (state: WorldState) => void;

export interface Connection {
  onWorldState(cb: WorldStateCallback): void;
  connect(): void;
  disconnect(): void;
}
