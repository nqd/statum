//go:build ignore

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
		sclosed state       = "closed"
		sopen   state       = "open"
		topen   transaction = "open"
		tclose  transaction = "close"
	)

	config := statum.NewStateMachineConfig[state, transaction]().
		AddState(sopen, statum.WithPermit(tclose, sclosed, nil, nil)).
		AddState(sclosed, statum.WithPermit(topen, sopen, nil, nil))
	fsm, err := statum.NewFSM[state, transaction](sopen, config)
	if err != nil {
		log.Panicln("failed to create new fsm", err)
	}

	log.Println("Current state:", fsm.Current())

	err = fsm.Fire(context.Background(), tclose)
	if err != nil {
		log.Panicln("Failed to create trigger close transaction")
	}

	log.Println("Current state:", fsm.Current())
}
