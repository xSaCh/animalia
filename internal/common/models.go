package common

// EntityType represents the type of entity
type EntityType string

// EntityState represents the current state of an entity
type EntityState string
type ObstacleType string

const (
	EntityTypeGoat EntityType = "goat"
	EntityTypeWolf EntityType = "wolf"

	// EntityState
	EntityStateIdle      EntityState = "idle"
	EntityStateMoving    EntityState = "moving"
	EntityStateEating    EntityState = "eating"
	EntityStateDrinking  EntityState = "drinking"
	EntityStateResting   EntityState = "resting"
	EntityStateSearching EntityState = "searching"

	// ObstacleType
	ObstacleTypeWall        ObstacleType = "wall"
	ObstacleTypeWaterSource ObstacleType = "water_source"
	ObstacleTypeFoodSource  ObstacleType = "food_source"
	ObstacleTypeRestArea    ObstacleType = "rest_area"
)

type Stats struct {
	Hunger    uint `json:"hunger"`    // 0-100, 0 = starving, 100 = full
	Thirst    uint `json:"thirst"`    // 0-100, 0 = dehydrated, 100 = hydrated
	Tiredness uint `json:"tiredness"` // 0-100, 0 = exhausted, 100 = fully rested
}

type StaticObstacle struct {
	Type     ObstacleType `json:"type"`
	Position Vector2D     `json:"position"`
}

type StaticObstacles struct {
	Walls        map[Vector2D]StaticObstacle `json:"walls"`
	WaterSources map[Vector2D]StaticObstacle `json:"water_sources"`
	FoodSources  map[Vector2D]StaticObstacle `json:"food_sources"`
	RestAreas    map[Vector2D]StaticObstacle `json:"rest_areas"`
}
