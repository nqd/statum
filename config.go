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
	states          map[S]*stateProperty[S, T]
	enterAnyStateCb CallbackNoReturn[S, T]
	leaveAnyStateCb Callback[S, T]
	nilCb           Callback[S, T] // the nil function used for all nil callback
	nilCbNoReturn   CallbackNoReturn[S, T]
}

type translationProperty[S, T constraints.Ordered] struct {
	toState S
}

type stateProperty[S, T constraints.Ordered] struct {
	events       map[T]*translationProperty[S, T]
	leaveStateCb Callback[S, T]         // fired when leaving current state S
	enterStateCb CallbackNoReturn[S, T] // fired when entering specific state S
}

// CallbackNoReturn is used as a callback with cancellation ability.
// Example: WithOnLeaveState, OnLeaveAnyState
type Callback[S, T constraints.Ordered] func(ctx context.Context, e *Event[S, T]) error

// CallbackNoReturn is used as a callback with no cancellation ability.
// Example: WithOnEnterState, OnEnterAnyState
type CallbackNoReturn[S, T constraints.Ordered] func(ctx context.Context, e *Event[S, T])

type StateOption[S, T constraints.Ordered] func(property *stateProperty[S, T])

func NewStateMachineConfig[S, T constraints.Ordered]() *Config[S, T] {
	states := make(map[S]*stateProperty[S, T], 0)
	nilCallback := nilCallback[S, T]
	nilCallbackNoReturn := nilCallbackNoReturn[S, T]

	return &Config[S, T]{
		states:          states,
		nilCb:           nilCallback,
		nilCbNoReturn:   nilCallbackNoReturn,
		enterAnyStateCb: nilCallbackNoReturn,
		leaveAnyStateCb: nilCallback,
	}
}

func (c *Config[S, T]) AddState(s S, opts ...StateOption[S, T]) *Config[S, T] {
	property, found := c.states[s]
	if !found {
		property = &stateProperty[S, T]{
			events:       make(map[T]*translationProperty[S, T], 0),
			leaveStateCb: c.nilCb,
			enterStateCb: c.nilCbNoReturn,
		}
		c.states[s] = property
	}

	for _, opt := range opts {
		opt(property)
	}

	return c
}

func (c *Config[S, T]) OnLeaveAnyState(f Callback[S, T]) *Config[S, T] {
	c.leaveAnyStateCb = f
	return c
}

func (c *Config[S, T]) OnEnterAnyState(f CallbackNoReturn[S, T]) *Config[S, T] {
	c.enterAnyStateCb = f
	return c
}

func WithPermit[S, T constraints.Ordered](t T, s S) StateOption[S, T] {
	return func(property *stateProperty[S, T]) {
		property.events[t] = &translationProperty[S, T]{
			toState: s,
		}
	}
}

func WithOnEnterState[S, T constraints.Ordered](f CallbackNoReturn[S, T]) StateOption[S, T] {
	return func(property *stateProperty[S, T]) {
		property.enterStateCb = f
	}
}

func WithOnLeaveState[S, T constraints.Ordered](f Callback[S, T]) StateOption[S, T] {
	return func(property *stateProperty[S, T]) {
		property.leaveStateCb = f
	}
}

func nilCallback[S, T constraints.Ordered](_ context.Context, _ *Event[S, T]) error {
	return nil
}

func nilCallbackNoReturn[S, T constraints.Ordered](_ context.Context, _ *Event[S, T]) {}
