package common

// EntityType represents the type of entity
type EntityType string

// EntityState represents the current state of an entity
type EntityState string
type ObstacleType string

const (
	// EntityType
	EntityTypeAnimal EntityType = "animal"
	EntityTypeAgent  EntityType = "agent"

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

type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type Vector2D struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type Stats struct {
	Hunger    uint `json:"hunger"`    // 0-100, 0 = starving, 100 = full
	Thirst    uint `json:"thirst"`    // 0-100, 0 = dehydrated, 100 = hydrated
	Tiredness uint `json:"tiredness"` // 0-100, 0 = exhausted, 100 = fully rested
}

type StaticObstacle struct {
	Type     ObstacleType `json:"type"`
	Position Position     `json:"position"`
	Size     Vector2D     `json:"size"`
}

type StaticObstacles struct {
	Walls        []StaticObstacle `json:"walls"`
	WaterSources []StaticObstacle `json:"water_sources"`
	FoodSources  []StaticObstacle `json:"food_sources"`
	RestAreas    []StaticObstacle `json:"rest_areas"`
}
