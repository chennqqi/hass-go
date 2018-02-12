package reporter

import (
	"fmt"
	"strings"
	"time"

	"github.com/jurgen-kluft/hass-go/state"
)

type Instance struct {
	history map[string]time.Time
}

func New() (c *Instance, err error) {
	c = &Instance{}
	return c, err
}

func (c *Instance) historyContainsID(ID string) bool {
	_, exists := c.history[ID]
	return exists
}

func (c *Instance) reportWeather(ID string, states *state.Instance) {
	if !c.historyContainsID(ID) {
		title := "Weather Report"

		// Detect rain between
		//  -  8:30 - 9:30
		//  - 12:00 - 13:00
		//  - 18:00 - 20:00
		report := title + "\n"
		report += "Change of rain is " + states.GetStringState("weather.currently:rain", "") + "\n"
		//fmt.Print(report)

		i := 1
		for true {
			key := fmt.Sprintf("weather.hourly[%d]:", i)
			if states.HasTimeState(key + "from") {
				hfrom := states.GetTimeState(key+"from", time.Now())
				huntil := states.GetTimeState(key+"until", time.Now())
				srain := states.GetStringState(key+"rain", "")
				scloud := states.GetStringState(key+"clouds", "")
				stemp := states.GetStringState(key+"temperature", "")
				temp := states.GetFloatState(key+"temperature", 0.0)
				line := fmt.Sprintf("%s, %s(%d), %s (%02d:%02d - %02d:%02d)\n", srain, stemp, int32(temp+0.5), scloud, hfrom.Hour(), hfrom.Minute(), huntil.Hour(), huntil.Minute())
				//fmt.Print(line)

				c.history[ID] = time.Now()
				report += line
			} else {
				break
			}
			i++
		}

		// Temperature morning - noon - evening

		// Weather report to
		states.SetStringState("shout.msg:weather", report)
	}
}

func (c *Instance) Process(states *state.Instance) time.Duration {
	for n := range states.Properties {
		if strings.HasPrefix(n, "reports.weather") {
			if states.HasStringState(n + ".ID") {
				id := states.GetStringState(n+".ID", "")
				c.reportWeather(id, states)
			}
		}
	}
	return 30 * time.Second
}
