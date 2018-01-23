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

func main() {

	// TODO: implement the main hass-go function
	time.LoadLocation("Asia/Shanghai")

	// Create:
	state := state.New()
	calendar, _ := calendar.New()
	// im,  := im.New()
	weather, _ := weather.New()
	suncalc, _ := suncalc.New()
	sensors, _ := sensors.New()
	lighting, _ := lighting.New(state)

	// Process
	calendar.Process(state)
	suncalc.Process(state)
	weather.Process(state)
	lighting.Process(state)
	sensors.Process(state)

	state.Print()
}
