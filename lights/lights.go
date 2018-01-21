package lights

import "time"

// Color Temperature
// URL: https://panasonic.net/es/solution-works/jiyugaoka/

type lighttime struct {
	kelvin      float64
	brightness  float64
	startMoment string
	endMoment   string
}

var defaultTime time.Time

type addtime struct {
	name  string
	tag   string
	shift time.Duration
}

var additionaltimes = []addtime{
	{"sunrise:end", ":+2h", 2 * time.Hour},
	{"sun:noon:end", ":+2h", 2 * time.Hour},
	{"astronomical:dusk:end", ":+1h", 1 * time.Hour},
	{"astronomical:dusk:end", ":+2h", 2 * time.Hour},
	{"astronomical:dusk:end", ":+3h", 3 * time.Hour},
}

var lighttable = []lighttime{
	{0.0, 0.01, "night:darkest:end", "astronomical:dawn:begin"},
	{0.03, 0.50, "astronomical:dawn:begin", "astronomical:dawn:end"},
	{0.06, 0.60, "nautical:dawn:begin", "nautical:dawn:end"},
	{0.11, 0.80, "civil:dawn:begin", "civil:dawn:end"},
	{0.15, 0.90, "sunrise:begin", "sunrise:end:+2h"},
	{0.6, 0.80, "sunrise:end:+2h", "sun:noon:begin"},
	{0.7, 0.85, "sun:noon:begin", "sun:noon:end"},
	{0.7, 0.85, "sun:noon:end", "sun:noon:end:+2h"},
	{0.7, 0.85, "sun:noon:end:+2h", "sunset:begin"},
	{0.22, 0.85, "sunset:begin", "sunset:end"},
	{0.21, 0.80, "civil:dusk:begin", "civil:dusk:end"},
	{0.2, 0.80, "nautical:dusk:begin", "nautical:dusk:end"},
	{0.15, 0.80, "astronomical:dusk:begin", "astronomical:dusk:end"},
	{0.13, 0.80, "astronomical:dusk:end", "astronomical:dusk:end:+1h"},
	{0.08, 0.50, "astronomical:dusk:end:+1h", "astronomical:dusk:end:+2h"},
	{0.04, 0.20, "astronomical:dusk:end:+2h", "astronomical:dusk:end:+3h"},
	{0.0, 0.01, "astronomical:dusk:end:+3h", "night:darkest:begin"},
	{0.0, 0.01, "night:darkest:begin", "night:darkest:end"},
}

type lightcfg struct {
	name   string
	minCT  float64
	maxCT  float64
	minBRI float64
	maxBRI float64
}

var lightcfgs = []lightcfg{
	{"HUE", 2000, 6500, 0, 254},
	{"Yee", 2000, 6500, 0, 254},
}

var seasonModifiers = map[string]float64{
	"Winter": 0.8,
	"Spring": 0.75,
	"Summer": 0.70,
	"Autumn": 0.75,
}

type weathermod struct {
	weather float64 // Cloud factor
	ct      float64
	bri     float64
}

var weather = []weathermod{
	{0.0, 0.9, 1.0},
	{0.5, 0.92, 1.0},
	{0.15, 0.97, 1.0},
	{0.35, 1.04, 1.0},
	{0.5, 1.12, 1.0},
	{0.9, 1.2, 1.0},
}

type Lights struct {
}

func get(k string, m map[string]string) string {
	v, x := m[k]
	if x {
		return v
	}
	return ""
}

// Process will update 'string'states and 'float'states
// States are both input and output, for example as input
// there are Season/Weather states like 'Season':'Winter'
// and 'Clouds':0.5
func (l *Lights) Process(sstates *map[string]string, fstates *map[string]float64, tstates *map[string]time.Time) {

	season := get("Season", *sstates)
	seasonModifier := seasonModifiers[season]

}
