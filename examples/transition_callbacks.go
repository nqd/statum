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
	sstart    state = "start"
	send      state = "end"
	sfinished state = "finished"

	trun    transaction = "run"
	tfinish transaction = "finish"
	treset  transaction = "reset"
)

func main() {

	config := statum.NewStateMachineConfig[state, transaction]().
		AddState(sstart, statum.WithPermit(trun, send)).
		AddState(send,
			statum.WithPermit(tfinish, sfinished),
			statum.WithPermit(treset, sstart),
			statum.WithOnEnterState(enterEnd),
		).
		AddState(sfinished,
			statum.WithPermit(treset, sstart),
			statum.WithOnLeaveState(leaveFinished),
		)

	fsm, err := statum.NewFSM[state, transaction](sstart, config)
	if err != nil {
		log.Panicln("failed to create new fsm", err)
	}

	err = fsm.Fire(context.Background(), trun)
	if err != nil {
		log.Panicln("Failed to create trigger run transaction", err)
	}

	if sstart != fsm.Current() {
		log.Panicf("Expected state to be 'start', got: '%s'\n", fsm.Current())
	}

	log.Println("Successfully ran state machine")
}

func enterEnd(ctx context.Context, e *statum.Event[state, transaction]) error {
	err := e.FSM.Fire(ctx, tfinish)
	if err != nil {
		return fmt.Errorf("fire finish: %w", err)
	}
	return nil
}

func leaveFinished(ctx context.Context, e *statum.Event[state, transaction]) error {
	if e.Src != send {
		log.Panicf("Source should have been '%s', but was '%s'\n", send, e.Src)
	}
	err := e.FSM.Fire(ctx, treset)
	if err != nil {
		return fmt.Errorf("fire reset: %w", err)
	}
	return nil
}
