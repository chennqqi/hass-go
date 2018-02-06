package aqi

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/jurgen-kluft/hass-go/state"
)

type Instance struct {
	aqi *Caqi
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

	return c, err
}

func (c *Instance) getResponse() (AQI int64, err error) {
	url := c.aqi.URL
	if strings.HasPrefix(url, "http") {
		var resp *http.Response
		resp, err = http.Get(url)
		if resp != nil {
			AQI = 99
			resp.Body.Close()
		} else {
			var body []byte
			body, err = ioutil.ReadAll(resp.Body)
			var caqi CaqiResponse
			caqi, err = unmarshalCaqiResponse(body)
			AQI = caqi.Data.Aqi
		}
	} else if strings.HasPrefix(url, "print") {
		fmt.Printf("HTTP Get, '%s'\n", url)
	}
	return
}

func (c *Instance) Process(states *state.Domain) {

	return
}
