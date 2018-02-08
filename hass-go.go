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

	for true {
		now := time.Now()
		states.SetTimeState("time", "now", now)

		fmt.Println("----- UPDATE -------")

		// Process
		calerr := calendarInstance.Process(states)
		if calerr != nil {
			fmt.Println("ERROR: Calendar error")
			//panic(calerr)
		}

		timeofdayInstance.Process(states)
		suncalcInstance.Process(states)
		weatherInstance.Process(states)
		aqiInstance.Process(states)
		lightingInstance.Process(states)
		sensorsInstance.PublishSensors(states)
		hassInstance.Process(states)
		reporterInstance.Process(states)

		shoutInstance.PublishMessages(states)

		states.PrintNamed("time")
		states.PrintNamed("hass")
		fmt.Println("")

		wait := now.Unix() + 2.0
		for true {
			t := time.Now().Unix()
			if t > wait {
				break
			}
		}
	}
}
