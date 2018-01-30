package weather

import (
	"fmt"
	"time"

	"github.com/adlio/darksky"
	"github.com/jurgen-kluft/hass-go/state"
	"github.com/spf13/viper"
)

func converFToC(fahrenheit float64) float64 {
	return ((fahrenheit - 32.0) * 5.0 / 9.0)
}

type Client struct {
	viper     *viper.Viper
	location  *time.Location
	darksky   *darksky.Client
	latitude  float64
	longitude float64
	darkargs  map[string]string
}

func New() (*Client, error) {
	c := &Client{}
	c.viper = viper.New()

	// Viper command-line package
	c.viper.SetConfigName("weather") // name of config file (without extension)
	c.viper.AddConfigPath("config/") // optionally look for config in the working directory
	err := c.viper.ReadInConfig()    // Find and read the config file
	if err != nil {                  // Handle errors reading the config file
		return nil, err
	}

	c.location, _ = time.LoadLocation(c.viper.GetString("location.timezone"))
	c.darksky = darksky.NewClient(c.viper.GetString("darksky.key"))
	c.darkargs = map[string]string{}
	c.darkargs["units"] = "si"

	return c, nil
}

func (c *Client) updateHourly(from time.Time, until time.Time, states *state.Domain, hourly *darksky.DataBlock) {

	for i, dp := range hourly.Data {
		hfrom := time.Unix(dp.Time.Unix(), 0)
		huntil := hoursLater(hfrom, 1.0)
		if hfrom.After(from) && huntil.Before(until) {
			states.SetTimeState("weather", fmt.Sprintf("hourly[%d]:from", i), hfrom)
			states.SetTimeState("weather", fmt.Sprintf("hourly[%d]:until", i), huntil)
			states.SetFloatState("weather", fmt.Sprintf("hourly[%d]:rain", i), dp.PrecipProbability)
			states.SetFloatState("weather", fmt.Sprintf("hourly[%d]:clouds", i), dp.CloudCover)
			states.SetFloatState("weather", fmt.Sprintf("hourly[%d]:temperature", i), dp.ApparentTemperature)
		}
	}
}

const (
	daySeconds = 60.0 * 60.0 * 24.0
)

func hoursLater(date time.Time, h float64) time.Time {
	return time.Unix(date.Unix()+int64(h*float64(daySeconds)/24.0), 0)
}

func (c *Client) Process(states *state.Domain) {
	lat := states.GetFloatState("geo", "latitude", c.latitude)
	lng := states.GetFloatState("geo", "longitude", c.longitude)
	forecast, err := c.darksky.GetForecast(fmt.Sprint(lat), fmt.Sprint(lng), c.darkargs)
	if err == nil {
		now := states.GetTimeState("time", "now", time.Now())

		from := now
		until := hoursLater(from, 1.0)

		states.SetTimeState("weather", "currently:from", from)
		states.SetTimeState("weather", "currently:until", until)
		states.SetFloatState("weather", "currently:rain", forecast.Currently.PrecipProbability)
		states.SetFloatState("weather", "currently:clouds", forecast.Currently.CloudCover)
		states.SetFloatState("weather", "currently:temperature", forecast.Currently.ApparentTemperature)

		c.updateHourly(now, hoursLater(now, 12.0), states, forecast.Hourly)
	}
	return
}
