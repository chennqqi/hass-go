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
	viper *viper.Viper
	subs  []EventSubscriber
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

func (c *Calendar) pushEvent(uuid string, title string, descr string, start time.Time, end time.Time) {
	for _, s := range c.subs {
		s.Handle(uuid, title, descr, start, end)
	}
}

func (c *Calendar) updateEvents(calendars []*ics.Calendar) error {
	when := time.Now()
	for _, cal := range calendars {
		// get events for time 'when'
		eventsForDay, errEvents := cal.GetEventsByDate(when)
		if errEvents != nil { // error -> error
			return errEvents
		}

		for _, e := range eventsForDay {
			title := e.GetSummary()
			descr := e.GetDescription()
			c.pushEvent(e.GenerateEventId(), title, descr, e.GetStart(), e.GetEnd())
		}
	}
	return nil
}

func (c *Calendar) download() (err error) {
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

	errEvents := c.updateEvents(cals)
	if errEvents != nil { // error -> error
		return errEvents
	}

	return nil
}

type EventSubscriber interface {
	Handle(UUID string, title string, description string, start time.Time, end time.Time)
}

func (c *Calendar) RegisterSubscriber(subscriber EventSubscriber) {
	c.subs = append(c.subs, subscriber)
}

// Process will update 'events' from the calendar
func (c *Calendar) Process() error {
	// Other general states
	now := time.Now()
	weekend := now.Weekday() == time.Saturday || now.Weekday() == time.Sunday
	varStart := now
	varEnd := now
	varStr1 := ""
	varStr2 := ""
	if weekend {
		day := now.Day()
		if now.Weekday() == time.Sunday {
			day--
		}
		varStart = time.Date(now.Year(), now.Month(), day, 0, 0, 0, 0, now.Location())
		varEnd = time.Date(now.Year(), now.Month(), day+1, 23, 59, 59, 0, now.Location())
		varStr1 = "var:weekend=true"
		varStr2 = "var:weekday=false"
	} else {
		varStart = time.Date(now.Year(), now.Month(), int(time.Monday), 0, 0, 0, 0, now.Location())
		varEnd = time.Date(now.Year(), now.Month(), int(time.Friday), 23, 59, 59, 0, now.Location())
		varStr1 = "var:weekend=false"
		varStr2 = "var:weekday=true"
	}
	c.pushEvent("calendar.weekend", varStr1, "Is it weekend?", varStart, varEnd)
	c.pushEvent("calendar.weekday", varStr2, "Is it a weekday?", varStart, varEnd)

	// Download calendar
	err := c.download()
	if err != nil {
		return err
	}

	return nil
}
