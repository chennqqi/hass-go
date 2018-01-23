package lighting

import (
	"fmt"
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
	ct          float64
	bri         float64
	darkorlight string
	startMoment string
	endMoment   string
}

func (l lighttime) Print() {
	fmt.Printf("lighttime.ct  : %f\n", l.ct)
	fmt.Printf("lighttime.bri : %f\n", l.bri)
	fmt.Printf("lighttime.dol : %s\n", l.darkorlight)
	fmt.Printf("lighttime.start : %s\n", l.startMoment)
	fmt.Printf("lighttime.end   : %s\n", l.endMoment)
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

type Instance struct {
	viper      *viper.Viper
	weather    []weathermod
	addtimes   []addtime
	lighttable []lighttime
	lighttypes []lighttype
}

func New(state *state.Instance) (*Instance, error) {
	l := &Instance{}
	l.viper = viper.New()
	l.weather = []weathermod{}
	l.addtimes = []addtime{}
	l.lighttable = []lighttime{}
	l.lighttypes = []lighttype{}

	// Viper command-line package
	l.viper.SetConfigName("hass-go-lighting")        // name of config file (without extension)
	l.viper.AddConfigPath("$HOME/.hass-go-lighting") // call multiple times to add many search paths
	l.viper.AddConfigPath(".")                       // optionally look for config in the working directory
	err := l.viper.ReadInConfig()                    // Find and read the config file
	if err != nil {                                  // Handle errors reading the config file
		return nil, err
	}

	dseason := dynamic.Dynamic{Item: l.viper.Get("season")}
	state.SetFloatState("Season.Winter", dseason.Get("winter").AsFloat64())
	state.SetFloatState("Season.Spring", dseason.Get("spring").AsFloat64())
	state.SetFloatState("Season.Summer", dseason.Get("summer").AsFloat64())
	state.SetFloatState("Season.Autumn", dseason.Get("autumn").AsFloat64())

	l.weather = []weathermod{}
	dweather := dynamic.Dynamic{Item: l.viper.Get("weather")}
	for _, dw := range dweather.ArrayIter() {
		w := weathermod{}
		w.clouds = dw.Get("clouds").AsFloat64()
		w.bri = dw.Get("bri").AsFloat64()
		w.ct = dw.Get("ct").AsFloat64()

		l.weather = append(l.weather, w)
	}

	daddtimes := dynamic.Dynamic{Item: l.viper.Get("addtime")}
	for _, dt := range daddtimes.ArrayIter() {
		t := addtime{}
		t.name = dt.Get("name").AsString()
		t.shift = dt.Get("shift").AsDuration()
		t.tag = dt.Get("tag").AsString()

		l.addtimes = append(l.addtimes, t)
	}

	dlighttypes := dynamic.Dynamic{Item: l.viper.Get("lighttype")}
	for _, dlt := range dlighttypes.ArrayIter() {
		lt := lighttype{}
		lt.name = dlt.Get("name").AsString()
		lt.minCT = dlt.Get("minCT").AsFloat64()
		lt.maxCT = dlt.Get("maxCT").AsFloat64()
		lt.minBRI = dlt.Get("minBRI").AsFloat64()
		lt.maxBRI = dlt.Get("maxBRI").AsFloat64()

		l.lighttypes = append(l.lighttypes, lt)
	}

	dlighttimes := dynamic.Dynamic{Item: l.viper.Get("lighttime")}
	for _, dlt := range dlighttimes.ArrayIter() {
		//kelvin      = 0.0
		//brightness  = 0.01
		//startMoment = "night:darkest:end"
		//endMoment   = "astronomical:dawn:begin"
		lt := lighttime{}
		lt.ct = dlt.Get("ct").AsFloat64()
		lt.bri = dlt.Get("bri").AsFloat64()
		lt.darkorlight = dlt.Get("darkorlight").AsString()
		lt.startMoment = dlt.Get("startMoment").AsString()
		lt.endMoment = dlt.Get("endMoment").AsString()
		//lt.Print()
		l.lighttable = append(l.lighttable, lt)
	}
	return l, nil
}

func inTimeSpan(start, end, check time.Time) bool {
	return check.After(start) && check.Before(end)
}

// Process will update 'string'states and 'float'states
// States are both input and output, for example as input
// there are Season/Weather states like 'Season':'Winter'
// and 'Clouds':0.5
func (l *Instance) Process(state *state.Instance) {
	now := time.Now()

	// Add our custom time-points
	for _, at := range l.addtimes {
		if state.HasTimeState(at.name) {
			t := state.GetTimeState(at.name, now)
			t = t.Add(at.shift)
			state.SetTimeState(at.name+at.tag, t)
		}
	}
	current := lighttime{}
	for _, lt := range l.lighttable {
		t0 := state.GetTimeState(lt.startMoment, now)
		t1 := state.GetTimeState(lt.endMoment, now)
		if inTimeSpan(t0, t1, now) {
			current = lt
			fmt.Printf("Current lighttime: %s -> %s\n\n", current.startMoment, current.endMoment)
			break
		}
	}

	season := state.GetStringState("Season", "Winter")
	seasonFac := state.GetFloatState("Season."+season, 0.0)
	cloudFac := state.GetFloatState("Clouds", 0.0)

	// Full cloud cover will increase color-temperature by 10% of (Max - Current)
	// NOTE: Only during the day (twilight + light)
	// TODO: when the moon is shining in the night the amount
	//       of blue-light is also higher than normal.
	// CT = 0 -> Coldest (>6500K)
	// CT = 1 -> Warmest (2000K)
	CT := (current.ct * seasonFac)
	if current.darkorlight != "dark" {
		// A bit colder color temperature when there are clouds during the day.
		CT = CT - cloudFac*0.1*CT
	}

	// Full cloud cover will increase brightness by 10% of (Max - Current)
	// BRI = 0 -> Very dim light
	// BRI = 1 -> Very bright light
	BRI := (current.bri * seasonFac)
	BRI = BRI + cloudFac*0.1*(1.0-BRI)
	state.SetFloatState("lights_BRI", BRI)

	// Update the state of the following string states
	// - lights_HUE_CT
	// - lights_HUE_BRI
	// - lights_YEE_CT
	// - lights_YEE_BRI
	for _, ltype := range l.lighttypes {
		lct := ltype.minCT + CT*(ltype.maxCT-ltype.minCT)
		state.SetFloatState("lights_"+ltype.name+"_CT", lct)
		lbri := ltype.minBRI + BRI*(ltype.maxBRI-ltype.minBRI)
		state.SetFloatState("lights_"+ltype.name+"_BRI", lbri)
	}

	// DOL = dark or light
	// States: 'dark', 'twilight', 'light'
	state.SetStringState("lights_DOL", current.darkorlight)

}
