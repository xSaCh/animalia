package bt

// How to use bt

/* Guard Behavior Tree
Root (Selector)
├─ Sequence: (player visible) -> Engage
│   ├─ Condition: player_visible?
│   └─ Selector
│       ├─ Sequence (attack)
│       │   ├─ Condition: player_in_range?
│       │   └─ Action: attack
│       └─ Action: chase
├─ Sequence: (low health) -> Heal
│   ├─ Condition: low_health?
│   └─ Action: find_cover_and_heal
└─ Action: patrol
*/

type GuardInfo struct {
	PlayerPos int

	Health int
}

type Guard struct {
	id   int
	info GuardInfo
	bt   Node
}

func NewGuard(id int, info GuardInfo) *Guard {
	isPlayerVisible := func(ctx *TickContext) bool {
		gInfo := ctx.BlackBoard.(*GuardInfo)
		if gInfo.PlayerPos >= 0 && gInfo.PlayerPos <= 10 {
			return true
		}
		return false
	}
	isPlayerInRange := func(ctx *TickContext) bool {
		gInfo := ctx.BlackBoard.(*GuardInfo)
		if gInfo.PlayerPos >= 0 && gInfo.PlayerPos <= 3 {
			return true
		}
		return false
	}

	attack := func(ctx *TickContext) Status {
		return Success
	}
	chase := func(ctx *TickContext) Status {
		return Success
	}
	lowHealth := func(ctx *TickContext) bool {
		return false
	}
	findCoverAndHeal := func(ctx *TickContext) Status {
		return Success
	}
	patrol := func(ctx *TickContext) Status {
		return Success
	}

	bt := NewSelector(1,
		NewSequence(2,
			NewCondition(3, isPlayerVisible),
			NewSelector(4,
				NewSequence(5,
					NewCondition(6, isPlayerInRange),
					NewAction(7, attack),
				),
				NewAction(8, chase),
			),
		),
		NewSequence(9,
			NewCondition(10, lowHealth),
			NewAction(11, findCoverAndHeal),
		),
		NewAction(12, patrol),
	)
	return &Guard{
		id:   id,
		info: info,
		bt:   bt,
	}
}

func usage() {
	guard := NewGuard(1, GuardInfo{PlayerPos: 5, Health: 100})
	ctx := &TickContext{
		BlackBoard: &guard.info,
		NodeStates: make([]int, 13), // allocate for node IDs 0-12
	}

	// Tick the behavior tree
	status := guard.bt.Tick(ctx)
	_ = status // use the status as needed

	// Example: update guard info and tick again
	guard.info.PlayerPos = 2 // player moved closer
	status = guard.bt.Tick(ctx)
	_ = status
}
