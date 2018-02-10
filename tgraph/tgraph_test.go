package tgraph

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHasProperty(t *testing.T) {
	assert := assert.New(t)

	p := New()
	exists := p.HasBoolState("test")
	assert.Equal(exists, false, "this bool state should not exist")
}

func TestNopProperty(t *testing.T) {
	assert := assert.New(t)

	p := New()
	state, exists := p.GetBoolState("test", false)
	assert.Equal(exists, false, "this bool state should not exist")
	assert.Equal(state, false, "non existing bool state should be false")
}
func TestSetProperty(t *testing.T) {
	assert := assert.New(t)

	p := New()
	oldstate, exists := p.SetBoolState("test", true)
	assert.Equal(exists, false, "this bool state should not exist")
	assert.Equal(oldstate, false, "non existing bool oldstate should be false")
	oldstate, exists = p.GetBoolState("test", false)
	assert.Equal(oldstate, true, "existing bool oldstate should be true")
}

type TestTrigger struct {
	t         *testing.T
	name      string
	triggered bool
}

func (t *TestTrigger) Trigger(name string, p *Property) {
	assert := assert.New(t.t)
	assert.Equal(name, t.name, "trigger should have property name 'test'")
	t.triggered = true
}

func TestSetPropertyAndTrigger(t *testing.T) {
	testtrigger := &TestTrigger{t: t, name: "test"}

	p := New()
	p.AddTrigger("test", "test", testtrigger)
	p.SetBoolState("test", true)
	assert.Equal(t, true, testtrigger.triggered, "trigger should have been triggered")
}
