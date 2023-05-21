package statum

import (
	"context"

	"golang.org/x/exp/constraints"
)

// Event is the info that get passed as a reference in the Callback
type Event[S, T constraints.Ordered] struct {
	FSM         *FSM[S, T]
	Transaction T
	Src         S
	Dst         S
}

type Config[S, T constraints.Ordered] struct {
	states map[S]*stateProperty[S, T]
}

type stateProperty[S, T constraints.Ordered] struct {
	events  map[T]S
	onLeave Callback[S, T] // fired when leaving current state S
	onEnter Callback[S, T] // fired when entering specific state S
}

type Callback[S, T constraints.Ordered] func(ctx context.Context, e Event[S, T])

type StateOption[S, T constraints.Ordered] func(property *stateProperty[S, T])

func NewStateMachineConfig[S, T constraints.Ordered]() *Config[S, T] {
	states := make(map[S]*stateProperty[S, T], 0)
	return &Config[S, T]{
		states: states,
	}
}

func (c *Config[S, T]) AddState(s S, opts ...StateOption[S, T]) *Config[S, T] {
	property, found := c.states[s]
	if !found {
		property = &stateProperty[S, T]{
			events: make(map[T]S, 0),
		}
		c.states[s] = property
	}

	for _, opt := range opts {
		opt(property)
	}

	return c
}

func WithPermit[S, T constraints.Ordered](t T, s S) StateOption[S, T] {
	return func(property *stateProperty[S, T]) {
		property.events[t] = s
	}
}

func WithOnEnter[S, T constraints.Ordered](f Callback[S, T]) StateOption[S, T] {
	return func(property *stateProperty[S, T]) {
		property.onEnter = f
	}
}

func WithOnLeave[S, T constraints.Ordered](f Callback[S, T]) StateOption[S, T] {
	return func(property *stateProperty[S, T]) {
		property.onLeave = f
	}
}
