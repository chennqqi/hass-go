package hass

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/jurgen-kluft/hass-go/dynamic"
	"github.com/jurgen-kluft/hass-go/state"
	"github.com/spf13/viper"
)

func postHttpSensor(url string, body string) (err error) {
	if strings.HasPrefix(url, "http") {
		var resp *http.Response
		resp, err = http.Post(url, "application/json", bytes.NewBufferString(body))
		if resp != nil {
			resp.Body.Close()
		}
	} else if strings.HasPrefix(url, "print") {
		fmt.Printf("HTTP Sensor, '%s', with message '%s'\n", url, body)
	}
	return err
}

type Instance struct {
	viper   *viper.Viper
	url     string
	vars    map[string]string
	body    string
	sensors []string
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

	sensor := dynamic.Dynamic{Item: s.viper.Get("hass")}
	s.url = sensor.Get("url").AsString()
	s.body = sensor.Get("body").AsString()
	vars := sensor.Get("vars")
	s.vars = map[string]string{}
	for _, v := range vars.ArrayIter() {
		kv := strings.Split(v.AsString(), "=")
		s.vars[kv[0]] = kv[1]
	}

	sensors := sensor.Get("sensors")
	s.sensors = []string{}
	for _, v := range sensors.ArrayIter() {
		s.sensors = append(s.sensors, v.AsString())
	}

	return s, nil
}

func (c *Instance) Process(states *state.Instance) time.Duration {

	for _, sn := range c.sensors {
		if states.HasStringState(sn) {
			surl := c.url
			sbody := c.body
			for vk, vv := range c.vars {
				vval := states.GetStringState(sn+"."+vk, "?") // Sensor value in string format
				surl = strings.Replace(surl, vv, vval, 1)
				sbody = strings.Replace(sbody, vv, vval, 1)
			}
			err := postHttpSensor(surl, sbody)
			if err != nil {
				break
			}
		}
	}
	return 30 * time.Second
}
