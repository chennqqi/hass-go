package state

import (
	"fmt"
	"strings"
	"time"
)

// Instance holds all state of our components in multiple 'map's
type Instance struct {
	Bools   map[string]bool
	Strings map[string]string
	Floats  map[string]float64
	Times   map[string]time.Time
}

// Domain is a map that holds multiple instances like:
// "sensor"
// "report"
type Domain struct {
	Domain map[string]*Instance
}

// New constructs a new Instance
func New() *Domain {
	d := &Domain{}
	d.Domain = map[string]*Instance{}
	return d
}

func (d *Domain) Add(domain string) *Instance {
	s := &Instance{}
	s.Bools = map[string]bool{}
	s.Strings = map[string]string{}
	s.Floats = map[string]float64{}
	s.Times = map[string]time.Time{}
	d.Domain[domain] = s
	return s
}

func (d *Domain) Get(domain string) *Instance {
	s, exists := d.Domain[domain]
	if !exists {
		s = d.Add(domain)
	}
	return s
}
func (d *Domain) Clear(domain string) *Instance {
	s, exists := d.Domain[domain]
	if !exists {
		s = d.Add(domain)
	} else {
		s.Clear()
	}
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

func (s *Instance) Clear() {
	s.Bools = map[string]bool{}
	s.Strings = map[string]string{}
	s.Floats = map[string]float64{}
	s.Times = map[string]time.Time{}
}

func (s *Instance) HasBoolState(name string) bool {
	_, exists := s.Strings[name]
	return exists
}
func (s *Instance) GetBoolState(name string, theDefault bool) bool {
	v, exists := s.Bools[name]
	if exists {
		return v
	}
	s.Bools[name] = theDefault
	return theDefault
}
func (s *Instance) SetBoolState(name string, state bool) (bool, bool) {
	v, exists := s.Bools[name]
	s.Bools[name] = state
	if !exists {
		v = state
	}
	return v, exists
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

// DOMAIN

func (d *Domain) HasBoolState(domain, name string) bool {
	s, exists := d.Domain[domain]
	if !exists {
		s = d.Add(domain)
	}
	_, exists = s.Bools[name]
	return exists
}
func (d *Domain) GetBoolState(domain, name string, defaultBool bool) bool {
	s, exists := d.Domain[domain]
	if !exists {
		s = d.Add(domain)
	}
	b, exists := s.Bools[name]
	if exists {
		return b
	}
	s.Bools[name] = defaultBool
	return defaultBool
}
func (d *Domain) SetBoolState(domain, name string, state bool) (bool, bool) {
	s, exists := d.Domain[domain]
	if !exists {
		s = d.Add(domain)
	}
	b, exists := s.Bools[name]
	s.Bools[name] = state
	if !exists {
		b = state
	}
	return b, exists
}

func (d *Domain) HasStringState(domain, name string) bool {
	s, exists := d.Domain[domain]
	if !exists {
		s = d.Add(domain)
	}
	_, exists = s.Strings[name]
	return exists
}
func (d *Domain) GetStringState(domain, name string, theDefault string) string {
	s, exists := d.Domain[domain]
	if !exists {
		s = d.Add(domain)
	}
	str, exists := s.Strings[name]
	if exists {
		return str
	}
	s.Strings[name] = theDefault
	return theDefault
}
func (d *Domain) SetStringState(domain, name string, state string) (string, bool) {
	s, exists := d.Domain[domain]
	if !exists {
		s = d.Add(domain)
	}
	str, exists := s.Strings[name]
	s.Strings[name] = state
	if !exists {
		str = state
	}
	return str, exists
}

func (d *Domain) HasFloatState(domain, name string) bool {
	s, exists := d.Domain[domain]
	if !exists {
		s = d.Add(domain)
	}
	_, exists = s.Floats[name]
	return exists
}
func (d *Domain) GetFloatState(domain, name string, theDefault float64) float64 {
	s, exists := d.Domain[domain]
	if !exists {
		s = d.Add(domain)
	}
	f, exists := s.Floats[name]
	if exists {
		return f
	}
	s.Floats[name] = theDefault
	return theDefault
}
func (d *Domain) SetFloatState(domain, name string, state float64) (float64, bool) {
	s, exists := d.Domain[domain]
	if !exists {
		s = d.Add(domain)
	}
	f, exists := s.Floats[name]
	s.Floats[name] = state
	if !exists {
		f = state
	}
	return f, exists
}

func (d *Domain) HasTimeState(domain, name string) bool {
	s, exists := d.Domain[domain]
	if !exists {
		s = d.Add(domain)
	}
	_, exists = s.Times[name]
	return exists
}
func (d *Domain) GetTimeState(domain, name string, theDefault time.Time) time.Time {
	s, exists := d.Domain[domain]
	if !exists {
		s = d.Add(domain)
	}
	t, exists := s.Times[name]
	if exists {
		return t
	}
	s.Times[name] = theDefault
	return theDefault
}
func (d *Domain) SetTimeState(domain, name string, state time.Time) (time.Time, bool) {
	s, exists := d.Domain[domain]
	if !exists {
		s = d.Add(domain)
	}
	t, exists := s.Times[name]
	s.Times[name] = state
	if !exists {
		t = state
	}
	return t, exists
}

func (d *Domain) Print() {
	for k, v := range d.Domain {
		v.Print(k)
	}
}
func (d *Domain) PrintNamed(domain string) {
	s, exists := d.Domain[domain]
	if exists {
		s.Print(domain)
	}
}

func (s *Instance) Print(header string) {
	for k, v := range s.Bools {
		fmt.Printf("%s : %s = (bool)%v\n", header, k, v)
	}
	for k, v := range s.Floats {
		fmt.Printf("%s : %s = (float)%f\n", header, k, v)
	}
	for k, v := range s.Strings {
		lines := strings.Split(v, "\n")
		for ln, line := range lines {
			if ln == 0 {
				fmt.Printf("%s : %s = '%s'\n", header, k, line)
			} else {
				fmt.Printf("     %s\n", line)
			}
		}
	}
	for k, v := range s.Times {
		fmt.Printf("%s : %s = %v\n", header, k, v)
	}
}
