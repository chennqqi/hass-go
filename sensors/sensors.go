package sensors

import (
	//"net/http"
	//"net/url"
	"github.com/jurgen-kluft/hass-go/dynamic"
	"github.com/spf13/viper"
)

// Sensors API
// - SetRain(rain int)
// - SetWind(wind int)
// - SetClouds(clouds int)
// - SetTemperature(celcius int)

// - SetCalEvent(event string)

// - UpdateToHASS

type sensor struct {
	name   string
	update bool
	state  string
	states []string
}

type Sensors struct {
	viper   *viper.Viper
	sensors map[string]*sensor
}

func New() (*Sensors, error) {
	s := &Sensors{}
	s.viper = viper.New()
	s.sensors = map[string]*sensor{}

	// Viper command-line package
	s.viper.SetConfigName("hass-go-sensors")        // name of config file (without extension)
	s.viper.AddConfigPath("$HOME/.hass-go-sensors") // call multiple times to add many search paths
	s.viper.AddConfigPath(".")                      // optionally look for config in the working directory
	err := s.viper.ReadInConfig()                   // Find and read the config file
	if err != nil {                                 // Handle errors reading the config file
		return nil, err
	}

	sensors := dynamic.Dynamic{s.viper.Get("sensor")}
	for _, e := range sensors.ArrayIter() {
		o := &sensor{}
		o.name = e.Get("name").AsString()
		o.update = true
		o.state = e.Get("default").AsString()
		s.sensors[o.name] = o
	}

	return s, nil
}

func (s *Sensors) SetRain(rain int) {

}
func (s *Sensors) SetWind(rain int) {

}
func (s *Sensors) SetTemperature(rain int) {

}
func (s *Sensors) SetCalEvent(event string) {

}

func (s *Sensors) UpdateToHASS() {

	// req, err := http.NewRequest("POST", HassURL+v.Sensor, bytes.NewBuffer(sensorJSON))
	// req.Header.Set("Content-Type", "application/json")
	// client := &http.Client{}
	// _, err = client.Do(req)

}
