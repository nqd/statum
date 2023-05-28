//--go:build ignore

package main

import (
	"context"
	"fmt"
	"log"

	"github.com/nqd/statum"
)

type state string
type transaction string

const (
	stateStart    state = "start"
	stateEnd      state = "end"
	stateFinished state = "finished"

	tranRun    transaction = "run"
	tranFinish transaction = "finish"
	tranReset  transaction = "reset"
)

func main() {
	config := statum.NewStateMachineConfig[state, transaction]().
		AddState(stateStart, statum.WithPermit(tranRun, stateEnd, nil, nil)).
		AddState(stateEnd,
			statum.WithPermit(tranFinish, stateFinished, nil, afterFinishTransaction),
			statum.WithPermit(tranReset, stateStart, nil, nil),
			statum.WithOnEnterState(enterEnd),
		).
		AddState(stateFinished,
			statum.WithPermit(tranReset, stateStart, nil, nil),
		)

	fsm, err := statum.NewFSM[state, transaction](stateStart, config)
	if err != nil {
		log.Panicln("failed to create new fsm", err)
	}

	err = fsm.Fire(context.Background(), tranRun)
	if err != nil {
		log.Panicln("Failed to create trigger run transaction", err)
	}

	if stateStart != fsm.Current() {
		log.Panicf("Expected state to be 'start', got: '%s'\n", fsm.Current())
	}

	log.Println("Successfully ran state machine")
}

func enterEnd(ctx context.Context, e *statum.Event[state, transaction]) error {
	err := e.FSM.Fire(ctx, tranFinish)
	if err != nil {
		return fmt.Errorf("fire finish: %w", err)
	}
	return nil
}

func afterFinishTransaction(ctx context.Context, e *statum.Event[state, transaction]) error {
	if e.Src != stateEnd {
		log.Panicf("Source should have been '%s', but was '%s'\n", stateEnd, e.Src)
	}
	err := e.FSM.Fire(ctx, tranReset)
	if err != nil {
		return fmt.Errorf("fire reset: %w", err)
	}
	return nil
}
