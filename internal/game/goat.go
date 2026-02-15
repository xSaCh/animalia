package game

import (
	"github.com/xSaCh/animalia/internal/common"
	"github.com/xSaCh/animalia/internal/game/bt"
)

type Goat struct {
	Entity
	bt bt.Node
}

// NewGoat creates a new Goat entity with appropriate initial values
func NewGoat(id int, position common.Vector2D) *Goat {
	return &Goat{
		Entity: Entity{
			ID:       id,
			Type:     common.EntityTypeGoat,
			Position: position,
			State:    common.EntityStateRoaming,
			Direction: common.Vector2D{
				X: 0,
				Y: 0,
			},
			TargetPos: nil,
			Stats: common.Stats{
				Hunger:    30, // Starting with low hunger (30/100)
				Thirst:    25, // Starting with low thirst (25/100)
				Tiredness: 20, // Starting with low tiredness (20/100)
			},
		},
		bt: createGoatBehaviorTree(),
	}
}

func createGoatBehaviorTree() bt.Node {
	return bt.NewSelector(1, []bt.Node{
		// Handle Thirst
		bt.NewSequence(2, []bt.Node{
			bt.NewCondition(3, func(ctx *bt.TickContext) bool {
				return ctx.BlackBoard.(*Goat).Stats.Thirst >= 80
			}),
			bt.NewAction(4, func(ctx *bt.TickContext) bt.Status {
				// TODO: Find Nearby WaterSource
				return bt.Success
			}),
			bt.NewAction(5, func(ctx *bt.TickContext) bt.Status {
				// TODO: Drank Water for x ticks
				return bt.Success
			}),
		}),

		// Handle Hunger
		bt.NewSequence(5, []bt.Node{
			bt.NewCondition(6, func(ctx *bt.TickContext) bool {
				return ctx.BlackBoard.(*Goat).Stats.Hunger >= 80
			}),
			bt.NewAction(4, func(ctx *bt.TickContext) bt.Status {
				// TODO: Find Nearby FoodSource
				return bt.Success
			}),
			bt.NewAction(7, func(ctx *bt.TickContext) bt.Status {
				// TODO: Eat food for x ticks
				return bt.Success
			}),
		}),

		// Handle Tiredness
		bt.NewSequence(8, []bt.Node{
			bt.NewCondition(9, func(ctx *bt.TickContext) bool {
				return ctx.BlackBoard.(*Goat).Stats.Tiredness >= 80
			}),
			bt.NewAction(4, func(ctx *bt.TickContext) bt.Status {
				// TODO: Find Nearby resting spot
				return bt.Success
			}),
			bt.NewAction(10, func(ctx *bt.TickContext) bt.Status {
				// TODO: rest for x ticks
				return bt.Success
			}),
		}),

		// default roaming
		bt.NewSequence(11, []bt.Node{
			bt.NewAction(12, func(tc *bt.TickContext) bt.Status {
				goat := tc.BlackBoard.(*Goat)
				_ = goat
				// TODO: Find roaming position
				return bt.Success
			}),
			bt.NewAction(13, func(ctx *bt.TickContext) bt.Status {
				goat := ctx.BlackBoard.(*Goat)
				goat.MoveTowardTarget()
				return bt.Success
			}),
		}),
	})
}
