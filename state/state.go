package state

import (
	"fmt"
	"strings"
	"time"
)

type Property struct {
	Bool   bool
	String string
	Float  float64
	Time   time.Time
}

// Instance holds all state of our components in multiple 'map's
type Instance struct {
	Properties map[string]Property
}

// New constructs a new Instance
func New() *Instance {
	s := &Instance{}
	s.Properties = map[string]Property{}
	return s
}

func (s *Instance) Get(domain string) *Instance {
	return s
}

func (s *Instance) RemoveAnyStartingWith(prefix string) {
	toremove := []string{}
	for k, _ := range s.Properties {
		if strings.HasPrefix(k, prefix) {
			toremove = append(toremove, k)
		}
	}
	for _, k := range toremove {
		delete(s.Properties, k)
	}
}

// BOOLEAN

func (s *Instance) HasBoolState(name string) bool {
	_, exists := s.Properties[name]
	return exists
}
func (s *Instance) GetBoolState(name string, theDefault bool) bool {
	v, exists := s.Properties[name]
	if exists {
		return v.Bool
	}
	s.Properties[name] = Property{Bool: theDefault}
	return theDefault
}
func (s *Instance) SetBoolState(name string, state bool) (bool, bool) {
	v, exists := s.Properties[name]
	if !exists {
		v = Property{Bool: state}
	} else if v.Bool != state {
		v.Bool = state
	}
	return v.Bool, exists
}

// FLOAT

func (s *Instance) HasFloatState(name string) bool {
	_, exists := s.Properties[name]
	return exists
}
func (s *Instance) GetFloatState(name string, theDefault float64) float64 {
	v, exists := s.Properties[name]
	if exists {
		return v.Float
	}
	s.Properties[name] = Property{Float: theDefault}
	return theDefault
}
func (s *Instance) SetFloatState(name string, state float64) (float64, bool) {
	v, exists := s.Properties[name]
	if !exists {
		v = Property{Float: state}
	} else if v.Float != state {
		v.Float = state
	}
	return v.Float, exists
}

// STRING

func (s *Instance) HasStringState(name string) bool {
	_, exists := s.Properties[name]
	return exists
}
func (s *Instance) GetStringState(name string, theDefault string) string {
	v, exists := s.Properties[name]
	if exists {
		return v.String
	}
	s.Properties[name] = Property{String: theDefault}
	return theDefault
}
func (s *Instance) SetStringState(name string, state string) (string, bool) {
	v, exists := s.Properties[name]
	if !exists {
		v = Property{String: state}
	} else if v.String != state {
		v.String = state
	}
	return v.String, exists
}

// TIME

func (s *Instance) HasTimeState(name string) bool {
	_, exists := s.Properties[name]
	return exists
}
func (s *Instance) GetTimeState(name string, theDefault time.Time) time.Time {
	v, exists := s.Properties[name]
	if exists {
		return v.Time
	}
	s.Properties[name] = Property{Time: theDefault}
	return theDefault
}
func (s *Instance) SetTimeState(name string, state time.Time) (time.Time, bool) {
	v, exists := s.Properties[name]
	if !exists {
		v = Property{Time: state}
	} else if v.Time != state {
		v.Time = state
	}
	return v.Time, exists
}

// PRINT

func (s *Instance) PrintNamed(domain string) {
	for k, v := range s.Properties {
		if strings.HasPrefix(k, domain) {
			fmt.Printf("%s = (bool)%v / (float)%.2f / (time)%v \n", k, v.Bool, v.Float, v.String, v.Time)
		}
	}
	for k, v := range s.Properties {
		if strings.HasPrefix(k, domain) {
			if len(v.String) > 0 {
				lines := strings.Split(v.String, "\n")
				for ln, line := range lines {
					if ln == 0 {
						fmt.Printf("%s = '%s'\n", k, line)
					} else {
						fmt.Printf("     %s\n", line)
					}
				}
			}
		}
	}
}

func (s *Instance) Print() {
	for k, v := range s.Properties {
		fmt.Printf("%s = (bool)%v / (float)%.2f / (time)%v \n", k, v.Bool, v.Float, v.String, v.Time)
	}
	for k, v := range s.Properties {
		if len(v.String) > 0 {
			lines := strings.Split(v.String, "\n")
			for ln, line := range lines {
				if ln == 0 {
					fmt.Printf("%s = '%s'\n", k, line)
				} else {
					fmt.Printf("     %s\n", line)
				}
			}
		}
	}
}
