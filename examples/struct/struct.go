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
	d := &Door{
		To: to,
	}

	conf := statum.NewStateMachineConfig[state, transaction]().
		AddState(stateClosed,
			statum.WithPermit(tranOpen, stateOpen),
		).
		AddState(stateOpen,
			statum.WithPermit(tranClose, stateClosed),
		).
		OnEnterAnyState(d.enterState)

	fsm, err := statum.NewFSM(stateClosed, conf)
	if err != nil {
		log.Panicln("failed to create new fsm", err)
	}

	d.FSM = fsm

	return d
}

func (d *Door) enterState(_ context.Context, e *statum.Event[state, transaction]) {
	fmt.Printf("The door to %s is %s\n", d.To, e.Dst)
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
