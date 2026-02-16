import type { Connection, WorldStateCallback } from "./types.js";
import type { WorldState } from "../models/world.js";

const DEFAULT_WS_URL = "ws://localhost:6969/ws";

export class WebSocketConnection implements Connection {
  private url: string;
  private ws: WebSocket | null = null;
  private callbacks: WorldStateCallback[] = [];

  constructor(url: string = DEFAULT_WS_URL) {
    this.url = url;
  }

  onWorldState(cb: WorldStateCallback): void {
    this.callbacks.push(cb);
  }

  connect(): void {
    this.ws = new WebSocket(this.url);
    this.ws.onmessage = (event) => {
      try {
        const state = JSON.parse(event.data as string) as WorldState;
        for (const cb of this.callbacks) cb(state);
      } catch {
        // ignore parse errors
      }
    };
    this.ws.onerror = () => {};
    this.ws.onclose = () => {
      this.ws = null;
    };
  }

  disconnect(): void {
    if (this.ws) {
      this.ws.close();
      this.ws = null;
    }
  }
}
