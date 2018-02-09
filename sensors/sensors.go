package sensors

import (
	"fmt"
	"time"

	"github.com/jurgen-kluft/hass-go/dynamic"
	"github.com/jurgen-kluft/hass-go/state"
	"github.com/spf13/viper"
)

type sensorStateAsString struct {
	domain string
	name   string
	descr  string
	unit   string
	update bool
}

// Sensors is an instance to track sensor state
type Sensors struct {
	viper    *viper.Viper
	ssensors map[string]*sensorStateAsString
	sstate   map[string]string
	fstate   map[string]float64
}

// New will return a new instance of 'Sensors'
func New() (*Sensors, error) {
	s := &Sensors{}
	s.viper = viper.New()
	s.ssensors = map[string]*sensorStateAsString{}
	s.sstate = map[string]string{}
	s.fstate = map[string]float64{}

	// Viper command-line package
	s.viper.SetConfigName("sensors") // name of config file (without extension)
	s.viper.AddConfigPath("config/") // optionally look for config in the working directory
	err := s.viper.ReadInConfig()    // Find and read the config file
	if err != nil {                  // Handle errors reading the config file
		return nil, err
	}

	sensors := dynamic.Dynamic{Item: s.viper.Get("sensor")}
	for _, e := range sensors.ArrayIter() {
		o := &sensorStateAsString{}
		o.domain = e.Get("domain").AsString()
		o.name = e.Get("name").AsString()
		o.descr = e.Get("descr").AsString()
		o.unit = e.Get("unit").AsString()

		s.ssensors[o.name] = o
	}

	return s, nil
}

// Process will write out the sensors to 'out'
func (s *Sensors) Process(states *state.Domain) time.Duration {
	hassState := states.Get("hass")
	sensorState := states.Get("sensor")
	if sensorState.HasChanged() {
		for _, sensor := range s.ssensors {
			if states.HasStringState(sensor.domain, sensor.name) {
				state := states.GetStringState(sensor.domain, sensor.name, "")
				hassState.SetStringState(sensor.name, state)
				sensorState.SetStringState(sensor.name+".name", sensor.name)
				sensorState.SetStringState(sensor.name+".state", state)
				sensorState.SetStringState(sensor.name+".descr", sensor.descr)
				sensorState.SetStringState(sensor.name+".unit", sensor.unit)
			} else if states.HasFloatState(sensor.domain, sensor.name) {
				state := states.GetFloatState(sensor.domain, sensor.name, 0.0)
				hassState.SetStringState(sensor.name, fmt.Sprintf("%.2f", state))
				sensorState.SetStringState(sensor.name+".name", sensor.name)
				sensorState.SetStringState(sensor.name+".state", fmt.Sprintf("%.2f", state))
				sensorState.SetStringState(sensor.name+".descr", sensor.descr)
				sensorState.SetStringState(sensor.name+".unit", sensor.unit)
			} else {
				fmt.Printf("ERROR: Sensor with name %s in domain %s doesn't exist\n", sensor.name, sensor.domain)
			}
		}
	}
	return 1 * time.Hour
}
