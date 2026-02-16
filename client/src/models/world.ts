/** World state types matching server output. Update when backend schema changes. */

export interface Vector2D {
  x: number;
  y: number;
}

export interface Stats {
  hunger: number;
  thirst: number;
  tiredness: number;
}

export interface Entity {
  id: number;
  type: string;
  position: Vector2D;
  state: string;
  direction: Vector2D;
  target_pos?: Vector2D;
  stats: Stats;
}

export interface StaticObstacle {
  type: string;
  position: Vector2D;
}

export interface StaticObstacles {
  walls: StaticObstacle[];
  water_sources: StaticObstacle[];
  food_sources: StaticObstacle[];
  rest_areas: StaticObstacle[];
}

export interface WorldConfig {
  tps: number;
}

export interface WorldState {
  id: number;
  width: number;
  height: number;
  navigation_grid: boolean[][];
  static_obstacles: StaticObstacles;
  entities: Entity[];
  config: WorldConfig;
}
