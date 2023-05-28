package main

import (
	"context"
	"log"

	"github.com/nqd/statum"
)

type state string
type transaction string

const (
	stateClosed state = "closed"
	stateOpen   state = "open"

	tranOpen  transaction = "open"
	tranClose transaction = "close"
)

type Door struct {
	To  string
	FSM *statum.FSM[state, transaction]
}

func NewDoor(to string) *Door {
	conf := statum.NewStateMachineConfig[state, transaction]().
		AddState(stateClosed,
			statum.WithPermit(tranOpen, stateOpen, nil, nil),
		).
		AddState(stateOpen,
			statum.WithPermit(tranClose, stateClosed, nil, nil),
		)

	fsm, err := statum.NewFSM(stateClosed, conf)
	if err != nil {
		log.Panicln("failed to create new fsm", err)
	}

	return &Door{
		To:  to,
		FSM: fsm,
	}
}

func main() {
	door := NewDoor("heaven")

	err := door.FSM.Fire(context.Background(), tranOpen)
	if err != nil {
		log.Panicln("failed to send trigger open", err)
	}

	err = door.FSM.Fire(context.Background(), tranClose)
	if err != nil {
		log.Panicln("failed to send trigger close", err)
	}
}
