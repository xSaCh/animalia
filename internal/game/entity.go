package game

import "github.com/xSaCh/animalia/internal/common"

// Entity represents a movable entity in the world
type Entity struct {
	ID        int                `json:"id"`
	Type      common.EntityType  `json:"type"`
	Position  common.Position    `json:"position"`
	State     common.EntityState `json:"state"`
	Direction common.Vector2D    `json:"direction"` // Direction vector/velocity
	TargetPos common.Position    `json:"target_pos,omitempty"`
	Stats     common.Stats       `json:"stats"`
}
