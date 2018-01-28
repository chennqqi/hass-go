package state

import (
	"fmt"
	"time"
)

// Instance holds all state of our components in multiple 'map's
type Instance struct {
	Strings map[string]string
	Floats  map[string]float64
	Times   map[string]time.Time
}

// Domain is a map that holds multiple instances like:
// "sensor"
// "report"
type Domain struct {
	Domain map[string]Instance
}

// New constructs a new Instance
func New() *Instance {
	s := &Instance{}
	s.Strings = map[string]string{}
	s.Floats = map[string]float64{}
	s.Times = map[string]time.Time{}
	return s
}

// Merge will take all states of 'from' and insert them into 's'
func (s *Instance) Merge(from *Instance) {
	for k, v := range from.Strings {
		s.Strings[k] = v
	}
	for k, v := range from.Floats {
		s.Floats[k] = v
	}
	for k, v := range from.Times {
		s.Times[k] = v
	}
}

func (s *Instance) HasStringState(name string) bool {
	_, exists := s.Strings[name]
	return exists
}
func (s *Instance) GetStringState(name string, theDefault string) string {
	str, exists := s.Strings[name]
	if exists {
		return str
	}
	s.Strings[name] = theDefault
	return theDefault
}
func (s *Instance) SetStringState(name string, state string) (string, bool) {
	str, exists := s.Strings[name]
	s.Strings[name] = state
	if !exists {
		str = state
	}
	return str, exists
}

func (s *Instance) HasFloatState(name string) bool {
	_, exists := s.Floats[name]
	return exists
}
func (s *Instance) GetFloatState(name string, theDefault float64) float64 {
	f, exists := s.Floats[name]
	if exists {
		return f
	}
	s.Floats[name] = theDefault
	return theDefault
}
func (s *Instance) SetFloatState(name string, state float64) (float64, bool) {
	f, exists := s.Floats[name]
	s.Floats[name] = state
	if !exists {
		f = state
	}
	return f, exists
}

func (s *Instance) HasTimeState(name string) bool {
	_, exists := s.Times[name]
	return exists
}
func (s *Instance) GetTimeState(name string, theDefault time.Time) time.Time {
	t, exists := s.Times[name]
	if exists {
		return t
	}
	s.Times[name] = theDefault
	return theDefault
}
func (s *Instance) SetTimeState(name string, state time.Time) (time.Time, bool) {
	t, exists := s.Times[name]
	s.Times[name] = state
	if !exists {
		t = state
	}
	return t, exists
}

func (s *Instance) Print() {
	for k, v := range s.Floats {
		fmt.Printf("%s : %f\n", k, v)
	}
	for k, v := range s.Strings {
		fmt.Printf("%s : %s\n", k, v)
	}
	for k, v := range s.Times {
		fmt.Printf("%s : %v\n", k, v)
	}
}
