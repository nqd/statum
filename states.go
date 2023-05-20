package statum

import "golang.org/x/exp/constraints"

//type States[S, T constraints.Ordered] map[S]struct {
//	Events  []event[S, T]
//	OnLeave func() // fired when leaving current state S
//	OnEnter func() // fired when entering specific state S
//}

//type States[S, T constraints.Ordered] map[S]*StateProperty[S, T]

type Config[S, T constraints.Ordered] struct {
	states map[S]*stateProperty[S, T]
}

type stateProperty[S, T constraints.Ordered] struct {
	events  map[T]S
	onLeave callback // fired when leaving current state S
	onEnter callback // fired when entering specific state S
}

type callback func()

type StateOption[S, T constraints.Ordered] func(property *stateProperty[S, T])

func NewStateMachineConfig[S, T constraints.Ordered]() *Config[S, T] {
	states := make(map[S]*stateProperty[S, T], 0)
	return &Config[S, T]{
		states: states,
	}
}

func (c *Config[S, T]) AddState(s S, opts ...StateOption[S, T]) {
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
}

func WithPermit[S, T constraints.Ordered](t T, s S) StateOption[S, T] {
	return func(property *stateProperty[S, T]) {
		property.events[t] = s
	}
}

func WithOnEnter[S, T constraints.Ordered](f callback) StateOption[S, T] {
	return func(property *stateProperty[S, T]) {
		property.onEnter = f
	}
}

func WithOnLeave[S, T constraints.Ordered](f callback) StateOption[S, T] {
	return func(property *stateProperty[S, T]) {
		property.onLeave = f
	}
}
