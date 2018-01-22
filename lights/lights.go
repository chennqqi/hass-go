package lights

import (
	"time"

	"github.com/jurgen-kluft/hass-go/dynamic"
	"github.com/jurgen-kluft/hass-go/state"
	"github.com/spf13/viper"
)

// Color Temperature
// URL: https://panasonic.net/es/solution-works/jiyugaoka/

var defaultTime time.Time

type addtime struct {
	name  string
	tag   string
	shift time.Duration
}

type lighttime struct {
	kelvin      float64
	brightness  float64
	startMoment string
	endMoment   string
}

type lighttype struct {
	name   string
	minCT  float64
	maxCT  float64
	minBRI float64
	maxBRI float64
}

type weathermod struct {
	clouds float64 // Cloud factor
	ct     float64
	bri    float64
}

type Lights struct {
	viper      *viper.Viper
	weather    []weathermod
	addtimes   []addtime
	lighttable []lighttime
	lighttypes []lighttype
}

func New(state *state.Instance) (*Lights, error) {
	l := &Lights{}
	l.viper = viper.New()

	// Viper command-line package
	l.viper.SetConfigName("hass-go-lighting")        // name of config file (without extension)
	l.viper.AddConfigPath("$HOME/.hass-go-lighting") // call multiple times to add many search paths
	l.viper.AddConfigPath(".")                       // optionally look for config in the working directory
	err := l.viper.ReadInConfig()                    // Find and read the config file
	if err != nil {                                  // Handle errors reading the config file
		return nil, err
	}

	l.weather = []weathermod{}
	dweather := dynamic.Dynamic{Item: viper.Get("weather")}
	for _, dw := range dweather.ArrayIter() {
		w := weathermod{}
		w.clouds = dw.Get("clouds").AsFloat64()
		w.bri = dw.Get("bri").AsFloat64()
		w.ct = dw.Get("ct").AsFloat64()

		l.weather = append(l.weather, w)
	}

	daddtimes := dynamic.Dynamic{Item: viper.Get("addtime")}
	for _, dt := range daddtimes.ArrayIter() {
		t := addtime{}
		t.name = dt.Get("name").AsString()
		t.shift = dt.Get("shift").AsDuration()
		t.tag = dt.Get("tag").AsString()

		l.addtimes = append(l.addtimes, t)
	}

	dlighttypes := dynamic.Dynamic{Item: viper.Get("lighttype")}
	for _, dlt := range dlighttypes.ArrayIter() {
		lt := lighttype{}
		lt.name = dlt.Get("name").AsString()
		lt.minCT = dlt.Get("minCT").AsFloat64()
		lt.maxCT = dlt.Get("maxCT").AsFloat64()
		lt.minBRI = dlt.Get("minBRI").AsFloat64()
		lt.maxBRI = dlt.Get("maxBRI").AsFloat64()

		l.lighttypes = append(l.lighttypes, lt)
	}

	dlighttimes := dynamic.Dynamic{Item: viper.Get("lighttime")}
	for _, dlt := range dlighttimes.ArrayIter() {
		//kelvin      = 0.0
		//brightness  = 0.01
		//startMoment = "night:darkest:end"
		//endMoment   = "astronomical:dawn:begin"
		lt := lighttime{}
		lt.kelvin = dlt.Get("kelvin").AsFloat64()
		lt.brightness = dlt.Get("brightness").AsFloat64()
		lt.startMoment = dlt.Get("startMoment").AsString()
		lt.endMoment = dlt.Get("endMoment").AsString()

		l.lighttable = append(l.lighttable, lt)
	}
	return l, nil
}

// Process will update 'string'states and 'float'states
// States are both input and output, for example as input
// there are Season/Weather states like 'Season':'Winter'
// and 'Clouds':0.5
func (l *Lights) Process(state *state.Instance) {
	now := time.Now()

	for _, at := range l.addtimes {
		if state.HasTimeState(at.name) {
			t := state.GetTimeState(at.name, now)
			t = t.Add(at.shift)
			state.SetTimeState(at.name+at.tag, t)
		}
	}

	season := state.GetStringState("Season", "Winter")
	seasonMod := viper.GetFloat64("Season." + season)

	clouds := state.GetFloatState("Clouds", 0.5)

}
