package statum

import (
	"context"
	"errors"

	"golang.org/x/exp/constraints"
)

var (
	ErrInvalidTransaction = errors.New("invalid transaction")
	ErrNotRegisteredState = errors.New("state is not registered")
)

type FSM[S, T constraints.Ordered] struct {
	currentState S
	config       *Config[S, T]
}

func NewFSM[S, T constraints.Ordered](initState S, config *Config[S, T]) (*FSM[S, T], error) {
	return &FSM[S, T]{
		currentState: initState,
		config:       config,
	}, nil
}

// Current returns the current fsm state
func (f *FSM[S, T]) Current() S {
	return f.currentState
}

// Fire sends a transition trigger to fsm
func (f *FSM[S, T]) Fire(ctx context.Context, t T) error {
	currentState := f.currentState

	events := f.config.states[currentState].events

	transactionProp, found := events[t]
	if !found {
		return ErrInvalidTransaction
	}

	event := &Event[S, T]{
		FSM:         f,
		Transaction: t,
		Src:         currentState,
		Dst:         transactionProp.toState,
	}

	// within the current state, the callbacks sequence is: (a) beforeTransaction, (b) leaveState
	if err := f.config.states[currentState].events[t].beforeTransaction(ctx, event); err != nil {
		return err
	}

	if err := f.config.states[currentState].leaveStateCb(ctx, event); err != nil {
		return err
	}

	f.setCurrentState(transactionProp.toState)

	// after move to next state, the callbacks sequence is: (c) enterState, (d) afterTransaction
	if err := f.config.states[transactionProp.toState].enterStateCb(ctx, event); err != nil {
		return err
	}
	if err := f.config.states[currentState].events[t].afterTransaction(ctx, event); err != nil {
		return err
	}

	return nil
}

func (f *FSM[S, T]) setCurrentState(s S) {
	f.currentState = s
}

// SetState move fsm to given state, do not trigger any Callback
func (f *FSM[S, T]) SetState(s S) error {
	_, found := f.config.states[s]
	if !found {
		return ErrNotRegisteredState
	}

	f.setCurrentState(s)
	return nil
}

//// Can returns true if
//func (f *FSM[S, T]) Can(t T) bool {
//	return false
//}
