package lights

import "time"

// Color Temperature
// URL: https://panasonic.net/es/solution-works/jiyugaoka/

type lightcfg struct {
	kelvin      int
	brightness  float64
	startMoment string
	startTime   time.Time
	endMoment   string
	endTime     time.Time
}

var defaultTime time.Time

var lighting = []lightcfg{
	{2000, 1, "night:dawn", defaultTime, "", defaultTime},
	{2100, 50, "astronomical:dawn", defaultTime, "", defaultTime},
	{2300, 100, "nautical:dawn", defaultTime, "", defaultTime},
	{2500, 100, "civil:dawn", defaultTime, "", defaultTime},
	{2700, 90, "sunrise", defaultTime, "3:00", defaultTime},
	{4700, 80, "", defaultTime, "", defaultTime},
	{6000, 80, "sun:noon", defaultTime, "sun:noon", defaultTime},
	{6000, 80, "sun:afternoon", defaultTime, "", defaultTime},
	{3500, 90, "sunset", defaultTime, "", defaultTime},
	{3200, 80, "civil:dusk", defaultTime, "", defaultTime},
	{2800, 80, "nautical:dusk", defaultTime, "", defaultTime},
	{2700, 80, "astronomical:dusk", defaultTime, "", defaultTime},
	{2700, 80, "", defaultTime, "1:00", defaultTime},
	{2700, 50, "", defaultTime, "1:00", defaultTime},
	{2000, 20, "", defaultTime, "1:00", defaultTime},
	{2000, 1, "", defaultTime, "", defaultTime},
	{2000, 1, "night:darkest", defaultTime, "", defaultTime},
}
