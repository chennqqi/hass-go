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
	if strings.HasPrefix(url, "http") {
		var resp *http.Response
		resp, err = http.Get(url)
		if resp != nil {
			AQI = 99.0
			resp.Body.Close()
		} else {
			var body []byte
			body, err = ioutil.ReadAll(resp.Body)
			var caqi CaqiResponse
			caqi, err = unmarshalCaqiResponse(body)
			AQI = float64(caqi.Data.Aqi)
		}
	} else if strings.HasPrefix(url, "print") {
		fmt.Printf("HTTP Get, '%s'\n", url)
	}
	return
}

// TODO: This could go in the json configuration file
func getAiqTagAndDescr(aiq float64) (tag, implications, caution string) {
	if aiq < 50.0 {
		return "Good", "Air quality is considered satisfactory, and air pollution poses little or no risk", "None"
	} else if aiq < 100 {
		return "Moderate", "Air quality is acceptable; however, for some pollutants there may be a moderate health concern for a very small number of people who are unusually sensitive to air pollution.", "Active children and adults, and people with respiratory disease, such as asthma, should limit prolonged outdoor exertion."
	} else if aiq < 150 {
		return "Unhealthy for Sensitive Groups", "Members of sensitive groups may experience health effects. The general public is not likely to be affected.", "Active children and adults, and people with respiratory disease, such as asthma, should limit prolonged outdoor exertion."
	} else if aiq < 200 {
		return "Unhealthy", "Everyone may begin to experience health effects; members of sensitive groups may experience more serious health effects", "Active children and adults, and people with respiratory disease, such as asthma, should avoid prolonged outdoor exertion; everyone else, especially children, should limit prolonged outdoor exertion"
	} else if aiq < 300 {
		return "Very Unhealthy", "Health warnings of emergency conditions. The entire population is more likely to be affected.", "Active children and adults, and people with respiratory disease, such as asthma, should avoid all outdoor exertion; everyone else, especially children, should limit outdoor exertion."
	}
	return "Hazardous", "Health alert: everyone may experience more serious health effects.", "Everyone should avoid all outdoor exertion"
}

// Process will get the AQI and post it in "weather"
func (c *Instance) Process(states *state.Domain) {
	now := states.GetTimeState("time", "now", time.Now())
	if now.Unix() >= c.update.Unix() {
		aqi, err := c.getResponse()
		if err == nil {
			states.SetFloatState("weather", "aqi", aqi)
			tag, implications, caution := getAiqTagAndDescr(aqi)
			states.SetStringState("weather", "aqi", tag)
			states.SetStringState("weather", "aqi.implications", implications)
			states.SetStringState("weather", "aqi.caution", caution)
		}
		c.update = now.Add(c.period)
	}
	return
}
