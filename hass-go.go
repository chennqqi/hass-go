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

type waiter struct {
	wait   time.Duration
	states *state.Domain
}

type process func(states *state.Domain) time.Duration

func (w *waiter) process(wait time.Duration) {
	if wait < w.wait {
		w.wait = wait
	}
}

func (w *waiter) sleep() {
	time.Sleep(w.wait)
	w.wait = 30 * time.Minute
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

	waiter := &waiter{}
	for true {
		now := time.Now()
		states.SetTimeState("time", "now", now)

		fmt.Println("----- UPDATE -------")

		// Process
		waiter.wait = 30 * time.Minute
		waiter.process(calendarInstance.Process(states))
		waiter.process(calendarInstance.Process(states))
		waiter.process(timeofdayInstance.Process(states))
		waiter.process(suncalcInstance.Process(states))
		waiter.process(weatherInstance.Process(states))
		waiter.process(aqiInstance.Process(states))
		waiter.process(lightingInstance.Process(states))
		waiter.process(sensorsInstance.Process(states))
		waiter.process(hassInstance.Process(states))
		waiter.process(reporterInstance.Process(states))
		waiter.process(shoutInstance.Process(states))

		states.PrintNamed("time")
		states.PrintNamed("hass")
		fmt.Println("")

		waiter.sleep()
	}
}
