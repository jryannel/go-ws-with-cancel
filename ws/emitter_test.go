package ws

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmitterRegister(t *testing.T) {
	type eventA struct {
		Message string
	}
	items := []eventA{}

	e := NewEmitter[eventA]()
	stop := e.On(func(data eventA) {
		items = append(items, data)
	})
	e.Emit(eventA{"hello"})
	assert.Len(t, items, 1)
	assert.Equal(t, items[0].Message, "hello")
	stop()
	e.Emit(eventA{"world"})
	assert.Len(t, items, 1)
}

func TestClearEmitter(t *testing.T) {
	type eventA struct {
		Message string
	}
	items := []eventA{}

	e := NewEmitter[eventA]()
	e.On(func(data eventA) {
		items = append(items, data)
	})
	e.Emit(eventA{"hello"})
	assert.Len(t, items, 1)
	assert.Equal(t, items[0].Message, "hello")
	e.Clear()
	e.Emit(eventA{"world"})
	assert.Len(t, items, 1)
}
