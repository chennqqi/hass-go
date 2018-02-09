# task graph with dependency and property changed triggers


NOTE: A tick property has additional properties 'update:last' and 'update:next'

type Trigger interface {
    func Trigger(previous, current string)
}

type Type int
const (
    TypeBool Type = iota
    TypeString Type
    TypeInt Type
    
)

type Property struct {
    // B bool      -> encoded into I
    // S string
    // I int64     
    // F float64   -> encoded into I
    // T time.Time -> encoded in F
    T Type
}


aqi has the following triggers
- property aqi.tick

calendar has the following triggers
- property calendar.tick

hass has the following triggers
NOTE: Maybe a hass.sensor makes more sense
- anything written to a 'sensor.NAME' property

lighting has the following triggers
- property weather.currently:clouds
- property time.season
- property lighting.tick

reporter has the following triggers
- anything written to the 'reporter' domain

sensors has the following triggers
NOTE: Maybe a sensors.sensor makes more sense to be triggered
- all properties as listed in the configuration

suncalc has the following triggers
- property time.day

timeofday has the following triggers
- property timeofday.tick

weather has the following triggers
- property weather.tick

