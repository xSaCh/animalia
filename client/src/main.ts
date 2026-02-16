import { Scene } from "./scene/Scene.js";
import { WebSocketConnection } from "./connection/websocket.js";
import { MockConnection } from "./connection/mock.js";

const canvas = document.getElementById("canvas") as HTMLCanvasElement;
if (!canvas) throw new Error("canvas not found");

const scene = new Scene(canvas);

const useMock = !window.location.search.includes("ws");
const connection = useMock
  ? new MockConnection()
  : new WebSocketConnection("ws://localhost:6969/ws");

connection.onWorldState((state) => scene.updateWorldState(state));
connection.connect();

function loop(): void {
  scene.render();
  requestAnimationFrame(loop);
}
requestAnimationFrame(loop);

window.addEventListener("resize", () => {
  scene.resize(canvas.clientWidth, canvas.clientHeight);
});
