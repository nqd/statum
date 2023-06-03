package statum_test

import (
	"context"
	"testing"

	"github.com/nqd/statum"
	"github.com/stretchr/testify/assert"
)

func TestFSM(t *testing.T) {
	type state string
	type transaction string
	var (
		solid       state = "solid"
		liquid      state = "liquid"
		gas         state = "gas"
		notRegister state = "not register state"

		melt     transaction = "melt"
		freeze   transaction = "freeze"
		vaporize transaction = "vaporize"
		condense transaction = "condense"
	)

	config := statum.NewStateMachineConfig[state, transaction]().
		AddState(liquid,
			statum.WithPermit(freeze, solid),
			statum.WithPermit(vaporize, gas)).
		AddState(gas,
			statum.WithPermit(condense, liquid)).
		AddState(solid,
			statum.WithPermit(melt, liquid))

	t.Run("Current", func(t *testing.T) {
		t.Run("should return init state", func(t *testing.T) {
			fsm, err := statum.NewFSM(liquid, config)
			assert.Nil(t, err)
			assert.Equal(t, liquid, fsm.Current())
		})
	})

	t.Run("Fire", func(t *testing.T) {
		t.Run("should return error when the transition is invalid for current state", func(t *testing.T) {
			fsm, err := statum.NewFSM(liquid, config)
			assert.Nil(t, err)

			err = fsm.Fire(context.Background(), melt)
			assert.Equal(t, err, statum.ErrInvalidTransaction)
		})

		t.Run("should accept event and move to new state", func(t *testing.T) {
			fsm, err := statum.NewFSM(liquid, config)
			assert.Nil(t, err)

			err = fsm.Fire(context.Background(), vaporize)
			assert.Nil(t, err)

			assert.Equal(t, gas, fsm.Current())
		})

		t.Run("should trigger callbacks", func(t *testing.T) {
			calledCount := 1
			notCalled := 0

			var liquidOnEnter statum.CallbackNoReturn[state, transaction] = func(ctx context.Context, e *statum.Event[state, transaction]) {
				notCalled += 1
			}
			var liquidOnLeave statum.Callback[state, transaction] = func(ctx context.Context, e *statum.Event[state, transaction]) error {
				calledCount += 1
				return nil
			}
			var solidOnEnter statum.CallbackNoReturn[state, transaction] = func(ctx context.Context, e *statum.Event[state, transaction]) {
				calledCount *= 2
			}
			var solidOnLeave statum.Callback[state, transaction] = func(ctx context.Context, e *statum.Event[state, transaction]) error {
				notCalled += 2
				return nil
			}

			config := statum.NewStateMachineConfig[state, transaction]().
				AddState(liquid,
					statum.WithPermit(freeze, solid),
					statum.WithOnEnterState(liquidOnEnter),
					statum.WithOnLeaveState(liquidOnLeave),
				).
				AddState(solid,
					statum.WithPermit(melt, liquid),
					statum.WithOnEnterState(solidOnEnter),
					statum.WithOnLeaveState(solidOnLeave),
				)

			fsm, err := statum.NewFSM(liquid, config)
			assert.Nil(t, err)

			err = fsm.Fire(context.Background(), freeze)
			assert.Nil(t, err)

			// liquidOnLeave (calledCount = 2) -> solidOnEnter (calledCount = 4)
			assert.Equal(t, 4, calledCount)
			assert.Equal(t, 0, notCalled)
		})
	})

	t.Run("SetState", func(t *testing.T) {
		t.Run("should move fsm to given state", func(t *testing.T) {
			fsm, err := statum.NewFSM(liquid, config)
			assert.Nil(t, err)

			err = fsm.SetState(gas)
			assert.Nil(t, err)

			assert.Equal(t, gas, fsm.Current())
		})

		t.Run("should return error when state is not registered", func(t *testing.T) {
			fsm, err := statum.NewFSM(liquid, config)
			assert.Nil(t, err)

			err = fsm.SetState(notRegister)
			assert.Error(t, statum.ErrNotRegisteredState, err)
		})

		// todo: test not trigger callback
	})
}
