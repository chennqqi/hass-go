package calendar

import (
	"errors"
	"fmt"
	"time"

	"github.com/PuloV/ics-golang"
	"github.com/jurgen-kluft/hass-go/dynamic"
	"github.com/spf13/viper"
)

type Calendar struct {
	viper  *viper.Viper
	events map[string]string
}

func New() (*Calendar, error) {
	c := &Calendar{}

	c.viper = viper.New()

	// Viper command-line package
	c.viper.SetConfigName("hass-go-calendar")        // name of config file (without extension)
	c.viper.AddConfigPath("$HOME/.hass-go-calendar") // call multiple times to add many search paths
	c.viper.AddConfigPath(".")                       // optionally look for config in the working directory
	err := c.viper.ReadInConfig()                    // Find and read the config file
	if err != nil {                                  // Handle errors reading the config file
		return nil, err
	}
	events := dynamic.Dynamic{Item: c.viper.Get("event")}
	for _, event := range events.ArrayIter() {
		ename := event.Get("name").AsString()
		estate := event.Get("state").AsString()
		c.events[ename] = estate
	}

	return c, nil
}

func updateEvents(calendars []*ics.Calendar, events map[string]string) error {
	when := time.Now()
	for _, cal := range calendars {
		// get events for time 'when'
		eventsForDay, errEvents := cal.GetEventsByDate(when)
		if errEvents != nil { // error -> error
			return errEvents
		}

		for _, e := range eventsForDay {
			title := e.GetSummary()
			var name, state string
			n, e := fmt.Sscanf(title, "%s.%s", &name, &state)
			if n == 2 && e == nil {
				events[name] = state
			}
		}
	}
	return nil
}

func (c *Calendar) download() (events map[string]string, err error) {
	parser := ics.New()
	input := parser.GetInputChan()
	dcals := dynamic.Dynamic{Item: c.viper.Get("calendar")}
	for _, c := range dcals.ArrayIter() {
		input <- c.Get("url").AsString()
	}
	parser.Wait()

	// 10 X 13 + 6 = 11 * 12 + 4

	// get all calendars from parser
	cals, errCals := parser.GetCalendars()

	events = map[string]string{}

	// if error or no calendars, error
	if errCals != nil {
		return events, errCals
	} else if len(cals) == 0 {
		return events, errors.New("No calendars (need one)")
	}

	// get events for time 'when' (using first calendar)
	errEvents := updateEvents(cals, events)
	if errEvents != nil { // error -> error
		return events, errEvents
	}

	return events, nil
}

// Process will update 'events' from the calendar
func (c *Calendar) Process(events *map[string]string) error {
	calendarEvents, err := c.download()
	if err != nil {
		return err
	}

	// First set all states we are tracking to their default
	// because not every event might occur in the calendar at
	// this specific date/time.
	for k, v := range c.events {
		(*events)[k] = v
	}

	// Then update all the states from the calendar events
	for k, v := range calendarEvents {
		(*events)[k] = v
	}
	return nil
}
