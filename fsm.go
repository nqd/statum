package statum

import (
	"context"
	"errors"
	"sync"

	"golang.org/x/exp/constraints"
)

var (
	ErrInvalidTransaction = errors.New("invalid transaction")
	ErrNotRegisteredState = errors.New("state is not registered")
)

type FSM[S, T constraints.Ordered] struct {
	mu           sync.RWMutex
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
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.currentState
}

// Fire sends a transition trigger to fsm
func (f *FSM[S, T]) Fire(ctx context.Context, t T) error {
	f.mu.RLock()
	currentState := f.currentState
	f.mu.RUnlock()

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

	// callback sequence is:
	// 1. leave current state cb
	// 2. leave any state cb
	// 3. enter next state cb
	// 4. enter any state cb

	if err := f.config.states[currentState].leaveStateCb(ctx, event); err != nil {
		return err
	}

	if err := f.config.leaveAnyStateCb(ctx, event); err != nil {
		return err
	}

	f.setCurrentState(transactionProp.toState)

	f.config.states[transactionProp.toState].enterStateCb(ctx, event)

	f.config.enterAnyStateCb(ctx, event)

	return nil
}

func (f *FSM[S, T]) setCurrentState(s S) {
	f.mu.Lock()
	f.currentState = s
	f.mu.Unlock()
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
