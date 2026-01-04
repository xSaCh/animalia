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

	tick uint
}

func NewWorld() *World {
	SIZE := 30
	grid := make([][]bool, SIZE)
	for i := range grid {
		grid[i] = make([]bool, SIZE)
	}
	for x := range SIZE {
		for y := range SIZE {
			grid[y][x] = true
		}
	}
	// Random 5 water sources
	waters := make([]common.StaticObstacle, 0)
	for range 5 {
		x := rand.Intn(SIZE)
		y := rand.Intn(SIZE)
		waters = append(waters, common.StaticObstacle{
			Position: common.Vector2D{X: float64(x), Y: float64(y)},
			Type:     common.ObstacleTypeWaterSource})
		grid[y][x] = false
	}
	// Random 10 food sources
	foods := make([]common.StaticObstacle, 0)
	for range 10 {
		x := rand.Intn(SIZE)
		y := rand.Intn(SIZE)
		foods = append(foods, common.StaticObstacle{
			Position: common.Vector2D{X: float64(x), Y: float64(y)},
			Type:     common.ObstacleTypeFoodSource})
		grid[y][x] = false
	}
	// Random 10 obstacles
	obstacles := make([]common.StaticObstacle, 0)
	for range 10 {
		x := rand.Intn(SIZE)
		y := rand.Intn(SIZE)
		obstacles = append(obstacles, common.StaticObstacle{
			Position: common.Vector2D{X: float64(x), Y: float64(y)},
			Type:     common.ObstacleTypeWall})
		grid[y][x] = false
	}
	return &World{
		ID:             001,
		Width:          float64(SIZE),
		Height:         float64(SIZE),
		NavigationGrid: grid,
		StaticObstacles: common.StaticObstacles{
			Walls:        obstacles,
			WaterSources: waters,
			FoodSources:  foods,
			RestAreas:    make([]common.StaticObstacle, 0),
		},
		Entities: make([]*Entity, 0),
		Config:   Config{TPS: 10},
	}
}

func (w *World) GetTick() uint {
	return w.tick
}

func (w *World) Tick() {

	/*
		- Excute Current State
		- Update Entity State
		- Handle Collisions
		- Update Entity Stats
	*/
	w.tick++

	for _, e := range w.Entities {
		switch e.State {
		case common.EntityStateRoaming:
			if e.TargetPos.IsZero() || e.Position.SameAs(e.TargetPos) {
				e.TargetPos = w.GetRandomWalkablePosition()
			}
			e.MoveTowardTarget()
		case common.EntityStateFindFood:
			if e.TargetPos.IsZero() || e.Position.SameAs(e.TargetPos) {
				e.TargetPos = w.GetRandomFoodSourcePos()
			}
			e.MoveTowardTarget()
		case common.EntityStateFindWater:
			if e.TargetPos.IsZero() || e.Position.SameAs(e.TargetPos) {
				e.TargetPos = w.GetRandomWaterSourcePos()
			}
			e.MoveTowardTarget()
		}
	}

	for i := range w.Entities {
		if i == 0 {
			continue
		}
		w.Entities[i].State = w.Entities[i].GetNextState(w.tick)

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

func (w *World) GetRandomWaterSourcePos() common.Vector2D {
	if len(w.StaticObstacles.WaterSources) == 0 {
		return common.Vector2D{}
	}
	return w.StaticObstacles.WaterSources[rand.Intn(len(w.StaticObstacles.WaterSources))].Position
}

func (w *World) GetRandomFoodSourcePos() common.Vector2D {
	if len(w.StaticObstacles.FoodSources) == 0 {
		return common.Vector2D{}
	}
	return w.StaticObstacles.FoodSources[rand.Intn(len(w.StaticObstacles.FoodSources))].Position
}

// Temp
func (w *World) RandomGoatEntity() *Entity {
	id := len(w.Entities) + 1
	pos := w.GetRandomWalkablePosition()
	return &Entity{
		ID:        id,
		Type:      common.EntityTypeGoat,
		Position:  pos,
		State:     common.EntityStateRoaming,
		Direction: common.Vector2D{X: 0, Y: 0},
		Stats: common.Stats{
			Hunger:    0,
			Thirst:    0,
			Tiredness: 0,
		},
	}
}
func (w *World) PrintEntities() {
	for _, e := range w.Entities {
		fmt.Printf("ID: %d, Position: (%.2f, %.2f), State: %v | %v\n", e.ID, e.Position.X, e.Position.Y, e.State, e.State == "")
	}
}

func (w *World) DrawAsciiWorld() {
	// Walkable means ` `
	// Obstacle means `#` (white means wall, blue means water, orange means food)
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
	// Place obstacles on grid
	for _, o := range w.StaticObstacles.Walls {
		x, y := int(o.Position.X), int(o.Position.Y)
		if x >= 0 && x < int(w.Width) && y >= 0 && y < int(w.Height) {
			grid[y][x] = "#"
		}
	}

	// Place water on grid
	for _, o := range w.StaticObstacles.WaterSources {
		x, y := int(o.Position.X), int(o.Position.Y)
		if x >= 0 && x < int(w.Width) && y >= 0 && y < int(w.Height) {
			grid[y][x] = "\033[34m#\033[0m" // Blue #
		}
	}
	// Place food on grid
	for _, o := range w.StaticObstacles.FoodSources {
		x, y := int(o.Position.X), int(o.Position.Y)
		if x >= 0 && x < int(w.Width) && y >= 0 && y < int(w.Height) {
			grid[y][x] = "\033[33m#\033[0m" // Orange #
		}
	}

	// Place entities on grid
	for _, e := range w.Entities {
		x, y := int(e.Position.X), int(e.Position.Y)
		if x >= 0 && x < int(w.Width) && y >= 0 && y < int(w.Height) {
			if e.State == common.EntityStateRoaming {
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
