package calendar

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/PuloV/ics-golang"
	"github.com/jurgen-kluft/hass-go/dynamic"
	"github.com/jurgen-kluft/hass-go/state"
	"github.com/spf13/viper"
)

type Calendar struct {
	viper   *viper.Viper
	sevents map[string]string
	fevents map[string]float64
}

func New() (*Calendar, error) {
	c := &Calendar{}
	c.sevents = map[string]string{}
	c.fevents = map[string]float64{}

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
		etype := event.Get("typeof").AsString()
		if etype == "string" {
			estate := event.Get("state").AsString()
			c.sevents[ename] = estate
		} else if etype == "float" {
			estate := event.Get("state").AsFloat64()
			c.fevents[ename] = estate
		}
	}

	return c, nil
}

func updateEvents(calendars []*ics.Calendar, state *state.Instance) error {
	when := time.Now()
	for _, cal := range calendars {
		// get events for time 'when'
		eventsForDay, errEvents := cal.GetEventsByDate(when)
		if errEvents != nil { // error -> error
			return errEvents
		}

		for _, e := range eventsForDay {
			title := e.GetSummary()
			//fmt.Printf("Calendar title: '%s'\n", title)
			e := strings.Split(title, ".")
			if len(e) == 2 {
				//fmt.Printf("Calendar event: %s-%s\n", e[0], e[1])
				if state.HasStringState(e[0]) {
					state.SetStringState(e[0], e[1])
				} else if state.HasFloatState(e[0]) {
					f, _ := strconv.ParseFloat(e[1], 64)
					state.SetFloatState(e[0], f)
				}
			}
		}
	}
	return nil
}

func (c *Calendar) download(state *state.Instance) (err error) {
	parser := ics.New()
	ics.DeleteTempFiles = false
	input := parser.GetInputChan()
	dcals := dynamic.Dynamic{Item: c.viper.Get("calendars")}
	for _, c := range dcals.ArrayIter() {
		calurl := c.Get("url").AsString()
		input <- string(calurl)
	}
	//input <- "tmp/test.ics"
	parser.Wait()
	cerrors, err := parser.GetErrors()
	if err == nil {
		for _, err := range cerrors {
			fmt.Printf("Calendar - ERROR: %s\n", err.Error())
		}
	}

	// get all calendars from parser
	cals, errCals := parser.GetCalendars()

	// if error or no calendars, error
	if errCals != nil {
		return errCals
	} else if len(cals) == 0 {
		return errors.New("No calendars (need one)")
	}

	// get events for time 'when' (using first calendar)
	errEvents := updateEvents(cals, state)
	if errEvents != nil { // error -> error
		return errEvents
	}

	return nil
}

// Process will update 'events' from the calendar
func (c *Calendar) Process(state *state.Instance) error {
	// First set all states we are tracking to their default
	// because not every event might occur in the calendar at
	// this specific date/time.
	for k, v := range c.sevents {
		state.SetStringState(k, v)
	}
	for k, v := range c.fevents {
		state.SetFloatState(k, v)
	}

	// Download calendar
	err := c.download(state)
	if err != nil {
		return err
	}

	return nil
}
