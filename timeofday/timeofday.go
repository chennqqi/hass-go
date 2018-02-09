package timeofday

import (
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/jurgen-kluft/hass-go/state"
)

type Instance struct {
	tod *Ctimeofday
}

func (c *Instance) readConfig() (*Ctimeofday, error) {
	jsonBytes, err := ioutil.ReadFile("config/timeofday.json")
	if err != nil {
		return nil, fmt.Errorf("ERROR: failed to read timeofday config ( %s )", err)
	}
	ctod, err := unmarshalctimeofday(jsonBytes)
	return ctod, err
}

func New() (c *Instance, err error) {
	c = &Instance{}
	c.tod, err = c.readConfig()
	if err != nil {
		fmt.Println(err.Error())
	}
	//c.tod.print()
	return c, err
}

func isTimeofday(now time.Time, tod Ctime) (bool, int64) {
	t0 := int64(now.Hour())*3600 + int64(now.Minute())*60 + int64(now.Second())
	t1 := int64(tod.Hour())*3600 + int64(tod.Minute())*60 + int64(tod.Second())
	return t0 < t1, t1 - t0
}

func (c *Instance) Process(states *state.Domain) time.Duration {
	now := states.GetTimeState("time", "now", time.Now())

	weekday := strings.ToLower(now.Weekday().String())
	tods, exists := c.tod.Weekday[weekday]
	if !exists {
		tods, _ = c.tod.Weekday["anyday"]
	}

	wait := int64(10)
	for _, tod := range tods {
		var istod bool
		istod, wait = isTimeofday(now, tod.From)
		if istod {
			states.SetStringState("time", "tod", tod.Name)
			break
		}
	}
	return time.Duration(wait) * time.Second
}
