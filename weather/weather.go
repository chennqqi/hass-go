package weather

import (
	"fmt"
	"time"

	"github.com/adlio/darksky"
	"github.com/jurgen-kluft/hass-go/dynamic"
	"github.com/spf13/viper"
)

func converFToC(fahrenheit float64) float64 {
	return ((fahrenheit - 32.0) * 5.0 / 9.0)
}

type Client struct {
	viper    *viper.Viper
	location *time.Location
	darksky  *darksky.Client
	darkargs map[string]string
}

func New() (*Client, error) {
	c := &Client{}
	c.viper = viper.New()

	// Viper command-line package
	c.viper.SetConfigName("hass-go-weather")        // name of config file (without extension)
	c.viper.AddConfigPath("$HOME/.hass-go-weather") // call multiple times to add many search paths
	c.viper.AddConfigPath(".")                      // optionally look for config in the working directory
	err := c.viper.ReadInConfig()                   // Find and read the config file
	if err != nil {                                 // Handle errors reading the config file
		return nil, err
	}

	c.location, _ = time.LoadLocation(c.viper.GetString("location.timezone"))
	c.darksky = darksky.NewClient(c.viper.GetString("darksky.key"))
	c.darkargs = map[string]string{}
	c.darkargs["units"] = "si"

	return c, nil
}

func (c *Client) DetermineRain(d darksky.DataPoint) {

	pi := d.PrecipIntensity
	rain := dynamic.Dynamic{Item: c.viper.Get("rain")}
	for _, e := range rain.ArrayIter() {
		min := e.Get("intensity_min").AsFloat64()
		max := e.Get("intensity_max").AsFloat64()
		if pi >= min && pi < max {
			fmt.Printf("Rain: %s\n", e.Get("name").AsString())
		}
	}
}

type Report struct {
	From        time.Time
	Until       time.Time
	Rain        float64
	Clouds      float64
	Temperature float64
}

const (
	daySeconds = 60.0 * 60.0 * 24.0
)

func timeLater(date time.Time, t float64) time.Time {
	return time.Unix(date.Unix()+int64(t*float64(daySeconds)/24.0), 0)
}

func (c *Client) Process() []Report {
	loc := dynamic.Dynamic{Item: c.viper.Get("location")}
	forecast, err := c.darksky.GetForecast(loc.Get("latitude").AsString(), loc.Get("longitude").AsString(), c.darkargs)
	report := []Report{}
	if err == nil {
		from := time.Now()
		until := timeLater(from, 1.0)
		w := Report{}
		w.From = from
		w.Until = until
		w.Rain = forecast.Currently.PrecipProbability
		w.Clouds = forecast.Currently.CloudCover
		w.Temperature = forecast.Currently.ApparentTemperature
		report = append(report, w)
	}
	return report
}
