package main

import (
	"context"
	"log"

	"github.com/nqd/statum"
)

func main() {
	type state string
	type transaction string

	const (
		stateClosed state       = "closed"
		stateOpen   state       = "open"
		tranOpen    transaction = "open"
		tranClose   transaction = "close"
	)

	config := statum.NewStateMachineConfig[state, transaction]().
		AddState(stateOpen, statum.WithPermit(tranClose, stateClosed)).
		AddState(stateClosed, statum.WithPermit(tranOpen, stateOpen))

	fsm, err := statum.NewFSM(stateOpen, config)
	if err != nil {
		log.Panicln("failed to create new fsm", err)
	}

	log.Println("Current state:", fsm.Current())

	err = fsm.Fire(context.Background(), tranClose)
	if err != nil {
		log.Panicln("Failed to create trigger close transaction")
	}

	log.Println("Current state:", fsm.Current())
}
