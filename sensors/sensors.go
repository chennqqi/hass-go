package sensors

import (
	"github.com/jurgen-kluft/hass-go/state"
	"strconv"

	"github.com/jurgen-kluft/hass-go/dynamic"
	"github.com/spf13/viper"
)

type sensorStateAsString struct {
	name           string
	typeof         string
	unit           string
	update         bool
	defaultState   string
	possibleStates []string
}
type sensorStateAsFloat struct {
	name         string
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
	s.viper.SetConfigName("hass-go-sensors")        // name of config file (without extension)
	s.viper.AddConfigPath("$HOME/.hass-go-sensors") // call multiple times to add many search paths
	s.viper.AddConfigPath(".")                      // optionally look for config in the working directory
	err := s.viper.ReadInConfig()                   // Find and read the config file
	if err != nil {                                 // Handle errors reading the config file
		return nil, err
	}

	sensors := dynamic.Dynamic{Item: s.viper.Get("sensor")}
	for _, e := range sensors.ArrayIter() {
		typeof := e.Get("typeof").AsString()
		if typeof == "string" {
			o := &sensorStateAsString{}
			o.name = e.Get("name").AsString()
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
			o.name = e.Get("name").AsString()
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

func (s *Sensors) getStateString(sensorName string) string {
	state, exists := s.sstate[sensorName]
	if !exists {
		var sensor *sensorStateAsString
		sensor, _ = s.ssensors[sensorName]
		s.sstate[sensorName] = sensor.defaultState
		state = sensor.defaultState
	}
	return state
}
func (s *Sensors) setStateString(sensorName string, sensorState string) {
	s.sstate[sensorName] = sensorState
}

func (s *Sensors) getStateFloat(sensorName string) float64 {
	state, exists := s.fstate[sensorName]
	if !exists {
		var sensor *sensorStateAsFloat
		sensor, _ = s.fsensors[sensorName]
		s.fstate[sensorName] = sensor.min
		state = sensor.min
	}
	return state
}
func (s *Sensors) setStateFloat(sensorName string, sensorState float64) {
	s.fstate[sensorName] = sensorState
}

// UpdateSensor updates that state of sensor 'name' with value 'state'
func (s *Sensors) UpdateSensor(name string, state string) {
	stringSensor := s.getStringSensor(name)
	if stringSensor != nil {
		s.setStateString(name, state)
	} else {
		floatSensor := s.getFloatSensor(name)
		if floatSensor != nil {
			float, err := strconv.ParseFloat(state, 64)
			if err == nil {
				s.setStateFloat(name, float)
			}
		}
	}
}

// PublishSensors will write out the sensors to 'out'
func (s *Sensors) PublishSensors(out *state.Instance) {

	for name, value := range s.sstate {
		out.SetStringState("publish:"+name, value)
	}
	for name, value := range s.fstate {
		out.SetFloatState("publish:"+name, value)
	}
}
