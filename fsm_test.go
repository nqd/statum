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
			fsm, err := statum.NewFSM[state, transaction](liquid, config)
			assert.Nil(t, err)
			assert.Equal(t, liquid, fsm.Current())
		})
	})

	t.Run("Fire", func(t *testing.T) {
		t.Run("should return error when the transition is invalid for current state", func(t *testing.T) {
			fsm, err := statum.NewFSM[state, transaction](liquid, config)
			assert.Nil(t, err)

			err = fsm.Fire(context.Background(), melt)
			assert.Equal(t, err, statum.ErrInvalidTransaction)
		})

		t.Run("should accept event and move to new state", func(t *testing.T) {
			fsm, err := statum.NewFSM[state, transaction](liquid, config)
			assert.Nil(t, err)

			err = fsm.Fire(context.Background(), vaporize)
			assert.Nil(t, err)

			assert.Equal(t, gas, fsm.Current())
		})

		t.Run("should trigger callbacks", func(t *testing.T) {
			count := 1

			var liquidOnEnter statum.Callback[state, transaction] = func(ctx context.Context, e statum.Event[state, transaction]) {
			}
			var liquidOnLeave statum.Callback[state, transaction] = func(ctx context.Context, e statum.Event[state, transaction]) {
				count += 1
			}
			var solidOnEnter statum.Callback[state, transaction] = func(ctx context.Context, e statum.Event[state, transaction]) {
				count *= 2
			}
			var solidOnLeave statum.Callback[state, transaction] = func(ctx context.Context, e statum.Event[state, transaction]) {
			}

			config := statum.NewStateMachineConfig[state, transaction]().
				AddState(liquid,
					statum.WithPermit(freeze, solid),
					statum.WithOnEnter(liquidOnEnter),
					statum.WithOnLeave(liquidOnLeave),
				).
				AddState(solid,
					statum.WithPermit(melt, liquid),
					statum.WithOnEnter(solidOnEnter),
					statum.WithOnLeave(solidOnLeave),
				)

			fsm, err := statum.NewFSM[state, transaction](liquid, config)
			assert.Nil(t, err)

			err = fsm.Fire(context.Background(), freeze)
			assert.Nil(t, err)

			// liquidOnLeave (count = 2) -> solidOnEnter
			assert.Equal(t, 4, count)
		})
	})

	t.Run("SetState", func(t *testing.T) {
		t.Run("should move fsm to given state", func(t *testing.T) {
			fsm, err := statum.NewFSM[state, transaction](liquid, config)
			assert.Nil(t, err)

			err = fsm.SetState(gas)
			assert.Nil(t, err)

			assert.Equal(t, gas, fsm.Current())
		})

		t.Run("should return error when state is not registered", func(t *testing.T) {
			fsm, err := statum.NewFSM[state, transaction](liquid, config)
			assert.Nil(t, err)

			err = fsm.SetState(notRegister)
			assert.Error(t, statum.ErrNotRegisteredState, err)
		})

		// todo: test not trigger callback
	})
}

//func TestFSM_Event(t *testing.T) {
//	t.Run("should accept event and move to new state", func(t *testing.T) {
//		states := statum.States[string, string]{}
//		fsm, err := statum.NewFSM[string, string](
//			"solid",
//			states,
//		)
//		assert.Nil(t, err)
//
//		err = fsm.Event(context.Background(), "melt")
//		assert.Nil(t, err)
//		assert.Equal(t, "liquid", fsm.Current())
//	})
//
//	t.Run("should return error when event is invalid", func(t *testing.T) {
//		/*
//				states := NewStateMachineConfig()
//				states.AddState(state1, WithPermit(trigger1, state11), WithPermit(trigger1, state11,
//				WithOnEnter(fu1), WithOnExit(fu2)
//			)
//					.Permit(trigger1, state11)
//					.Permit(trigger2, state22)
//					.OnEnter(fun1)
//					.OnExit(fun2)
//				states.Configure()
//		*/
//		states := make(statum.States[string, string], 0)
//		states["solid"] = &statum.StateProperty[string, string]{
//			Events: []statum.Event[string, string]{
//				{
//					Transition: "melt",
//					To:         "liquid",
//				},
//			},
//			OnEnter: nil,
//			OnLeave: nil,
//		}
//		states["liquid"] = &statum.StateProperty[string, string]{
//			Events: []statum.Event[string, string]{
//				{
//					Transition: "freeze",
//					To:         "solid",
//				},
//			},
//		}
//
//		fsm, err := statum.NewFSM[string, string](
//			"solid",
//			states,
//		)
//		assert.Nil(t, err)
//
//		err = fsm.Event(context.Background(), "invalid_transition")
//		assert.Error(t, err)
//	})
//}
