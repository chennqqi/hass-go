package hass

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"

	"github.com/jurgen-kluft/hass-go/dynamic"
	"github.com/jurgen-kluft/hass-go/state"
	"github.com/spf13/viper"
)

func postHttpSensor(url string, body string) {
	if strings.HasPrefix(url, "http") {
		resp, err := http.Post(url, "application/json", bytes.NewBufferString(body))
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
	} else if strings.HasPrefix(url, "print") {
		fmt.Printf("HTTP Sensor, '%s', with message '%s'", url, body)
	}
}

type Instance struct {
	viper *viper.Viper
	url   string
	vars  map[string]string
	body  string
}

// New will return a new instance of 'Sensors'
func New() (*Instance, error) {
	s := &Instance{}
	s.viper = viper.New()

	// Viper command-line package
	s.viper.SetConfigName("hass")    // name of config file (without extension)
	s.viper.AddConfigPath("config/") // optionally look for config in the working directory
	err := s.viper.ReadInConfig()    // Find and read the config file
	if err != nil {                  // Handle errors reading the config file
		return nil, err
	}

	sensors := dynamic.Dynamic{Item: s.viper.Get("sensor")}
	for _, e := range sensors.ArrayIter() {
		s.url = e.Get("url").AsString()
		s.body = e.Get("body").AsString()
		s.vars = map[string]string{}

		varname := ""
		for i, v := range e.Get("vars").ArrayIter() {
			if i&1 == 0 {
				varname = v.AsString()
			} else {
				s.vars[varname] = v.AsString()
			}
		}
	}

	return s, nil
}

func (c *Instance) Process(states *state.Domain) {
	sensors := states.Get("hass")
	for sn, _ := range sensors.Strings {
		surl := c.url
		sbody := c.body
		for vk, vv := range c.vars {
			vval := states.GetStringState("sensor", sn+"."+vk, "") // Sensor value in string format
			surl = strings.Replace(surl, vv, vval, 1)
			sbody = strings.Replace(sbody, vv, vval, 1)
		}
		postHttpSensor(surl, sbody)
	}
}
