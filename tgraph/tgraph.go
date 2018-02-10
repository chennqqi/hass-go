package tgraph

import "time"

type Trigger interface {
	Trigger(name string, p *Property)
}

type Property struct {
	B bool
	F float64
	I int64
	S string
	T time.Time
	X map[string]Trigger
}

func New() *Properties {
	p := &Properties{}
	p.P = map[string]*Property{}
	return p
}

type Properties struct {
	P map[string]*Property
}

// PROPERTIES

func (p *Properties) AddProperty(name string) *Property {
	prop := &Property{X: map[string]Trigger{}}
	p.P[name] = prop
	return prop
}

// TRIGGER

func (p *Properties) AddTrigger(name string, tag string, trigger Trigger) {
	prop, exists := p.P[name]
	if !exists {
		prop = p.AddProperty(name)
	}
	prop.X[tag] = trigger
	return
}

func (p *Properties) RemoveTrigger(name string, tag string) {
	prop, exists := p.P[name]
	if exists {
		delete(prop.X, tag)
	}
}

// TRIGGER

func (p *Property) CallTriggers(name string) {
	for _, trigger := range p.X {
		trigger.Trigger(name, p)
	}
}

// BOOLEAN

func (p *Properties) HasBoolState(name string) bool {
	_, exists := p.P[name]
	return exists
}

func (p *Properties) GetBoolState(name string, defaultstate bool) (currentstate bool, exists bool) {
	prop, exists := p.P[name]
	if exists {
		currentstate = prop.B
	} else {
		currentstate = false
	}
	return
}

func (p *Properties) SetBoolState(name string, newstate bool) (oldstate bool, existed bool) {
	prop, existed := p.P[name]
	if existed {
		oldstate = prop.I != 0
	} else {
		prop = p.AddProperty(name)
		prop.B = newstate
		oldstate = !newstate
	}
	prop.B = newstate
	if oldstate != newstate {
		prop.CallTriggers(name)
	}
	return
}

// 64-BIT FLOAT

func (p *Properties) HasFloatState(name string) bool {
	_, exists := p.P[name]
	return exists
}

func (p *Properties) GetFloatState(name string, defaultstate float64) (currentstate float64, exists bool) {
	prop, exists := p.P[name]
	if exists {
		currentstate = prop.F
	} else {
		currentstate = 0
	}
	return
}

func (p *Properties) SetFloatState(name string, newstate float64) (oldstate float64, existed bool) {
	prop, existed := p.P[name]
	if existed {
		oldstate = prop.F
	} else {
		prop = p.AddProperty(name)
		prop.F = newstate
		oldstate = newstate
	}
	prop.F = newstate
	return
}

// 64-BIT INTEGER

func (p *Properties) HasIntState(name string) bool {
	_, exists := p.P[name]
	return exists
}

func (p *Properties) GetIntState(name string, defaultstate bool) (currentstate int64, exists bool) {
	prop, exists := p.P[name]
	if exists {
		currentstate = prop.I
	} else {
		currentstate = 0
	}
	return
}

func (p *Properties) SetIntState(name string, newstate int64) (oldstate int64, existed bool) {
	prop, existed := p.P[name]
	if existed {
		oldstate = prop.I
	} else {
		prop = p.AddProperty(name)
		prop.I = newstate
		oldstate = newstate
	}
	prop.I = newstate
	return
}

// STRING

func (p *Properties) HasStringState(name string) bool {
	_, exists := p.P[name]
	return exists
}

func (p *Properties) GetStringState(name string, defaultstate string) (currentstate string, exists bool) {
	prop, exists := p.P[name]
	if exists {
		currentstate = prop.S
	} else {
		currentstate = ""
	}
	return
}

func (p *Properties) SetStringState(name string, newstate string) (oldstate string, existed bool) {
	prop, existed := p.P[name]
	if existed {
		oldstate = prop.S
	} else {
		prop = p.AddProperty(name)
		prop.S = newstate
		oldstate = newstate
	}
	prop.S = newstate
	return
}

// TIME

func (p *Properties) HasTimeState(name string) bool {
	_, exists := p.P[name]
	return exists
}

func (p *Properties) GetTimeState(name string, defaultstate time.Time) (currentstate time.Time, exists bool) {
	prop, exists := p.P[name]
	if exists {
		currentstate = prop.T
	} else {
		currentstate = time.Now()
	}
	return
}

func (p *Properties) SetTimeState(name string, newstate time.Time) (oldstate time.Time, existed bool) {
	prop, existed := p.P[name]
	if existed {
		oldstate = prop.T
	} else {
		prop = p.AddProperty(name)
		prop.T = newstate
		oldstate = newstate
	}
	prop.T = newstate
	return
}
