package game

import (
	"math"

	"github.com/xSaCh/animalia/internal/common"
)

// Entity represents a movable entity in the world
type Entity struct {
	ID        int                `json:"id"`
	Type      common.EntityType  `json:"type"`
	Position  common.Vector2D    `json:"position"`
	State     common.EntityState `json:"state"`
	Direction common.Vector2D    `json:"direction"` // Direction vector/velocity
	TargetPos *common.Vector2D    `json:"target_pos,omitempty"`
	Stats     common.Stats       `json:"stats"`

	prevState         common.EntityState
	lastStateChangeAt uint // How entity will have access to tick ?

}

func (e *Entity) GetNextState(currentTick uint) common.EntityState {
	stats := e.Stats
	var nextState common.EntityState
	if stats.Thirst >= 80 {
		nextState = common.EntityStateFindWater
	}
	if stats.Hunger >= 70 {
		nextState = common.EntityStateFindFood
	}
	if stats.Tiredness >= 85 {
		nextState = common.EntityStateResting
	}

	// Exit conditions
	if nextState == common.EntityStateFindWater && stats.Thirst <= 20 {
		nextState = common.EntityStateRoaming
	}
	if nextState == common.EntityStateFindFood && stats.Hunger <= 25 {
		nextState = common.EntityStateRoaming
	}
	if nextState == common.EntityStateResting && stats.Tiredness <= 30 {
		nextState = common.EntityStateRoaming
	}

	if nextState != "" {
		e.prevState = e.State
		e.lastStateChangeAt = currentTick
		e.TargetPos = nil // Reset Target Position on state change
		return nextState
	}
	return e.State
}

func (e *Entity) MoveTowardTarget() {

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

func (e *Entity) UpdateStats(currentTick uint) {
	tickDiff := currentTick - e.lastStateChangeAt
	shouldUpdate := func(n uint) int8 {
		if tickDiff%n == 0 {
			return 1
		}
		return 0
	}
	switch e.State {
	case common.EntityStateRoaming:
		e.Stats.Hunger += 1 * shouldUpdate(1)
		e.Stats.Thirst += 2 * shouldUpdate(1)
		e.Stats.Tiredness += 1 * shouldUpdate(1)
	case common.EntityStateFindFood:
		e.Stats.Hunger -= 3 * shouldUpdate(1)
		e.Stats.Thirst += 1 * shouldUpdate(1)
		e.Stats.Tiredness += 1 * shouldUpdate(1)
	case common.EntityStateFindWater:
		e.Stats.Thirst -= 4 * shouldUpdate(1)
		e.Stats.Hunger += 1 * shouldUpdate(1)
		e.Stats.Tiredness += 1 * shouldUpdate(1)
	case common.EntityStateResting:
		e.Stats.Tiredness -= 4 * shouldUpdate(1)
		e.Stats.Hunger += 1 * shouldUpdate(1)
		e.Stats.Thirst += 1 * shouldUpdate(1)
	}

	e.Stats.Hunger = int8(math.Min(100, math.Max(0, float64(e.Stats.Hunger))))
	e.Stats.Thirst = int8(math.Min(100, math.Max(0, float64(e.Stats.Thirst))))
	e.Stats.Tiredness = int8(math.Min(100, math.Max(0, float64(e.Stats.Tiredness))))
}
