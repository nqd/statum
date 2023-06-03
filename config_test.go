package statum

import (
	"context"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewStateMachineConfig(t *testing.T) {
	type state string
	type transaction string
	var (
		solid  state = "solid"
		liquid state = "liquid"
		gas    state = "gas"

		melt     transaction = "melt"
		freeze   transaction = "freeze"
		vaporize transaction = "vaporize"
		condense transaction = "condense"
	)

	t.Run("AddState WithPermit", func(t *testing.T) {
		config := NewStateMachineConfig[state, transaction]()
		config.
			AddState(liquid,
				WithPermit(freeze, solid),
				WithPermit(vaporize, gas)).
			AddState(gas,
				WithPermit(condense, liquid)).
			AddState(solid,
				WithPermit(melt, liquid))

		assert.Equal(t, &translationProperty[state, transaction]{
			toState: liquid,
		}, config.states[solid].events[melt])
		assert.Equal(t, &translationProperty[state, transaction]{
			toState: gas,
		}, config.states[liquid].events[vaporize])
		assert.Equal(t, &translationProperty[state, transaction]{
			toState: solid,
		}, config.states[liquid].events[freeze])
		assert.Equal(t, &translationProperty[state, transaction]{
			toState: liquid,
		}, config.states[gas].events[condense])
	})

	t.Run("AddState WithOnEnterState", func(t *testing.T) {
		var cb1 CallbackNoReturn[state, transaction] = func(ctx context.Context, e *Event[state, transaction]) {}

		config := NewStateMachineConfig[state, transaction]()
		config.AddState(liquid,
			WithOnEnterState(cb1))

		assert.Equal(t, reflect.ValueOf(cb1).Pointer(),
			reflect.ValueOf(config.states[liquid].enterStateCb).Pointer())
	})

	t.Run("AddState WithOnLeaveState", func(t *testing.T) {
		var cb2 Callback[state, transaction] = func(ctx context.Context, e *Event[state, transaction]) error {
			return nil
		}
		config := NewStateMachineConfig[state, transaction]()
		config.AddState(liquid,
			WithOnLeaveState(cb2))

		assert.Equal(t, reflect.ValueOf(cb2).Pointer(),
			reflect.ValueOf(config.states[liquid].leaveStateCb).Pointer())
	})
}
