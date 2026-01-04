package game

import (
	"fmt"
	"math/rand"

	"github.com/xSaCh/animalia/internal/common"
)

type Config struct {
	TPS int `json:"tps"` // Ticks per second
}

// World represents the game world
type World struct {
	ID              int                    `json:"id"`
	Width           float64                `json:"width"`
	Height          float64                `json:"height"`
	NavigationGrid  [][]bool               `json:"navigation_grid"` // true = walkable, false = blocked
	StaticObstacles common.StaticObstacles `json:"static_obstacles"`
	Entities        []*Entity              `json:"entities"`
	Config          Config                 `json:"config"`
}

func NewWorld() *World {
	SIZE := 30
	grid := make([][]bool, SIZE)
	for i := range grid {
		grid[i] = make([]bool, SIZE)
	}
	for x := 0; x < SIZE; x++ {
		for y := 0; y < SIZE; y++ {
			grid[y][x] = true
		}
	}

	return &World{
		ID:             001,
		Width:          float64(SIZE),
		Height:         float64(SIZE),
		NavigationGrid: grid,
		// StaticObstacles: common.NewStaticObstacles(),
		Entities: make([]*Entity, 0),
		Config:   Config{TPS: 10},
	}
}

func (w *World) Tick() {

	/*
		- Excute Current State
		- Update Entity State
		- Handle Collisions
		- Update Entity Stats
	*/

	for _, e := range w.Entities {
		switch e.State {
		case common.EntityStateIdle:
			// Do nothing ???
			e.TargetPos = common.Vector2D{}
		case common.EntityStateMoving:
			if e.TargetPos.IsZero() || e.Position.SameAs(e.TargetPos) {
				e.TargetPos = w.GetRandomWalkablePosition()
			}
			e.MoveTowardTarget()
		}
	}

	for i := range w.Entities {
		if i == 0 {
			continue
		}
		shouldChangeState := rand.Int31n(100) < 60 // 60% chance to change state
		if shouldChangeState {
			w.Entities[i].State = w.Entities[i].GetNextState()
		}
	}

	//TODO: How to handle collisions?
	//TODO: How to update stats ?

}

func (w *World) GetRandomWalkablePosition() common.Vector2D {
	for {
		gridX := rand.Intn(len(w.NavigationGrid))
		gridY := rand.Intn(len(w.NavigationGrid[0]))
		pos := common.Vector2D{X: float64(gridX), Y: float64(gridY)}
		if w.NavigationGrid[gridX][gridY] {
			return pos
		}
	}
}

// Temp
func (w *World) RandomGoatEntity() *Entity {
	id := len(w.Entities) + 1
	pos := w.GetRandomWalkablePosition()
	return &Entity{
		ID:        id,
		Type:      common.EntityTypeGoat,
		Position:  pos,
		State:     common.EntityStateIdle,
		Direction: common.Vector2D{X: 0, Y: 0},
		Stats: common.Stats{
			Hunger:    100,
			Thirst:    100,
			Tiredness: 100,
		},
	}
}
func (w *World) PrintEntities() {
	for _, e := range w.Entities {
		fmt.Printf("ID: %d, Position: (%.2f, %.2f), State: %v\n", e.ID, e.Position.X, e.Position.Y, e.State)
	}
}

func (w *World) DrawAsciiWorld() {
	// Walkable means ` `
	// Obstacle means `#`
	// Entity means `<Entity_ID>` (white means idle state, green means moving state)
	// Target means `<Entity_ID>(but in red color)`

	// Create display grid
	grid := make([][]string, int(w.Height))
	for i := range grid {
		grid[i] = make([]string, int(w.Width))
		for j := range grid[i] {
			if w.NavigationGrid[i][j] {
				grid[i][j] = " "
			} else {
				grid[i][j] = "#"
			}
		}
	}

	// Place entities on grid
	for _, e := range w.Entities {
		x, y := int(e.Position.X), int(e.Position.Y)
		if x >= 0 && x < int(w.Width) && y >= 0 && y < int(w.Height) {
			if e.State == common.EntityStateMoving {
				grid[y][x] = fmt.Sprintf("\033[32m%d\033[0m", e.ID) // Green for moving
			} else {
				grid[y][x] = fmt.Sprintf("%d", e.ID) // White for idle
			}
		}

		// Place target position in red
		if !e.TargetPos.IsZero() {
			tx, ty := int(e.TargetPos.X), int(e.TargetPos.Y)
			if tx >= 0 && tx < int(w.Width) && ty >= 0 && ty < int(w.Height) {
				grid[ty][tx] = fmt.Sprintf("\033[31m%d\033[0m", e.ID)
			}
		}
	}

	// Print grid
	fmt.Println("+" + string(make([]byte, int(w.Width))) + "+")
	for _, row := range grid {
		fmt.Print("|")
		for _, cell := range row {
			if cell == " " {
				fmt.Print(" ")
			} else {
				fmt.Print(cell)
			}
		}
		fmt.Println("|")
	}
	fmt.Println("+" + string(make([]byte, int(w.Width))) + "+")
}
