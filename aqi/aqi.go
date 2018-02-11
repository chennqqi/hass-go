package aqi

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/jurgen-kluft/hass-go/state"
)

type Instance struct {
	aqi    *Caqi
	update time.Time
	period time.Duration
}

func (c *Instance) readConfig() (*Caqi, error) {
	jsonBytes, err := ioutil.ReadFile("config/aqi.json")
	if err != nil {
		return nil, fmt.Errorf("ERROR: failed to read aqi config ( %s )", err)
	}
	obj, err := unmarshalcaqi(jsonBytes)
	return obj, err
}

func New() (c *Instance, err error) {
	c = &Instance{}
	c.aqi, err = c.readConfig()
	if err != nil {
		fmt.Println(err.Error())
	}
	c.update = time.Now()
	c.period = time.Minute * 15
	return c, err
}

func (c *Instance) getResponse() (AQI float64, err error) {
	url := c.aqi.URL
	url = strings.Replace(url, "${CITY}", c.aqi.City, 1)
	url = strings.Replace(url, "${TOKEN}", c.aqi.Token, 1)
	if strings.HasPrefix(url, "http") {
		var resp *http.Response
		fmt.Printf("HTTP Get, '%s'\n", url)
		resp, err = http.Get(url)
		if err != nil {
			AQI = 99.0
			resp.Body.Close()
		} else {
			var body []byte
			body, err = ioutil.ReadAll(resp.Body)
			var caqi CaqiResponse
			caqi, err = unmarshalCaqiResponse(body)
			AQI = float64(caqi.Data.Aqi)
			if err != nil {
				fmt.Print(string(body))
			}
		}
	} else if strings.HasPrefix(url, "print") {
		fmt.Printf("HTTP Get, '%s'\n", url)
	}
	return
}

func (c *Instance) getAiqTagAndDescr(aiq float64) (level AqiLevel) {
	for _, l := range c.aqi.Levels {
		if aiq < l.LessThan {
			level = l
			return
		}
	}
	level = c.aqi.Levels[1]
	return
}

// Process will get the AQI and post it in "weather"
func (c *Instance) Process(states *state.Domain) time.Duration {
	now := states.GetTimeState("time", "now", time.Now())
	if now.Unix() >= c.update.Unix() {
		aqi, err := c.getResponse()
		weather := states.Get("weather")
		//weather.ResetChangeTracking()
		if err == nil {
			weather.SetFloatState("aqi", aqi)
			level := c.getAiqTagAndDescr(aqi)
			weather.SetStringState("aqi", level.Tag)
			weather.SetStringState("aqi.implications", level.Implications)
			weather.SetStringState("aqi.caution", level.Caution)
		} else {
			fmt.Println(err.Error())
		}
		c.update = now.Add(c.period)
	}
	return 5 * time.Minute
}
