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
}
