package main

import (
	"fmt"
	"time"

	"github.com/jurgen-kluft/hass-go/aqi"
	"github.com/jurgen-kluft/hass-go/calendar"
	"github.com/jurgen-kluft/hass-go/hass"
	"github.com/jurgen-kluft/hass-go/lighting"
	"github.com/jurgen-kluft/hass-go/reporter"
	"github.com/jurgen-kluft/hass-go/sensors"
	"github.com/jurgen-kluft/hass-go/shout"
	"github.com/jurgen-kluft/hass-go/state"
	"github.com/jurgen-kluft/hass-go/suncalc"
	"github.com/jurgen-kluft/hass-go/timeofday"
	"github.com/jurgen-kluft/hass-go/weather"
)

type scheduler struct {
	wait   map[string]time.Time
	states *state.Domain
}

type processfn func(states *state.Domain) time.Duration

func NewScheduler(states *state.Domain) *scheduler {
	s := &scheduler{}
	s.wait = map[string]time.Time{}
	s.states = states
	return s
}

func (s *scheduler) process(tag string, fn processfn) int {
	now := time.Now()
	updateAt, exists := s.wait[tag]
	if !exists {
		updateAt = now
		s.wait[tag] = updateAt
	}

	if now.Unix() >= updateAt.Unix() {
		wait := fn(s.states)
		now = time.Now()
		fmt.Printf("Scheduler updated '%s'\n", tag)
		s.wait[tag] = now.Add(wait)
		return 1
	}
	return 0
}

// sleep will figure out the earliest time a certain process needs
// to be updated and it will take that duration and sleep.
// Note: The maximum sleep time is 1 minute.
func (s *scheduler) sleep() {
	now := time.Now()
	earliest := now.Add(1 * time.Minute)
	for _, t := range s.wait {
		if t.Unix() < earliest.Unix() {
			earliest = t
		}
	}
	wait := time.Duration(earliest.Unix()-now.Unix()) * time.Second
	time.Sleep(wait)
}

func main() {

	// Create:
	states := state.New()

	shoutInstance, _ := shout.New()
	calendarInstance, _ := calendar.New()
	timeofdayInstance, _ := timeofday.New()
	weatherInstance, _ := weather.New()
	aqiInstance, _ := aqi.New()
	suncalcInstance, _ := suncalc.New()
	sensorsInstance, _ := sensors.New()
	lightingInstance, _ := lighting.New()
	hassInstance, _ := hass.New()
	reporterInstance, _ := reporter.New()

	scheduler := NewScheduler(states)
	for true {
		now := time.Now()
		states.SetTimeState("time", "now", now)

		fmt.Println("----- UPDATE -------")

		// Process
		updated := 0
		updated += scheduler.process("calendar", func(states *state.Domain) time.Duration { return calendarInstance.Process(states) })
		updated += scheduler.process("timeofday", func(states *state.Domain) time.Duration { return timeofdayInstance.Process(states) })
		updated += scheduler.process("suncalc", func(states *state.Domain) time.Duration { return suncalcInstance.Process(states) })
		updated += scheduler.process("weather", func(states *state.Domain) time.Duration { return weatherInstance.Process(states) })
		updated += scheduler.process("aqi", func(states *state.Domain) time.Duration { return aqiInstance.Process(states) })
		updated += scheduler.process("lighting", func(states *state.Domain) time.Duration { return lightingInstance.Process(states) })
		updated += scheduler.process("sensors", func(states *state.Domain) time.Duration { return sensorsInstance.Process(states) })
		updated += scheduler.process("hass", func(states *state.Domain) time.Duration { return hassInstance.Process(states) })
		updated += scheduler.process("reporter", func(states *state.Domain) time.Duration { return reporterInstance.Process(states) })
		updated += scheduler.process("shout", func(states *state.Domain) time.Duration { return shoutInstance.Process(states) })

		if updated > 0 {
			states.PrintNamed("time")
			states.PrintNamed("hass")
			fmt.Println("")
		}

		states.ResetChangeTracking()
		scheduler.sleep()
	}
}
