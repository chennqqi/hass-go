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
	viper     *viper.Viper
	calendars []*ics.Calendar
	events    map[string]string
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

	return c, nil
}

func (c *Calendar) update(events map[string]string) error {
	when := time.Now()
	for _, cal := range c.calendars {
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

func (c *Calendar) download() error {
	parser := ics.New()
	input := parser.GetInputChan()
	dcals := dynamic.Dynamic{Item: c.viper.Get("calendar")}
	for _, c := range dcals.ArrayIter() {
		input <- c.Get("url").AsString()
	}
	parser.Wait()

	// get all calendars from parser
	cals, errCals := parser.GetCalendars()
	c.calendars = cals

	// if error or no calendars, error
	if errCals != nil {
		return errCals
	} else if len(cals) == 0 {
		return errors.New("No calendars (need one)")
	}

	// get events for time 'when' (using first calendar)
	errEvents := c.update(c.events)
	if errEvents != nil { // error -> error
		return errEvents
	}

	return nil
}

func (c *Calendar) Process(events map[string]string) error {
	err := c.download()
	if err != nil {
		return err
	}
	for k, v := range c.events {
		events[k] = v
	}
	return nil
}
