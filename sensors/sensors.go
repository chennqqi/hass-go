package sensors

import (
	//"net/http"
	//"net/url"
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
}

// New will return a new instance of 'Sensors'
func New() (*Sensors, error) {
	s := &Sensors{}
	s.viper = viper.New()
	s.ssensors = map[string]*sensorStateAsString{}
	s.fsensors = map[string]*sensorStateAsFloat{}

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
			possibleStates := dynamic.Dynamic{Item: e.Get("states")}
			for _, state := range possibleStates.ArrayIter() {
				o.possibleStates = append(o.possibleStates, state.AsString())
			}

			s.ssensors[o.name] = o
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
		}
	}

	return s, nil
}

func (s *Sensors) Process(sstates *map[string]string, fstates [string]float64) {

	// req, err := http.NewRequest("POST", HassURL+v.Sensor, bytes.NewBuffer(sensorJSON))
	// req.Header.Set("Content-Type", "application/json")
	// client := &http.Client{}
	// _, err = client.Do(req)

}
