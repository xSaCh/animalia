package bt

/*
State
BTContext

Node
*/

type Status int
type ActionFn func(*TickContext) Status
type ConditionFn func(*TickContext) bool

const (
	Success Status = iota
	Failure
	Running
)

type TickContext struct {
	BlackBoard any
	NodeStates []int
}

type Node interface {
	ID() int
	Tick(*TickContext) Status
}

// Sequence Node
type Sequence struct {
	id       int
	children []Node
}

func (s *Sequence) ID() int {
	return s.id
}

func (s *Sequence) Tick(ctx *TickContext) Status {
	current := ctx.NodeStates[s.id]
	for current < len(s.children) {
		status := s.children[current].Tick(ctx)
		switch status {
		case Success:
			current++
		case Failure:
			ctx.NodeStates[s.id] = 0
			return Failure
		case Running:
			ctx.NodeStates[s.id] = current
			return Running
		}
	}
	ctx.NodeStates[s.id] = 0
	return Success
}

// Selector Node
type Selector struct {
	id       int
	children []Node
}

func (s *Selector) ID() int {
	return s.id
}

func (s *Selector) Tick(ctx *TickContext) Status {
	current := ctx.NodeStates[s.id]
	for current < len(s.children) {
		status := s.children[current].Tick(ctx)
		switch status {
		case Success:
			ctx.NodeStates[s.id] = 0
			return Success
		case Failure:
			current++
		case Running:
			ctx.NodeStates[s.id] = current
			return Running
		}
	}
	ctx.NodeStates[s.id] = 0
	return Failure
}

// Action Node
type Action struct {
	id int
	fn ActionFn
}

func (a *Action) ID() int {
	return a.id
}

func (a *Action) Tick(ctx *TickContext) Status {
	return a.fn(ctx)
}

// Condition Node
type Condition struct {
	id int
	fn ConditionFn
}

func (a *Condition) ID() int {
	return a.id
}

func (a *Condition) Tick(ctx *TickContext) Status {
	if a.fn(ctx) {
		return Success
	}
	return Failure
}

// Constructors

func NewSequence(id int, children []Node) *Sequence {
	return &Sequence{
		id:       id,
		children: children,
	}
}
func NewSelector(id int, children []Node) *Selector {
	return &Selector{
		id:       id,
		children: children,
	}
}
func NewAction(id int, fn ActionFn) *Action {
	return &Action{
		id: id,
		fn: fn,
	}
}
func NewCondition(id int, fn ConditionFn) *Condition {
	return &Condition{
		id: id,
		fn: fn,
	}
}
