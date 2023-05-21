package statum

import (
	"context"

	"golang.org/x/exp/constraints"
)

type FSM[S, T constraints.Ordered] struct {
	currentState S
	states       States[S, T]
}

type States[S, T constraints.Ordered] map[S]*StateProperty[S, T]

type StateProperty[S, T constraints.Ordered] struct {
	Events  []Event[S, T]
	OnLeave func() // fired when leaving current state S
	OnEnter func() // fired when entering specific state S
}

func NewFSM[S, T constraints.Ordered](initState S, states States[S, T]) (*FSM[S, T], error) {
	return &FSM[S, T]{
		currentState: initState,
		states:       states,
	}, nil
}

// Event sends a transition trigger to fsm
func (f *FSM[S, T]) Event(ctx context.Context, t T) error {
	return nil
}

// Current returns the current fsm state
func (f *FSM[S, T]) Current() S {
	return f.currentState
}

// SetState move fsm to given state, do not trigger any Callback
func (f *FSM[S, T]) SetState(s S) {}

// Can returns true if
func (f *FSM[S, T]) Can(t T) bool {
	return false
}
