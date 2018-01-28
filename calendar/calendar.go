package calendar

import (
	"github.com/jurgen-kluft/go-icloud-calendar"
	"github.com/jurgen-kluft/hass-go/dynamic"
	"github.com/spf13/viper"
	"time"
)

type Calendar struct {
	viper *viper.Viper
	cals  []*icalendar.Calendar
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

	dcals := dynamic.Dynamic{Item: c.viper.Get("calendars")}
	for _, dc := range dcals.ArrayIter() {
		url := dc.Get("url").AsString()
		cal := icalendar.NewURLCalendar(url)
		c.cals = append(c.cals, cal)
	}

	return c, nil
}

func (c *Calendar) updateEvents(when time.Time) ([]CEvent, error) {
	events := []CEvent{}
	for _, cal := range c.cals {
		eventsForDay := cal.GetEventsByDate(when)

		for _, e := range eventsForDay {
			event := CEvent{}
			event.UUID = e.GenerateUUID()
			event.Title = e.Summary
			event.Description = e.Description
			event.Start = e.Start
			event.End = e.End
			events = append(events, event)
		}
	}
	return events, nil
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

type CEvent struct {
	UUID        string
	Title       string
	Description string
	Start       time.Time
	End         time.Time
}

// Process will update 'events' from the calendar
func (c *Calendar) Process() ([]CEvent, error) {
	now := time.Now()

	events := []CEvent{}
	// Download calendar
	err := c.load()
	if err != nil {
		return events, err
	}
	// Update events
	events, err = c.updateEvents(now)
	if err != nil {
		return events, err
	}

	// Other general states
	varStr1 := ""
	varStr2 := ""
	weekend, varStart, varEnd := weekOrWeekEndStartEnd(now)
	if weekend {
		varStr1 = "var:weekend=true"
		varStr2 = "var:weekday=false"
	} else {
		varStr1 = "var:weekend=false"
		varStr2 = "var:weekday=true"
	}
	event := CEvent{}
	event.UUID = "calendar.weekend"
	event.Title = varStr1
	event.Description = "Is it weekend?"
	event.Start = varStart
	event.End = varEnd
	events = append(events, event)

	event = CEvent{}
	event.UUID = "calendar.weekday"
	event.Title = varStr2
	event.Description = "Is it a weekday?"
	event.Start = varStart
	event.End = varEnd
	events = append(events, event)

	return events, err
}
