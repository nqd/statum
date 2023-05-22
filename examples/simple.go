//go:build ignore

package main

import "github.com/nqd/statum"

func main() {
	type state string
	type transaction string

	const (
		closed state       = "closed"
		open   state       = "open"
		open   transaction = "open"
		close  transaction = "close"
	)

	config := statum.NewStateMachineConfig[state, transaction]().
		AddState(open, statum.WithPermit(close, closed)).
		AddState(closed, statum.WithPermit(open, open))
	statum.NewFSM[state, transaction](closed, config)
}
