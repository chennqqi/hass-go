package calendar

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jurgen-kluft/go-icloud-calendar"
	"github.com/jurgen-kluft/hass-go/dynamic"
	"github.com/jurgen-kluft/hass-go/state"
	"github.com/spf13/viper"
)

type Calendar struct {
	viper  *viper.Viper
	events map[string]cevent
	cals   []*icalendar.Calendar
}

type cevent struct {
	calendar string
	domain   string
	name     string
	state    string
	typeof   string
	values   []string
}

func New() (*Calendar, error) {
	c := &Calendar{}
	c.viper = viper.New()
	c.events = map[string]cevent{}

	// Viper command-line package
	c.viper.SetConfigName("calendar") // name of config file (without extension)
	c.viper.AddConfigPath("config/")  // optionally look for config in the working directory
	err := c.viper.ReadInConfig()     // Find and read the config file
	if err != nil {                   // Handle errors reading the config file
		return nil, err
	}

	devents := dynamic.Dynamic{Item: c.viper.Get("events")}
	for _, de := range devents.ArrayIter() {
		e := cevent{}
		e.calendar = de.Get("calendar").AsString()
		e.domain = de.Get("domain").AsString()
		e.name = de.Get("name").AsString()
		e.state = de.Get("state").AsString()
		e.typeof = de.Get("typeof").AsString()
		e.values = []string{}
		for _, dv := range de.Get("values").ArrayIter() {
			e.values = append(e.values, dv.AsString())
		}
		ekey := e.domain + ":" + e.name
		c.events[ekey] = e
	}
	return c, nil
}

func (c *Calendar) updateEvents(when time.Time, states *state.Domain) error {
	for _, cal := range c.cals {
		eventsForDay := cal.GetEventsByDate(when)
		for _, e := range eventsForDay {
			var domain string
			var dname string
			var dstate string
			title := strings.Replace(e.Summary, ":", " : ", 1)
			title = strings.Replace(title, "=", " = ", 1)
			fmt.Sscanf(title, "%s : %s = %s", &domain, &dname, &dstate)
			//fmt.Printf("Parsed: '%s' - '%s' - '%s'\n", domain, dname, dstate)

			ekey := domain + ":" + dname
			ce, exists := c.events[ekey]
			if exists {
				if ce.typeof == "string" {
					states.SetStringState(domain, dname, dstate)
				} else if ce.typeof == "float" {
					fstate, err := strconv.ParseFloat(dstate, 64)
					if err == nil {
						states.SetFloatState(domain, dname, fstate)
					}
				}
			}
		}
	}
	return nil
}

func (c *Calendar) load() (err error) {
	for _, cal := range c.cals {
		err = cal.Load()
	}
	return err
}

func weekOrWeekEndStartEnd(now time.Time) (weekend bool, start, end time.Time) {
	day := now.Day()
	if now.Weekday() == time.Saturday || now.Weekday() == time.Sunday {
		if now.Weekday() == time.Sunday {
			day--
		}
		start = time.Date(now.Year(), now.Month(), day, 0, 0, 0, 0, now.Location())
		end = time.Date(now.Year(), now.Month(), day+1, 23, 59, 59, 0, now.Location())
		weekend = true
		return
	}

	weekend = false
	if now.Weekday() == time.Tuesday {
		day--
	} else if now.Weekday() == time.Wednesday {
		day -= 2
	} else if now.Weekday() == time.Thursday {
		day -= 3
	} else if now.Weekday() == time.Friday {
		day -= 4
	}
	start = time.Date(now.Year(), now.Month(), day, 0, 0, 0, 0, now.Location())
	end = time.Date(now.Year(), now.Month(), day+4, 23, 59, 59, 0, now.Location())

	return weekend, start, end
}

// Process will update 'events' from the calendar
func (c *Calendar) Process(states *state.Domain) error {
	now := states.GetTimeState("time", "now", time.Now())

	// Download calendar
	err := c.load()
	if err != nil {
		return err
	}

	// Default all states before updating them
	for _, eevent := range c.events {
		if eevent.typeof == "string" {
			states.SetStringState(eevent.domain, eevent.name, eevent.state)
		} else if eevent.typeof == "float" {
			fstate, err := strconv.ParseFloat(eevent.state, 64)
			if err == nil {
				states.SetFloatState(eevent.domain, eevent.name, fstate)
			}
		}
	}

	// Update events
	err = c.updateEvents(now, states)
	if err != nil {
		return err
	}

	// Other general states
	weekend, varStart, varEnd := weekOrWeekEndStartEnd(now)

	states.SetBoolState("calendar", "weekend", weekend)
	states.SetBoolState("calendar", "weekday", !weekend)

	states.SetTimeState("calendar", "weekend.start", varStart)
	states.SetTimeState("calendar", "weekend.end", varEnd)

	states.SetTimeState("calendar", "weekday.start", varStart)
	states.SetTimeState("calendar", "weekday.end", varEnd)

	states.SetStringState("calendar", "weekend.title", "Weekend")
	states.SetStringState("calendar", "weekday.title", "Weekday")

	states.SetStringState("calendar", "weekend.description", "Saturday and Sunday")
	states.SetStringState("calendar", "weekday.description", "Monday to Friday")

	return err
}
