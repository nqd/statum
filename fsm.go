package fsm

import "golang.org/x/exp/constraints"

type FSM[S, T constraints.Ordered] struct {
	currentState S
	states       states[S, T]
}

type transition[S, T constraints.Ordered] struct {
	trigger T
	to      S
}

type states[S, T constraints.Ordered] map[S]struct {
	transitions []transition[S, T]
	onLeave     func() // fired when leaving current state S
	onEnter     func() // fired when entering specific state S
}

func NewFSM[S, T constraints.Ordered](initState S, states states[S, T]) (*FSM[S, T], error) {
	return &FSM[S, T]{
		currentState: initState,
		states:       states[S, T],
	}, nil
}
