package calendar

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"github.com/jurgen-kluft/go-icloud-calendar"
	"github.com/jurgen-kluft/hass-go/state"
)

type Calendar struct {
	ccal   *ccalendar
	events map[string]cevent
	cals   []*icalendar.Calendar
}

func (c *Calendar) readConfig() (*ccalendar, error) {
	jsonBytes, err := ioutil.ReadFile("config/calendar.json")
	if err != nil {
		return nil, fmt.Errorf("Failed to read calendar config ( %s )", err)
	}
	ccal, err := unmarshalccalendar(jsonBytes)
	return ccal, err
}

func New() (*Calendar, error) {
	var err error

	c := &Calendar{}
	c.events = map[string]cevent{}
	c.ccal, err = c.readConfig()
	for _, cal := range c.ccal.event {
		c.events[cal.name] = cal
	}

	return c, err
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
