package fsm

import (
	"context"

	"golang.org/x/exp/constraints"
)

type FSM[S, T constraints.Ordered] struct {
	currentState S
	states       states[S, T]
}

type event[S, T constraints.Ordered] struct {
	transition T
	to         S
}

type states[S, T constraints.Ordered] map[S]struct {
	events  []event[S, T]
	onLeave func() // fired when leaving current state S
	onEnter func() // fired when entering specific state S
}

func NewFSM[S, T constraints.Ordered](initState S, states states[S, T]) (*FSM[S, T], error) {
	return &FSM[S, T]{
		currentState: initState,
		states:       states[S, T],
	}, nil
}

// Event sends a transition trigger to fsm
func (f *FSM[S, T]) Event(ctx context.Context, t T) {}

// Current returns the current fsm state
func (f *FSM[S, T]) Current() S {}

// SetState move fsm to given state, do not trigger any callback
func (f *FSM[S, T]) SetState(s S) {}

// Can returns true if
func (f *FSM[S, T]) Can(t T) bool {
	return false
}
