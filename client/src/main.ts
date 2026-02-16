import { Scene } from "./scene/Scene.js";
import { WebSocketConnection } from "./connection/websocket.js";
import { MockConnection } from "./connection/mock.js";

const canvas = document.getElementById("canvas") as HTMLCanvasElement;
if (!canvas) throw new Error("canvas not found");

const hud = document.getElementById("hud") as HTMLDivElement;
const hudTitle = document.getElementById("hud-title") as HTMLHeadingElement;
const hudStats = document.getElementById("hud-stats") as HTMLDivElement;

const scene = new Scene(canvas);

let selectedEntityId: number | null = null;

scene.onEntitySelect = (id) => {
  selectedEntityId = id;
  if (id === null) {
    hud.style.display = "none";
  } else {
    hud.style.display = "block";
    updateHud();
  }
};

function updateHud() {
  if (selectedEntityId === null) return;
  const world = scene.getLatestWorld();
  if (!world) return;
  const entity = world.entities.find((e) => e.id === selectedEntityId);
  if (!entity) {
    // Entity might have died or disappeared
    hud.style.display = "none";
    selectedEntityId = null;
    return;
  }

  hudTitle.innerText = `${entity.type} #${entity.id}`;
  hudStats.innerHTML = `
    <div class="stat-row"><span>State:</span> <span>${entity.state}</span></div>
    <div class="stat-row"><span>Pos:</span> <span>(${entity.position.x.toFixed(
      1
    )}, ${entity.position.y.toFixed(1)})</span></div>
    <div class="stat-row"><span>Hunger:</span> <span>${
      entity.stats.hunger
    }</span></div>
    <div class="stat-row"><span>Thirst:</span> <span>${
      entity.stats.thirst
    }</span></div>
    <div class="stat-row"><span>Tiredness:</span> <span>${
      entity.stats.tiredness
    }</span></div>
  `;
}

const useMock = !window.location.search.includes("ws");
const connection = useMock
  ? new MockConnection()
  : new WebSocketConnection("ws://localhost:6969/ws");

connection.onWorldState((state) => scene.updateWorldState(state));
connection.connect();

function loop(): void {
  scene.render();
  if (selectedEntityId !== null) updateHud();
  requestAnimationFrame(loop);
}
requestAnimationFrame(loop);

window.addEventListener("resize", () => {
  scene.resize(canvas.clientWidth, canvas.clientHeight);
});
