package lighting

import (
	"fmt"
	"math"
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
	ct          [2]float64
	bri         [2]float64
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

type seasonmod struct {
	name   string
	minCT  float64
	maxCT  float64
	minBRI float64
	maxBRI float64
}

type weathermod struct {
	clouds  float64 // Cloud factor
	ct_pct  float64 // Percentage change +/-
	bri_pct float64 // Percentage change +/-
}

type Instance struct {
	viper      *viper.Viper
	season     map[string]seasonmod
	weather    []weathermod
	addtimes   []addtime
	lighttable []lighttime
	lighttypes []lighttype
}

func New() (*Instance, error) {
	l := &Instance{}
	l.viper = viper.New()
	l.season = map[string]seasonmod{}
	l.weather = []weathermod{}
	l.addtimes = []addtime{}
	l.lighttable = []lighttime{}
	l.lighttypes = []lighttype{}

	// Viper command-line package
	l.viper.SetConfigName("lighting") // name of config file (without extension)
	l.viper.AddConfigPath("config/")  // optionally look for config in the working directory
	err := l.viper.ReadInConfig()     // Find and read the config file
	if err != nil {                   // Handle errors reading the config file
		return nil, err
	}

	dseasonmod := dynamic.Dynamic{Item: l.viper.Get("season")}
	for _, ds := range dseasonmod.ArrayIter() {
		sm := seasonmod{}
		sm.name = ds.Get("name").AsString()
		sm.minCT = ds.Get("minCT").AsFloat64()
		sm.maxCT = ds.Get("maxCT").AsFloat64()
		sm.minBRI = ds.Get("minBRI").AsFloat64()
		sm.maxBRI = ds.Get("maxBRI").AsFloat64()

		l.season[sm.name] = sm
	}

	l.weather = []weathermod{}
	dweather := dynamic.Dynamic{Item: l.viper.Get("weather")}
	for _, dw := range dweather.ArrayIter() {
		w := weathermod{}
		w.clouds = dw.Get("clouds").AsFloat64()
		w.bri_pct = dw.Get("bri_pct").AsFloat64()
		w.ct_pct = dw.Get("ct_pct").AsFloat64()

		l.weather = append(l.weather, w)
	}

	dextratimes := dynamic.Dynamic{Item: l.viper.Get("extra_suncalc_time")}
	for _, dt := range dextratimes.ArrayIter() {
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
		dcta := dlt.Get("ct")
		for i, v := range dcta.ArrayIter() {
			lt.ct[i] = v.AsFloat64()
		}
		cbria := dlt.Get("bri")
		for i, v := range cbria.ArrayIter() {
			lt.bri[i] = v.AsFloat64()
		}
		lt.darkorlight = dlt.Get("darkorlight").AsString()
		lt.startMoment = dlt.Get("startMoment").AsString()
		lt.endMoment = dlt.Get("endMoment").AsString()
		//lt.Print()
		l.lighttable = append(l.lighttable, lt)
	}
	return l, nil
}

func inTimeSpan(start, end, t time.Time) bool {
	return t.After(start) && t.Before(end)
}

// Return the factor 0.0 - 1.0 that indicates where we are in between start - end
func computeTimeSpanX(start, end, t time.Time) float64 {
	sh, sm, sc := start.Clock()
	sx := float64(sh*60*60) + float64(sm*60) + float64(sc)
	eh, em, ec := end.Clock()
	ex := float64(eh*60*60) + float64(em*60) + float64(ec)
	th, tm, tc := t.Clock()
	tx := float64(th*60*60) + float64(tm*60) + float64(tc)
	x := (tx - sx) / (ex - sx)
	return x
}

// Process will update 'string'states and 'float'states
// States are both input and output, for example as input
// there are Season/Weather states like 'Season':'Winter'
// and 'Clouds':0.5
func (l *Instance) Process(states *state.Instance) time.Duration {
	// Update our internal state with 'state'
	now := states.GetTimeState("time.now", time.Now())

	// Add our custom time-points
	for _, at := range l.addtimes {
		t := states.GetTimeState("sun."+at.name, now)
		t = t.Add(at.shift)
		states.SetTimeState("sun."+at.name+at.tag, t)
	}

	current := lighttime{}
	currentx := 0.0 // Time interpolation factor, where are we between startMoment - endMoment
	for _, lt := range l.lighttable {
		t0 := states.GetTimeState("sun."+lt.startMoment, now)
		t1 := states.GetTimeState("sun."+lt.endMoment, now)
		if inTimeSpan(t0, t1, now) {
			current = lt
			currentx = computeTimeSpanX(t0, t1, now)
			currentx = float64(int64(currentx*100.0)) / 100.0
			//fmt.Printf("Current lighttime: %s -> %s (x: %f)\n\n", current.startMoment, current.endMoment, currentx)
			states.SetStringState("lighting.current", fmt.Sprintf("%s -> %s (x: %f)", current.startMoment, current.endMoment, currentx))
			break
		}
	}

	seasonName := states.GetStringState("time.season", "winter")
	season := l.season[seasonName]
	clouds := weathermod{clouds: 0.0, ct_pct: 0.0, bri_pct: 0.0}
	cloudFac := states.GetFloatState("weather.currently:clouds", 0.0)
	for _, w := range l.weather {
		if cloudFac <= w.clouds {
			clouds = w
			break
		}
	}

	// Full cloud cover will increase color-temperature by 10% of (Max - Current)
	// NOTE: Only during the day (twilight + light)
	// TODO: when the moon is shining in the night the amount
	//       of blue-light is also higher than normal.
	// CT = 0.0 -> Coldest (>6500K)
	// CT = 1.0 -> Warmest (2000K)
	CT := current.ct[0] + currentx*(current.ct[1]-current.ct[0])
	if current.darkorlight != "dark" {
		if clouds.ct_pct >= 0 {
			CT = CT + clouds.ct_pct*(1.0-CT)
		} else {
			CT = CT - clouds.ct_pct*CT
		}
	}
	CT = season.minCT + (CT * (season.maxCT - season.minCT))

	// Full cloud cover will increase brightness by 10% of (Max - Current)
	// BRI = 0 -> Very dim light
	// BRI = 1 -> Very bright light
	BRI := current.bri[0] + currentx*(current.bri[1]-current.bri[0])
	BRI = BRI + cloudFac*0.1*(1.0-BRI)
	if current.darkorlight != "dark" {
		// A bit brighter lights when there are clouds during the day.
		if clouds.ct_pct >= 0 {
			BRI = BRI + clouds.bri_pct*(1.0-BRI)
		} else {
			BRI = BRI - clouds.bri_pct*BRI
		}
	}
	BRI = season.minBRI + (BRI * (season.maxBRI - season.minBRI))

	// TODO: Put into configuration TOML
	// Update 'state' for the following
	// - sensor:lights_HUE_CT
	// - sensor:lights_HUE_BRI
	// - sensor:lights_YEE_CT
	// - sensor:lights_YEE_BRI
	for _, ltype := range l.lighttypes {
		lct := ltype.minCT + CT*(ltype.maxCT-ltype.minCT)
		lbri := ltype.minBRI + BRI*(ltype.maxBRI-ltype.minBRI)
		states.SetFloatState("lighting.lights_"+ltype.name+"_ct", math.Floor(lct))
		states.SetFloatState("lighting.lights_"+ltype.name+"_bri", math.Floor(lbri))
	}

	states.SetFloatState("lighting.lights_ct", float64(int64(CT*100.0))/100.0)
	states.SetFloatState("lighting.lights_bri", float64(int64(BRI*100.0))/100.0)
	states.SetStringState("lighting.darklight", current.darkorlight)

	return 30 * time.Second
}
