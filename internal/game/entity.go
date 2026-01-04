package game

import "github.com/xSaCh/animalia/internal/common"

// Entity represents a movable entity in the world
type Entity struct {
	ID        int                `json:"id"`
	Type      common.EntityType  `json:"type"`
	Position  common.Vector2D    `json:"position"`
	State     common.EntityState `json:"state"`
	Direction common.Vector2D    `json:"direction"` // Direction vector/velocity
	TargetPos common.Vector2D    `json:"target_pos,omitempty"`
	Stats     common.Stats       `json:"stats"`

	prevState common.EntityState 
	// lastStateChangeAt int // How entity will have access to tick ?
	
}

func (e *Entity) GetNextState() common.EntityState {
	e.prevState = e.State
	switch e.State {
	case common.EntityStateIdle:
		return common.EntityStateMoving
	case common.EntityStateMoving:
		return common.EntityStateIdle
	}
	//TODO: Implement finite state machine for entity behavior
	switch e.State {
	case common.EntityStateIdle:
		return common.EntityStateMoving
	case common.EntityStateMoving:
		return common.EntityStateEating
	case common.EntityStateEating:
		return common.EntityStateDrinking
	case common.EntityStateDrinking:
		return common.EntityStateResting
	case common.EntityStateResting:
		return common.EntityStateSearching
	case common.EntityStateSearching:
		return common.EntityStateIdle
	default:
		return common.EntityStateIdle
	}
}

func (e *Entity) MoveTowardTarget() {
	
	// Calculate direction vector from current position to target
	dir := e.TargetPos.Subtract(e.Position)
	distance := dir.Length()

	// If close enough to target, snap to target position
	if distance <= 0.5 {
		e.Position = e.TargetPos
		return
	}

	// Normalize direction vector and apply speed
	e.Direction.X = dir.X / distance
	e.Direction.Y = dir.Y / distance

	// Move towards target with speed of 1
	e.Position.X += e.Direction.X
	e.Position.Y += e.Direction.Y
}
