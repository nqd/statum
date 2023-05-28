//go:build ignore

package main

import (
	"context"
	"fmt"
	"log"

	"github.com/nqd/statum"
)

func main() {
	type state string
	type transaction string

	const (
		stateIdle     state = "idle"
		stateScanning state = "scanning"

		tranScan      transaction = "scan"
		tranWorking   transaction = "working"
		tranSituation transaction = "situation"
		tranFinish    transaction = "finish"
	)

	afterScan := func(ctx context.Context, e *statum.Event[state, transaction]) error {
		fmt.Println("after scan:", e.FSM.Current())
		return nil
	}
	afterWorking := func(ctx context.Context, e *statum.Event[state, transaction]) error {
		fmt.Println("after working:", e.FSM.Current())
		return nil
	}
	afterSituation := func(ctx context.Context, e *statum.Event[state, transaction]) error {
		fmt.Println("after situation:", e.FSM.Current())
		return nil
	}
	afterFinish := func(ctx context.Context, e *statum.Event[state, transaction]) error {
		fmt.Println("after finish:", e.FSM.Current())
		return nil
	}

	conf := statum.NewStateMachineConfig[state, transaction]().
		AddState(stateIdle,
			statum.WithPermit(tranScan, stateScanning, nil, afterScan),
			statum.WithPermit(tranSituation, stateIdle, nil, nil),
		).
		AddState(stateScanning,
			statum.WithPermit(tranWorking, stateScanning, nil, afterWorking),
			statum.WithPermit(tranSituation, stateScanning, nil, afterSituation),
			statum.WithPermit(tranFinish, stateIdle, nil, afterFinish),
		)

	fsm, err := statum.NewFSM(stateIdle, conf)
	if err != nil {
		log.Panicln("failed to create new fsm", err)
	}

	fmt.Println("Current state:", fsm.Current())

	err = fsm.Fire(context.Background(), tranScan)
	if err != nil {
		log.Panicln("Failed to create trigger scan transaction")
	}

	fmt.Println("1:" + fsm.Current())

	err = fsm.Fire(context.Background(), tranWorking)
	if err != nil {
		log.Panicln("Failed to create trigger working transaction")
	}

	fmt.Println("2:" + fsm.Current())

	err = fsm.Fire(context.Background(), tranSituation)
	if err != nil {
		log.Panicln("Failed to create trigger situation transaction")
	}

	fmt.Println("3:" + fsm.Current())

	err = fsm.Fire(context.Background(), tranFinish)
	if err != nil {
		log.Panicln("Failed to create trigger finish transaction")
	}

	fmt.Println("4:" + fsm.Current())
}
