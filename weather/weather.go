package weather

import (
	"fmt"
	"time"

	"github.com/adlio/darksky"
	"github.com/jurgen-kluft/hass-go/dynamic"
	"github.com/jurgen-kluft/hass-go/im"
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

func (c *Client) Process(im *im.IM) {
	forecast, err := c.darksky.GetForecast(c.viper.GetString("location.latitude"), c.viper.GetString("location.longitude"), c.darkargs)
	if err == nil {
		username := "The Weather"
		msg := forecast.Currently.Summary
		pretext := "Details"
		prebody := fmt.Sprintf("Rain: %d%%\n", int(forecast.Currently.PrecipProbability))
		prebody += fmt.Sprintf("Clouds: %d%%\n", int(forecast.Currently.CloudCover*100.0))
		prebody += fmt.Sprintf("Temperature: %d C\n", int(converFToC(forecast.Currently.Temperature+0.5)))

		im.PostMessage(c.viper.GetString("slack.channel"), username, msg, pretext, prebody)
	}

}
