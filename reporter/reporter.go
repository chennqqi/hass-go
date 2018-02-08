package reporter

import (
	"fmt"
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

func (c *Instance) reportWeather(ID string, states *state.Domain) {
	if !c.historyContainsID(ID) {
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

				c.history[ID] = time.Now()
				report += line
			} else {
				break
			}
			i++
		}

		// Temperature morning - noon - evening

		// Weather report to
		states.SetStringState("shout", "msg:weather", report)
	}
}

func (c *Instance) Process(states *state.Domain) time.Duration {

	// In calendar these are written as:
	//
	//     report:weather=

	reports := states.Get("report")
	for _, r := range reports.Strings {
		if r == "weather" {
			if reports.HasStringState(r + ".ID") {
				id := reports.GetStringState(r+".ID", "")
				c.reportWeather(id, states)
			}
		}
	}

	return 1 * time.Second
}
