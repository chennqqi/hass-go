package main

import (
	"fmt"
	"strings"
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

func PrintCalEventsToConsole(events []calendar.CEvent) {
	for _, e := range events {
		var domain string
		var dname string
		var dstate string
		title := strings.Replace(e.Title, ":", " : ", 1)
		title = strings.Replace(title, "=", " = ", 1)
		fmt.Sscanf(title, "%s : %s = %s", &domain, &dname, &dstate)
		//fmt.Printf("Parsed: '%s' - '%s' - '%s'\n", domain, dname, dstate)

		fmt.Printf("CalEvent %s: %s = %s (%s)\n", domain, dname, dstate, e.Description)
	}
}
func UpdateCalEventsToState(events []calendar.CEvent, sensors *sensors.Sensors) {
	for _, e := range events {
		var domain string
		var dname string
		var dstate string
		title := strings.Replace(e.Title, ":", " : ", 1)
		title = strings.Replace(title, "=", " = ", 1)
		fmt.Sscanf(title, "%s : %s = %s", &domain, &dname, &dstate)
		//fmt.Printf("Parsed: '%s' - '%s' - '%s'\n", domain, dname, dstate)
		if domain == "sensor" {
			sensors.UpdateSensor(dname, dstate)
		}
	}
}

func HandleWeatherReport(report []weather.Report) {

	// Get hourly report, cut off the head of anything that is before 'from'
	// Trim the tail of anything that is beyond 'until'

	// Report is like this:
	//   today            : Rain = 50%, Temperature = min - max
	//   air quality      : 50-100, Moderate
	//   sunrise   ( 6- 8): Fog, Soft breeze, 10 Celcius (Cool)
	//   morning   ( 8-10): Light Drizzle, Soft breeze, 8 to 13 Celcius (Cool)
	//   morning   (10-12):
	//   noon      (12-14):
	//   afternoon (14-16):
	//   afternoon (16-18):
	//   evening   (18-20):
	//

}

func main() {

	// TODO: implement the main hass-go function
	time.LoadLocation("Asia/Shanghai")

	// Create:
	stateInstance := state.New()
	stateInstance.SetTimeState("Now", time.Now())

	calendarInstance, _ := calendar.New()
	// im,  := im.New()
	weatherInstance, _ := weather.New()
	suncalcInstance, _ := suncalc.New()
	sensorsInstance, _ := sensors.New()
	lightingInstance, _ := lighting.New()
	lightingState := state.New()
	sensorState := state.New()

	// Process
	calEvents, calerr := calendarInstance.Process()
	if calerr != nil {
		panic(calerr)
	}
	PrintCalEventsToConsole(calEvents)
	UpdateCalEventsToState(calEvents, sensorsInstance)

	suncalcInstance.Process(stateInstance)
	weatherReport := weatherInstance.Process()
	HandleWeatherReport(weatherReport)

	lightingInstance.Process(stateInstance)
	for k, v := range stateInstance.Strings {
		if strings.HasPrefix(k, "sensor:") {
			fmt.Printf("Lighting State: %s = %s\n", k, v)
		}
	}
	for k, v := range lightingState.Floats {
		if strings.HasPrefix(k, "sensor:") {
			fmt.Printf("Lighting State: %s = %f\n", k, v)
		}
	}

	sensorsInstance.PublishSensors(sensorState)
	for k, v := range sensorState.Strings {
		fmt.Printf("Publish Sensor: %s = %s\n", k, v)
	}
	for k, v := range sensorState.Floats {
		fmt.Printf("Publish Sensor: %s = %f\n", k, v)
	}

	stateInstance.Print()

}
