package main

import (
	"fmt"
	"github.com/jurgen-kluft/hass-go/calendar"
	"github.com/jurgen-kluft/hass-go/hass"
	"github.com/jurgen-kluft/hass-go/lighting"
	"github.com/jurgen-kluft/hass-go/sensors"
	"github.com/jurgen-kluft/hass-go/state"
	"github.com/jurgen-kluft/hass-go/suncalc"
	"github.com/jurgen-kluft/hass-go/weather"
	"time"
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
	title := "Weather Report"

	// Detect rain between
	//  -  8:30 - 9:30
	//  - 12:00 - 13:00
	//  - 18:00 - 20:00
	weather := states.Get("weather")

	report := title + "\n"
	report += "Change of rain is " + weather.GetStringState("currently:rain", "") + "\n"
	//fmt.Print(report)

	i := 1
	for true {
		key := fmt.Sprintf("hourly[%d]:", i)
		if weather.HasTimeState(key + "from") {
			hfrom := weather.GetTimeState(key+"from", time.Now())
			huntil := weather.GetTimeState(key+"until", time.Now())
			srain := weather.GetStringState(key+"rain", "")
			scloud := weather.GetStringState(key+"clouds", "")
			stemp := weather.GetStringState(key+"temperature", "")
			temp := weather.GetFloatState(key+"temperature", 0.0)
			line := fmt.Sprintf("%s, %s(%d), %s (%02d:%02d - %02d:%02d)\n", srain, stemp, int32(temp+0.5), scloud, hfrom.Hour(), hfrom.Minute(), huntil.Hour(), huntil.Minute())
			//fmt.Print(line)
			report += line
		} else {
			break
		}
		i++
	}

	// Temperature morning - noon - evening

	// Weather report to
	states.SetStringState("slack", "weather", report)
}

func main() {

	// Create:
	states := state.New()

	// im,  := im.New()
	calendarInstance, _ := calendar.New()
	weatherInstance, _ := weather.New()
	suncalcInstance, _ := suncalc.New()
	sensorsInstance, _ := sensors.New()
	lightingInstance, _ := lighting.New()
	hassInstance, _ := hass.New()

	for true {
		now := time.Now()
		states.SetTimeState("time", "now", now)

		fmt.Println("----- UPDATE -------")
		states.PrintNamed("time")

		// Process
		calerr := calendarInstance.Process(states)
		if calerr != nil {
			panic(calerr)
		}

		suncalcInstance.Process(states)
		weatherInstance.Process(states)
		lightingInstance.Process(states)
		sensorsInstance.PublishSensors(states)
		hassInstance.Process(states)

		buildWeatherReport(states)

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
