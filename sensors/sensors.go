package sensors

import (
	"fmt"

	"github.com/jurgen-kluft/hass-go/state"

	"github.com/jurgen-kluft/hass-go/dynamic"
	"github.com/spf13/viper"
)

type sensorStateAsString struct {
	domain         string
	name           string
	descr          string
	typeof         string
	unit           string
	update         bool
	defaultState   string
	possibleStates []string
}
type sensorStateAsFloat struct {
	domain       string
	name         string
	descr        string
	typeof       string
	unit         string
	update       bool
	defaultState float64
	min          float64
	max          float64
}

// Sensors is an instance to track sensor state
type Sensors struct {
	viper    *viper.Viper
	ssensors map[string]*sensorStateAsString
	fsensors map[string]*sensorStateAsFloat
	sstate   map[string]string
	fstate   map[string]float64
}

// New will return a new instance of 'Sensors'
func New() (*Sensors, error) {
	s := &Sensors{}
	s.viper = viper.New()
	s.ssensors = map[string]*sensorStateAsString{}
	s.fsensors = map[string]*sensorStateAsFloat{}
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
		typeof := e.Get("typeof").AsString()
		if typeof == "string" || typeof == "school" || typeof == "work" {
			o := &sensorStateAsString{}
			o.domain = e.Get("domain").AsString()
			o.name = e.Get("name").AsString()
			o.descr = e.Get("descr").AsString()
			o.typeof = typeof
			o.unit = e.Get("unit").AsString()
			o.update = true
			o.possibleStates = []string{}
			o.defaultState = e.Get("default").AsString()
			possibleStates := e.Get("states")
			for _, state := range possibleStates.ArrayIter() {
				o.possibleStates = append(o.possibleStates, state.AsString())
			}

			s.ssensors[o.name] = o
			s.sstate[o.name] = o.defaultState
		} else if typeof == "float" {
			o := &sensorStateAsFloat{}
			o.domain = e.Get("domain").AsString()
			o.name = e.Get("name").AsString()
			o.descr = e.Get("descr").AsString()
			o.typeof = typeof
			o.unit = e.Get("unit").AsString()
			o.update = true
			o.defaultState = e.Get("default").AsFloat64()
			o.min = e.Get("min").AsFloat64()
			o.max = e.Get("max").AsFloat64()

			s.fsensors[o.name] = o
			s.fstate[o.name] = o.defaultState
		}
	}

	return s, nil
}

func (s *Sensors) getStringSensor(sensorName string) *sensorStateAsString {
	sensor, exists := s.ssensors[sensorName]
	if !exists {
		return sensor
	}
	return nil
}
func (s *Sensors) getFloatSensor(sensorName string) *sensorStateAsFloat {
	sensor, exists := s.fsensors[sensorName]
	if !exists {
		return sensor
	}
	return nil
}

// PublishSensors will write out the sensors to 'out'
func (s *Sensors) PublishSensors(states *state.Domain) {
	for _, sensor := range s.ssensors {
		state := states.GetStringState(sensor.domain, sensor.name, sensor.defaultState)
		states.SetStringState("hass", sensor.name, state)

		states.SetStringState("sensor", sensor.name+".name", sensor.name)
		states.SetStringState("sensor", sensor.name+".typeof", sensor.typeof)
		states.SetStringState("sensor", sensor.name+".state", state)
		states.SetStringState("sensor", sensor.name+".descr", sensor.descr)
		states.SetStringState("sensor", sensor.name+".unit", sensor.unit)
	}
	for _, sensor := range s.fsensors {
		state := states.GetFloatState(sensor.domain, sensor.name, sensor.defaultState)
		states.SetFloatState("hass", sensor.name, state)

		states.SetStringState("sensor", sensor.name+".name", sensor.name)
		states.SetStringState("sensor", sensor.name+".typeof", sensor.typeof)
		states.SetStringState("sensor", sensor.name+".state", fmt.Sprintf("%f", state))
		states.SetStringState("sensor", sensor.name+".descr", sensor.descr)
		states.SetStringState("sensor", sensor.name+".unit", sensor.unit)
	}
}
