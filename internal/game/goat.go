package game

import (
	"fmt"
	"math"

	"github.com/xSaCh/animalia/internal/common"
	"github.com/xSaCh/animalia/internal/game/btree"
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

func createGoatBehaviorTree() btree.Node {
	idGen := btree.NewIDGenerator()
	
	findWaterSource := func(ctx *btree.TickContext) btree.Status {
		goat := ctx.BlackBoard.(*Goat)
		world := ctx.World.(*World)

		// Find nearest water source
		waterPos := world.GetRandomWaterSourcePos()
		if waterPos.IsZero() {
			return btree.Failure // No water source available
		}
		goat.TargetPos = &waterPos
		goat.State = common.EntityStateFindWater
		return btree.Success
	}
	findFoodSource := func(ctx *btree.TickContext) btree.Status {
		goat := ctx.BlackBoard.(*Goat)
		world := ctx.World.(*World)

		// Find nearest food source
		foodPos := world.GetRandomFoodSourcePos()
		if foodPos.IsZero() {
			return btree.Failure // No food source available
		}
		goat.TargetPos = &foodPos
		goat.State = common.EntityStateFindFood
		return btree.Success
	}
	findRestingSpot := func(ctx *btree.TickContext) btree.Status {
		goat := ctx.BlackBoard.(*Goat)
		world := ctx.World.(*World)

		// Find a random walkable position to rest
		restPos := world.GetRandomWalkablePosition()
		goat.TargetPos = &restPos
		goat.State = common.EntityStateResting
		return btree.Success
	}
	findRoamingPosition := func(ctx *btree.TickContext) btree.Status {
		goat := ctx.BlackBoard.(*Goat)
		world := ctx.World.(*World)

		// Find random roaming position
		if goat.TargetPos == nil || goat.Position.Distance(*goat.TargetPos) < 0.5 {
			roamPos := world.GetRandomWalkablePosition()
			goat.TargetPos = &roamPos
		}
		goat.State = common.EntityStateRoaming
		return btree.Success
	}

	moveToWaterAndDrink := func(ctx *btree.TickContext) btree.Status {
		goat := ctx.BlackBoard.(*Goat)

		// Check if reached target
		if goat.TargetPos != nil && goat.Position.Distance(*goat.TargetPos) > 0.5 {
			goat.MoveTowardTarget()
			goat.updateStatsDuringWalk()
			return btree.Running
		}

		fmt.Printf("Goat %d is thirsty. Moving to water source at (%.2f, %.2f)\n", goat.ID, goat.TargetPos.X, goat.TargetPos.Y)
		// Drink water
		goat.Stats.Thirst -= 2
		if goat.Stats.Thirst <= 20 {
			goat.TargetPos = nil
			return btree.Success
		}
		return btree.Running
	}

	moveToFoodAndEat := func(ctx *btree.TickContext) btree.Status {
		goat := ctx.BlackBoard.(*Goat)

		// Check if reached target
		if goat.TargetPos != nil && goat.Position.Distance(*goat.TargetPos) > 0.5 {
			goat.MoveTowardTarget()
			goat.updateStatsDuringWalk()
			return btree.Running
		}

		fmt.Printf("Goat %d is hungry. Moving to food source at (%.2f, %.2f)\n", goat.ID, goat.TargetPos.X, goat.TargetPos.Y)
		// Eat food
		goat.Stats.Hunger -= 2
		if goat.Stats.Hunger <= 25 {
			goat.TargetPos = nil
			return btree.Success
		}
		return btree.Running
	}

	moveToRestingSpotAndRest := func(ctx *btree.TickContext) btree.Status {
		goat := ctx.BlackBoard.(*Goat)

		// Check if reached resting spot
		if goat.TargetPos != nil && goat.Position.Distance(*goat.TargetPos) > 0.5 {
			goat.MoveTowardTarget()
			goat.updateStatsDuringWalk()
			return btree.Running
		}

		fmt.Printf("Goat %d is tired. Moving to resting spot at (%.2f, %.2f)\n", goat.ID, goat.TargetPos.X, goat.TargetPos.Y)
		// Rest
		goat.Stats.Tiredness -= 2
		if goat.Stats.Tiredness <= 30 {
			goat.TargetPos = nil
			return btree.Success
		}
		return btree.Running
	}

	moveWhileRoaming := func(ctx *btree.TickContext) btree.Status {
		goat := ctx.BlackBoard.(*Goat)
		world := ctx.World.(*World)

		// Check if reached target pos
		if goat.TargetPos != nil && goat.Position.Distance(*goat.TargetPos) > 0.5 {
			goat.MoveTowardTarget()
			goat.updateStatsDuringWalk()
			return btree.Running
		}
		fmt.Printf("Goat %d is idling at (%.2f, %.2f)\n", goat.ID, goat.TargetPos.X, goat.TargetPos.Y)
		if world.GetTick()%40 == 0 {
			return btree.Success
		}
		return btree.Running
	}

	return btree.NewSelector(idGen.Next(),
		// Handle Thirst
		btree.NewSequence(idGen.Next(),
			btree.NewCondition(idGen.Next(), func(ctx *btree.TickContext) bool {
				return ctx.BlackBoard.(*Goat).Stats.Thirst >= 80
			}),
			btree.NewAction(idGen.Next(), findWaterSource),
			btree.NewAction(idGen.Next(), moveToWaterAndDrink),
		),

		// Handle Hunger
		btree.NewSequence(idGen.Next(),
			btree.NewCondition(idGen.Next(), func(ctx *btree.TickContext) bool {
				return ctx.BlackBoard.(*Goat).Stats.Hunger >= 80
			}),
			btree.NewAction(idGen.Next(), findFoodSource),
			btree.NewAction(idGen.Next(), moveToFoodAndEat),
		),

		// Handle Tiredness
		btree.NewSequence(idGen.Next(),
			btree.NewCondition(idGen.Next(), func(ctx *btree.TickContext) bool {
				return ctx.BlackBoard.(*Goat).Stats.Tiredness >= 80
			}),
			btree.NewAction(idGen.Next(), findRestingSpot),
			btree.NewAction(idGen.Next(), moveToRestingSpotAndRest),
		),

		// default roaming
		btree.NewSequence(idGen.Next(),
			btree.NewAction(idGen.Next(), findRoamingPosition),
			btree.NewAction(idGen.Next(), moveWhileRoaming),
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
	ctx := &btree.TickContext{
		BlackBoard: g,
		World:      world,
		NodeStates: g.btState,
	}
	g.bt.Tick(ctx)
}
