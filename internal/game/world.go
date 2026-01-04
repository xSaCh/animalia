package game

import "github.com/xSaCh/animalia/internal/common"

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
	Entities        []Entity               `json:"entities"`
	Config          Config                 `json:"config"`
}

func NewWorld() *World {
	return &World{
		ID:             001,
		Width:          100,
		Height:         100,
		NavigationGrid: make([][]bool, 100),
		// StaticObstacles: common.NewStaticObstacles(),
		Entities: make([]Entity, 0),
		Config:   Config{TPS: 10},
	}
}

func (w *World) Tick() {
	/*    - Update Entity State
	- Update Entity Position
	- Handle Collisions
	- Update Entity Stats
	*/
}
