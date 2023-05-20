package statum_test

import (
	"context"
	"testing"

	"github.com/nqd/statum"
	"github.com/stretchr/testify/assert"
)

func TestFSM_Current(t *testing.T) {
	t.Run("should return init state", func(t *testing.T) {

		states := statum.States[string, string]{}
		fsm, err := statum.NewFSM[string, string](
			"solid",
			states,
		)
		assert.Nil(t, err)
		assert.Equal(t, "solid", fsm.Current())
	})
}

func TestFSM_Event(t *testing.T) {
	t.Run("should accept event and move to new state", func(t *testing.T) {
		states := statum.States[string, string]{}
		fsm, err := statum.NewFSM[string, string](
			"solid",
			states,
		)
		assert.Nil(t, err)

		err = fsm.Event(context.Background(), "melt")
		assert.Nil(t, err)
		assert.Equal(t, "liquid", fsm.Current())
	})

	t.Run("should return error when event is invalid", func(t *testing.T) {
		/*
				states := NewStateMachineConfig()
				states.AddState(state1, WithPermit(trigger1, state11), WithPermit(trigger1, state11,
				WithOnEnter(fu1), WithOnExit(fu2)
			)
					.Permit(trigger1, state11)
					.Permit(trigger2, state22)
					.OnEnter(fun1)
					.OnExit(fun2)
				states.Configure()
		*/
		states := make(statum.States[string, string], 0)
		states["solid"] = &statum.StateProperty[string, string]{
			Events: []statum.Event[string, string]{
				{
					Transition: "melt",
					To:         "liquid",
				},
			},
			OnEnter: nil,
			OnLeave: nil,
		}
		states["liquid"] = &statum.StateProperty[string, string]{
			Events: []statum.Event[string, string]{
				{
					Transition: "freeze",
					To:         "solid",
				},
			},
		}

		fsm, err := statum.NewFSM[string, string](
			"solid",
			states,
		)
		assert.Nil(t, err)

		err = fsm.Event(context.Background(), "invalid_transition")
		assert.Error(t, err)
	})
}
