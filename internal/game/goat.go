package game

import (
	"fmt"
	"math"

	"github.com/xSaCh/animalia/internal/common"
	"github.com/xSaCh/animalia/internal/game/bt"
)

type Goat struct {
	BaseEntity
}

// NewGoat creates a new Goat entity with appropriate initial values
func NewGoat(id int, position common.Vector2D) *Goat {
	return &Goat{
		BaseEntity: BaseEntity{
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
			bt:      createGoatBehaviorTree(),
			btState: make([]int, 20), // Initialize with enough slots for all node IDs
		},
	}
}

func createGoatBehaviorTree() bt.Node {
	findWaterSource := func(ctx *bt.TickContext) bt.Status {
		goat := ctx.BlackBoard.(*Goat)
		world := ctx.World.(*World)

		// Find nearest water source
		waterPos := world.GetRandomWaterSourcePos()
		if waterPos.IsZero() {
			return bt.Failure // No water source available
		}
		goat.TargetPos = &waterPos
		goat.State = common.EntityStateFindWater
		return bt.Success
	}
	findFoodSource := func(ctx *bt.TickContext) bt.Status {
		goat := ctx.BlackBoard.(*Goat)
		world := ctx.World.(*World)

		// Find nearest food source
		foodPos := world.GetRandomFoodSourcePos()
		if foodPos.IsZero() {
			return bt.Failure // No food source available
		}
		goat.TargetPos = &foodPos
		goat.State = common.EntityStateFindFood
		return bt.Success
	}
	findRestingSpot := func(ctx *bt.TickContext) bt.Status {
		goat := ctx.BlackBoard.(*Goat)
		world := ctx.World.(*World)

		// Find a random walkable position to rest
		restPos := world.GetRandomWalkablePosition()
		goat.TargetPos = &restPos
		goat.State = common.EntityStateResting
		return bt.Success
	}
	findRoamingPosition := func(ctx *bt.TickContext) bt.Status {
		goat := ctx.BlackBoard.(*Goat)
		world := ctx.World.(*World)

		// Find random roaming position
		if goat.TargetPos == nil || goat.Position.Distance(*goat.TargetPos) < 0.5 {
			roamPos := world.GetRandomWalkablePosition()
			goat.TargetPos = &roamPos
		}
		goat.State = common.EntityStateRoaming
		return bt.Success
	}

	moveToWaterAndDrink := func(ctx *bt.TickContext) bt.Status {
		goat := ctx.BlackBoard.(*Goat)

		// Check if reached target
		if goat.TargetPos != nil && goat.Position.Distance(*goat.TargetPos) > 0.5 {
			goat.MoveTowardTarget()
			goat.updateStatsDuringWalk()
			return bt.Running
		}

		fmt.Printf("Goat %d is thirsty. Moving to water source at (%.2f, %.2f)\n", goat.ID, goat.TargetPos.X, goat.TargetPos.Y)
		// Drink water
		goat.Stats.Thirst -= 2
		if goat.Stats.Thirst <= 20 {
			goat.TargetPos = nil
			return bt.Success
		}
		return bt.Running
	}

	moveToFoodAndEat := func(ctx *bt.TickContext) bt.Status {
		goat := ctx.BlackBoard.(*Goat)

		// Check if reached target
		if goat.TargetPos != nil && goat.Position.Distance(*goat.TargetPos) > 0.5 {
			goat.MoveTowardTarget()
			goat.updateStatsDuringWalk()
			return bt.Running
		}

		fmt.Printf("Goat %d is hungry. Moving to food source at (%.2f, %.2f)\n", goat.ID, goat.TargetPos.X, goat.TargetPos.Y)
		// Eat food
		goat.Stats.Hunger -= 2
		if goat.Stats.Hunger <= 25 {
			goat.TargetPos = nil
			return bt.Success
		}
		return bt.Running
	}

	moveToRestingSpotAndRest := func(ctx *bt.TickContext) bt.Status {
		goat := ctx.BlackBoard.(*Goat)

		// Check if reached resting spot
		if goat.TargetPos != nil && goat.Position.Distance(*goat.TargetPos) > 0.5 {
			goat.MoveTowardTarget()
			goat.updateStatsDuringWalk()
			return bt.Running
		}

		fmt.Printf("Goat %d is tired. Moving to resting spot at (%.2f, %.2f)\n", goat.ID, goat.TargetPos.X, goat.TargetPos.Y)
		// Rest
		goat.Stats.Tiredness -= 2
		if goat.Stats.Tiredness <= 30 {
			goat.TargetPos = nil
			return bt.Success
		}
		return bt.Running
	}

	moveWhileRoaming := func(ctx *bt.TickContext) bt.Status {
		goat := ctx.BlackBoard.(*Goat)
		world := ctx.World.(*World)

		// Check if reached target pos
		if goat.TargetPos != nil && goat.Position.Distance(*goat.TargetPos) > 0.5 {
			goat.MoveTowardTarget()
			goat.updateStatsDuringWalk()
			return bt.Running
		}
		fmt.Printf("Goat %d is idling at (%.2f, %.2f)\n", goat.ID, goat.TargetPos.X, goat.TargetPos.Y)
		if world.GetTick()%40 == 0 {
			return bt.Success
		}
		return bt.Running
	}

	return bt.NewSelector(1,
		// Handle Thirst
		bt.NewSequence(2,
			bt.NewCondition(3, func(ctx *bt.TickContext) bool {
				return ctx.BlackBoard.(*Goat).Stats.Thirst >= 80
			}),
			bt.NewAction(4, findWaterSource),
			bt.NewAction(5, moveToWaterAndDrink),
		),

		// Handle Hunger
		bt.NewSequence(6,
			bt.NewCondition(7, func(ctx *bt.TickContext) bool {
				return ctx.BlackBoard.(*Goat).Stats.Hunger >= 80
			}),
			bt.NewAction(8, findFoodSource),
			bt.NewAction(9, moveToFoodAndEat),
		),

		// Handle Tiredness
		bt.NewSequence(10,
			bt.NewCondition(11, func(ctx *bt.TickContext) bool {
				return ctx.BlackBoard.(*Goat).Stats.Tiredness >= 80
			}),
			bt.NewAction(12, findRestingSpot),
			bt.NewAction(13, moveToRestingSpotAndRest),
		),

		// default roaming
		bt.NewSequence(14,
			bt.NewAction(15, findRoamingPosition),
			bt.NewAction(16, moveWhileRoaming),
		),
	)
}

func (g *Goat) updateStatsDuringWalk() {
	g.Stats.Hunger += 1
	g.Stats.Thirst += 2
	g.Stats.Tiredness += 1

	g.Stats.Hunger = int8(math.Min(100, math.Max(0, float64(g.Stats.Hunger))))
	g.Stats.Thirst = int8(math.Min(100, math.Max(0, float64(g.Stats.Thirst))))
	g.Stats.Tiredness = int8(math.Min(100, math.Max(0, float64(g.Stats.Tiredness))))
}

// Tick executes the ent's behavior tree
func (g *Goat) Tick(world *World) {
	ctx := &bt.TickContext{
		BlackBoard: g,
		World:      world,
		NodeStates: g.btState,
	}
	g.bt.Tick(ctx)
}
