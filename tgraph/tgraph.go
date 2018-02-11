package tgraph

import "time"

type EventListener interface {
	OnAdded(name string, p *Property)
	OnRemoved(name string, p *Property)
	OnChanged(name string, p *Property)
}

type Property struct {
	B bool
	F float64
	I int64
	S string
	T time.Time
}

type Properties struct {
	P map[string]*Property
	L map[string]EventListener
}

func New() *Properties {
	p := &Properties{}
	p.P = map[string]*Property{}
	return p
}

// PROPERTIES

func (p *Properties) HasProperty(name string) bool {
	_, exists := p.P[name]
	return exists
}

func (p *Properties) AddProperty(name string) *Property {
	prop := &Property{}
	p.P[name] = prop
	p.onAdded(name, prop)
	return prop
}

func (p *Properties) RemoveProperty(name string) (removed bool) {
	prop, exists := p.P[name]
	if exists {
		p.onRemoved(name, prop)
		delete(p.P, name)
	}
	return exists
}

// TRIGGER

func (p *Properties) AddListener(tag string, listener EventListener) {
	p.L[tag] = listener
	return
}

func (p *Properties) RemoveListener(tag string) {
	delete(p.L, tag)
}

func (p *Properties) onAdded(name string, property *Property) {
	for _, listener := range p.L {
		listener.OnAdded(name, property)
	}
}

func (p *Properties) onRemoved(name string, property *Property) {
	for _, listener := range p.L {
		listener.OnRemoved(name, property)
	}
}

func (p *Properties) onChanged(name string, property *Property) {
	for _, listener := range p.L {
		listener.OnChanged(name, property)
	}
}

// BOOLEAN

func (p *Properties) GetBool(name string, defaultstate bool) (currentstate bool, exists bool) {
	prop, exists := p.P[name]
	if exists {
		currentstate = prop.B
	} else {
		currentstate = false
	}
	return
}

func (p *Properties) SetBool(name string, newstate bool) (oldstate bool, existed bool) {
	prop, existed := p.P[name]
	if existed {
		oldstate = prop.I != 0
		if oldstate != newstate {
			p.onChanged(name, prop)
		}
	} else {
		prop = p.AddProperty(name)
	}
	prop.B = newstate
	return
}

// 64-BIT FLOAT

func (p *Properties) GetFloat(name string, defaultstate float64) (currentstate float64, exists bool) {
	prop, exists := p.P[name]
	if exists {
		currentstate = prop.F
	} else {
		currentstate = 0
	}
	return
}

func (p *Properties) SetFloat(name string, newstate float64) (oldstate float64, existed bool) {
	prop, existed := p.P[name]
	if existed {
		oldstate = prop.F
		if oldstate != newstate {
			p.onChanged(name, prop)
		}
	} else {
		prop = p.AddProperty(name)
	}
	prop.F = newstate
	return
}

// 64-BIT INTEGER

func (p *Properties) GetInt(name string, defaultstate bool) (currentstate int64, exists bool) {
	prop, exists := p.P[name]
	if exists {
		currentstate = prop.I
	} else {
		currentstate = 0
	}
	return
}

func (p *Properties) SetInt(name string, newstate int64) (oldstate int64, existed bool) {
	prop, existed := p.P[name]
	if existed {
		oldstate = prop.I
		if oldstate != newstate {
			p.onChanged(name, prop)
		}
	} else {
		prop = p.AddProperty(name)
	}
	prop.I = newstate
	return
}

// STRING

func (p *Properties) GetString(name string, defaultstate string) (currentstate string, exists bool) {
	prop, exists := p.P[name]
	if exists {
		currentstate = prop.S
	} else {
		currentstate = ""
	}
	return
}

func (p *Properties) SetString(name string, newstate string) (oldstate string, existed bool) {
	prop, existed := p.P[name]
	if existed {
		oldstate = prop.S
		if oldstate != newstate {
			p.onChanged(name, prop)
		}
	} else {
		prop = p.AddProperty(name)
	}
	prop.S = newstate
	return
}

// TIME

func (p *Properties) GetTime(name string, defaultstate time.Time) (currentstate time.Time, exists bool) {
	prop, exists := p.P[name]
	if exists {
		currentstate = prop.T
	} else {
		currentstate = time.Now()
	}
	return
}

func (p *Properties) SetTime(name string, newstate time.Time) (oldstate time.Time, existed bool) {
	prop, existed := p.P[name]
	if existed {
		oldstate = prop.T
		if oldstate != newstate {
			p.onChanged(name, prop)
		}
	} else {
		prop = p.AddProperty(name)
	}
	prop.T = newstate
	return
}
