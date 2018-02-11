package tgraph

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHasProperty(t *testing.T) {
	assert := assert.New(t)

	p := New()
	exists := p.HasProperty("test")
	assert.Equal(exists, false, "this bool state should not exist")
}

func TestNopProperty(t *testing.T) {
	assert := assert.New(t)

	p := New()
	state, exists := p.GetBool("test", false)
	assert.Equal(exists, false, "this bool state should not exist")
	assert.Equal(state, false, "non existing bool state should be false")
}
func TestSetProperty(t *testing.T) {
	assert := assert.New(t)

	p := New()
	oldstate, exists := p.SetBool("test", true)
	assert.Equal(exists, false, "this bool state should not exist")
	assert.Equal(oldstate, false, "non existing bool oldstate should be false")
	oldstate, exists = p.GetBool("test", false)
	assert.Equal(oldstate, true, "existing bool oldstate should be true")
}

type TestListener struct {
	t         *testing.T
	name      string
	triggered bool
}

func (t *TestListener) OnAdded(name string, p *Property) {
	assert := assert.New(t.t)
	assert.Equal(name, t.name, "listener should have got a property with name 'test'")
	t.triggered = true
}

func (t *TestListener) OnRemoved(name string, p *Property) {
}

func (t *TestListener) OnChanged(name string, p *Property) {
}

func TestSetPropertyAndTrigger(t *testing.T) {
	testlistener := &TestListener{t: t, name: "test"}

	p := New()
	p.AddListener("test", testlistener)
	p.SetBool("test", true)
	assert.Equal(t, true, testlistener.triggered, "trigger should have been triggered")
}
