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
	EntityStateRoaming  EntityState = "roaming"
	EntityStateMoving   EntityState = "moving"
	EntityStateDrinking EntityState = "drinking"
	EntityStateEating   EntityState = "eating"
	EntityStateResting  EntityState = "resting"
	EntityStateIdle     EntityState = "idle"

	// ObstacleType
	ObstacleTypeWall        ObstacleType = "wall"
	ObstacleTypeWaterSource ObstacleType = "water_source"
	ObstacleTypeFoodSource  ObstacleType = "food_source"
	ObstacleTypeRestArea    ObstacleType = "rest_area"
)

type Stats struct {
	Hunger    int8 `json:"hunger"`    // 0-100, 0 = full, 100 = starving
	Thirst    int8 `json:"thirst"`    // 0-100, 0 = hydrated, 100 = dehydrated
	Tiredness int8 `json:"tiredness"` // 0-100, 0 = fully rested, 100 = exhausted
}

type StaticObstacle struct {
	Type     ObstacleType `json:"type"`
	Position Vector2D     `json:"position"`
}

type StaticObstacles struct {
	Walls        []StaticObstacle `json:"walls"`
	WaterSources []StaticObstacle `json:"water_sources"`
	FoodSources  []StaticObstacle `json:"food_sources"`
	RestAreas    []StaticObstacle `json:"rest_areas"`
}
