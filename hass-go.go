package main

import (
	"time"

	"github.com/jurgen-kluft/hass-go/calendar"
	"github.com/jurgen-kluft/hass-go/lighting"
	"github.com/jurgen-kluft/hass-go/sensors"
	"github.com/jurgen-kluft/hass-go/state"
	"github.com/jurgen-kluft/hass-go/suncalc"
	"github.com/jurgen-kluft/hass-go/weather"
)

// This can be computed from sun-rise, sun-set
// - sensor.dark_or_light (Dark, Twilight, Light)

// The data for these sensors can also come from a calendar
// In icloud we have the following calendars already:
// - TimeOfDay
// - Season
// - Weather
// - Lighting (NA)
// Note: These calendars can be downloaded on a low frequency.
// - Weather
// - TimeOfDay (Breakfast, Morning, Lunch, Afternoon, Evening, ..)
// - Season (Summer, Winter, Spring, Autumn)
// - Lighting: Color-Temperature and Brightness

// Every X minutes:
// - Update Calendar
// - Update SunCalc
// - Update DarkOrLight
// - Update Weather
// - Update Lighting
// - Update Sensors (this will HTTP to Home-Assistant)
// - Update Notifications and Reports
//   - Weather warnings
//   - Weather reports
//   - Alarms
// - Sleep

const (
	daySeconds = 60.0 * 60.0 * 24.0
)

func hoursLater(date time.Time, h float64) time.Time {
	return time.Unix(date.Unix()+int64(h*float64(daySeconds)/24.0), 0)
}

func buildWeatherReport(states *state.Domain) {
	report := "Weather Report"

	// Detect rain between
	//  -  8:30 - 9:30
	//  - 12:00 - 13:00
	//  - 18:00 - 20:00

	// Temperature morning - noon - evening

	// Weather report to
	states.SetStringState("slack", "weather", report)
}

func main() {

	// TODO: implement the main hass-go function
	//time.LoadLocation("Asia/Shanghai")
	now := time.Now()

	// Create:
	states := state.New()
	states.SetTimeState("time", "now", now)

	// im,  := im.New()
	calendarInstance, _ := calendar.New()
	weatherInstance, _ := weather.New()
	suncalcInstance, _ := suncalc.New()
	sensorsInstance, _ := sensors.New()
	lightingInstance, _ := lighting.New()

	// Process
	calerr := calendarInstance.Process(states)
	if calerr != nil {
		panic(calerr)
	}

	suncalcInstance.Process(states)
	weatherInstance.Process(states)
	lightingInstance.Process(states)
	sensorsInstance.PublishSensors(states)

	states.Print()

	states.Clear("publish")
}
