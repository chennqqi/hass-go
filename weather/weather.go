package weather

import (
	"fmt"
	"time"

	"github.com/adlio/darksky"
	"github.com/jurgen-kluft/hass-go/dynamic"
	"github.com/jurgen-kluft/hass-go/state"
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
	rain := dynamic.Dynamic{c.viper.Get("rain")}
	for _, e := range rain.ArrayIter() {
		min := e.Get("intensity_min").AsFloat64()
		max := e.Get("intensity_max").AsFloat64()
		if pi >= min && pi < max {
			fmt.Printf("Rain: %s\n", e.Get("name").AsString())
		}
	}
}

func (c *Client) CreateReport(from time.Time, until time.Time, d darksky.DataPoint) {

	// Get hourly report, cut off the head of anything that is before 'from'
	// Trim the tail of anything that is beyond 'until'

	// Report is like this:
	//   today            : Rain = 50%, Temperature = min - max
	//   air quality      : 50-100, Moderate
	//   sunrise   ( 6- 8): Fog, Soft breeze, 10 Celcius (Cool)
	//   morning   ( 8-10): Light Drizzle, Soft breeze, 8 to 13 Celcius (Cool)
	//   morning   (10-12):
	//   noon      (12-14):
	//   afternoon (14-16):
	//   afternoon (16-18):
	//   evening   (18-20):
	//

}

func (c *Client) Process(state *state.Instance) {
	loc := dynamic.Dynamic{Item: c.viper.Get("location")}
	forecast, err := c.darksky.GetForecast(loc.Get("latitude").AsString(), loc.Get("longitude").AsString(), c.darkargs)
	if err == nil {
		// username := "The Weather"
		// msg := forecast.Currently.Summary
		// pretext := "Details"
		// prebody := fmt.Sprintf("Rain: %d%%\n", int(forecast.Currently.PrecipProbability))
		// prebody += fmt.Sprintf("Clouds: %d%%\n", int(forecast.Currently.CloudCover*100.0))
		// prebody += fmt.Sprintf("Temperature: %d C\n", int(converFToC(forecast.Currently.Temperature+0.5)))

		state.SetFloatState("Clouds", forecast.Currently.CloudCover)

		//im.PostMessage(c.viper.GetString("slack.channel"), username, msg, pretext, prebody)
	}

}
