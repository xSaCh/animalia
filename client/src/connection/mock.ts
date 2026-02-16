import type { Connection, WorldStateCallback } from "./types.js";
import type { WorldState } from "../models/world.js";
import logData from "../../a2.log?raw";

/** Replays or generates world state for development when backend is not connected. */
export class MockConnection implements Connection {
  private callbacks: WorldStateCallback[] = [];
  private intervalId: ReturnType<typeof setInterval> | null = null;
  private readonly events: WorldState[];
  private eventIndex = 0;

  constructor() {
    this.events = this.parseEvents(logData);
  }

  onWorldState(cb: WorldStateCallback): void {
    this.callbacks.push(cb);
  }

  connect(): void {
    if (this.events.length === 0) return;

    this.intervalId = setInterval(() => {
      const state = this.events[this.eventIndex];
      for (const cb of this.callbacks) cb(state);
      this.eventIndex = (this.eventIndex + 1) % this.events.length;
    }, 500);
  }

  disconnect(): void {
    if (this.intervalId !== null) {
      clearInterval(this.intervalId);
      this.intervalId = null;
    }
  }

  private parseEvents(rawLog: string): WorldState[] {
    const lines = rawLog
      .split(/\r?\n/)
      .map((line) => line.trim())
      .filter((line) => line.length > 0);

    const parsed: WorldState[] = [];
    for (const line of lines) {
      try {
        parsed.push(JSON.parse(line) as WorldState);
      } catch {
        // Ignore malformed lines in log replay.
      }
    }

    return parsed;
  }
}
