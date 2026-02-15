package game

import (
	"github.com/xSaCh/animalia/internal/common"
	"github.com/xSaCh/animalia/internal/game/bt"
)

type Entity interface {
	Tick(world *World)
	GetBaseEntity() *BaseEntity
}

// BaseEntity represents a movable entity in the world
type BaseEntity struct {
	ID        int                `json:"id"`
	Type      common.EntityType  `json:"type"`
	Position  common.Vector2D    `json:"position"`
	State     common.EntityState `json:"state"`
	Direction common.Vector2D    `json:"direction"` // Direction vector/velocity
	TargetPos *common.Vector2D   `json:"target_pos,omitempty"`
	Stats     common.Stats       `json:"stats"`

	prevState         common.EntityState
	lastStateChangeAt uint // How entity will have access to tick ?

	bt      bt.Node
	btState []int // Track state for each node in behavior tree

}

func (e *BaseEntity) GetBaseEntity() *BaseEntity {
	return e
}

func (e *BaseEntity) MoveTowardTarget() {

	// Calculate direction vector from current position to target
	dir := e.TargetPos.Subtract(e.Position)
	distance := dir.Length()

	// If close enough to target, snap to target position
	if distance <= 0.5 {
		e.Position = *e.TargetPos
		return
	}

	// Normalize direction vector and apply speed
	e.Direction.X = dir.X / distance
	e.Direction.Y = dir.Y / distance

	// Move towards target with speed of 1
	e.Position.X += e.Direction.X
	e.Position.Y += e.Direction.Y
}