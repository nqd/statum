package statum

import (
	"context"
	"reflect"
	"runtime"
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

	assertTranslationProperty := func(t *testing.T, t1 translationProperty[state, transaction], t2 translationProperty[state, transaction]) bool {
		assert.Equal(t, t1.toState, t1.toState)
		assertTwoFunsEqual(t, t1.afterTransaction, t2.afterTransaction)
		assertTwoFunsEqual(t, t1.beforeTransaction, t2.beforeTransaction)

		return true
	}

	t.Run("AddState WithPermit", func(t *testing.T) {
		var cb1 Callback[state, transaction] = func(ctx context.Context, e *Event[state, transaction]) error { return nil }
		cb2 := func(ctx context.Context, e *Event[state, transaction]) error { return nil }

		config := NewStateMachineConfig[state, transaction]()
		config.
			AddState(liquid,
				WithPermit(freeze, solid, cb1, nil),
				WithPermit(vaporize, gas, nil, nil)).
			AddState(gas,
				WithPermit(condense, liquid, nil, cb2)).
			AddState(solid,
				WithPermit(melt, liquid, nil, nil))

		assertTranslationProperty(t, translationProperty[state, transaction]{
			toState:           liquid,
			afterTransaction:  nil,
			beforeTransaction: nil,
		}, config.states[solid].events[melt])
		assertTranslationProperty(t, translationProperty[state, transaction]{
			toState:           gas,
			afterTransaction:  nil,
			beforeTransaction: nil,
		}, config.states[liquid].events[vaporize])
		assertTranslationProperty(t, translationProperty[state, transaction]{
			toState:           solid,
			afterTransaction:  nil,
			beforeTransaction: cb1,
		}, config.states[liquid].events[freeze])
		assertTranslationProperty(t, translationProperty[state, transaction]{
			toState:           liquid,
			afterTransaction:  cb2,
			beforeTransaction: nil,
		}, config.states[gas].events[condense])
	})

	t.Run("AddState WithOnEnterState", func(t *testing.T) {
		var cb1 Callback[state, transaction] = func(ctx context.Context, e *Event[state, transaction]) error {
			return nil
		}
		config := NewStateMachineConfig[state, transaction]()
		config.AddState(liquid,
			WithOnEnterState(cb1))

		assert.Equal(t, reflect.ValueOf(cb1).Pointer(),
			reflect.ValueOf(config.states[liquid].onEnterState).Pointer())
	})

	t.Run("AddState WithOnLeaveState", func(t *testing.T) {
		var cb2 Callback[state, transaction] = func(ctx context.Context, e *Event[state, transaction]) error {
			return nil
		}
		config := NewStateMachineConfig[state, transaction]()
		config.AddState(liquid,
			WithOnLeaveState(cb2))

		assert.Equal(t, reflect.ValueOf(cb2).Pointer(),
			reflect.ValueOf(config.states[liquid].onLeaveState).Pointer())
	})
}

func assertTwoFunsEqual(t *testing.T, func1, func2 interface{}) bool {
	funcName1 := runtime.FuncForPC(reflect.ValueOf(func1).Pointer()).Name()
	funcName2 := runtime.FuncForPC(reflect.ValueOf(func2).Pointer()).Name()

	assert.Equal(t, funcName1, funcName2)

	return true
}
